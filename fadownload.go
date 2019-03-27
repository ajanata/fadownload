/*
 *
 * Copyright (c) 2019, Andy Janata
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted
 * provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this list of conditions
 *   and the following disclaimer.
 * * Redistributions in binary form must reproduce the above copyright notice, this list of
 *   conditions and the following disclaimer in the documentation and/or other materials provided
 *   with the distribution.
 * * Neither the name of the copyright holder nor the names of its contributors may be used to
 *   endorse or promote products derived from this software without specific prior written
 *   permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
 * FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR
 * CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
 * WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY
 * WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ajanata/faapi"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Load up the config
	c := loadConfig()

	if c.User == "" || c.Output == "" {
		fmt.Println("No -user or -output flags")
		return
	}

	// Configure logging
	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		log.WithField("level", c.LogLevel).Warn("Unable to parse log level, using INFO")
		level = log.InfoLevel
	}
	log.SetLevel(level)

	if c.LogJSON {
		log.SetFormatter(&log.JSONFormatter{})
	} else if c.LogForceColors {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Turn on pprof debugging if requested
	if c.Debug {
		go func() {
			log.Info(http.ListenAndServe("localhost:6680", nil))
		}()
	}

	// Create FurAffinity API client.
	fa, err := faapi.New(c.FA.faAPIConfig())
	if err != nil {
		log.WithError(err).Fatal("Unable to create faapi client!")
	}
	defer fa.Close()

	username, err := fa.GetUsername()
	if err != nil {
		log.WithError(err).Warn("Not logged in to FurAffinity! Only artwork rated general will be downloaded.")
	} else {
		log.WithField("username", username).Info("Logged in to FurAffinity.")
	}

	err = os.MkdirAll(c.Output, 0755)
	if err != nil {
		log.WithError(err).Error("Unable to create output directory")
		os.Exit(1)
	}

	u := fa.NewUser(c.User)
	p := uint(0)
	for {
		p++
		subs, err := u.GetSubmissions(p)
		if err != nil {
			log.WithError(err).Error("Unable to get submissions")
			os.Exit(1)
		}
		if len(subs) == 0 {
			break
		}

		log.WithField("page", p).Info("Got page")
		for _, sub := range subs {
			slog := log.WithField("title", sub.Title)
			slog.Info("Downloading")
			details, err := sub.Details()
			if err != nil {
				slog.WithError(err).Error("Unable to load details")
				os.Exit(1)
			}

			bb, err := details.Download()
			if err != nil {
				slog.WithError(err).Error("Unable to download")
			}

			split := strings.Split(details.DownloadURL, "/")
			name := fmt.Sprintf("%s/%s", c.Output, split[len(split)-1])
			err = ioutil.WriteFile(name, bb, 0644)
			if err != nil {
				slog.WithError(err).Error("Unable to save download")
				os.Exit(1)
			}
		}
	}
	log.Info("Done.")
}
