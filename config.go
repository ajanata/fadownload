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
	"time"

	"github.com/koding/multiconfig"

	"github.com/ajanata/faapi"
)

type (
	// Config is the configuration for the bot.
	Config struct {
		User   string
		Output string `default:"."`
		Limit  int    `default:"2147483647"`

		Debug          bool   `default:"false"`
		LogLevel       string `default:"INFO"`
		LogForceColors bool   `default:"false"`
		LogJSON        bool   `default:"false"`

		FA FA
	}

	// FA is the configuration for FurAffinity.
	FA struct {
		Cookies   []Cookie
		Proxy     string
		RateLimit duration `required:"true"`
		// RequestTimeout is the timeout for a single attempt at the request.
		RequestTimeout duration
		RetryDelay     duration
		RetryLimit     int
		// Timeout is the timeout on the entire request, including retries.
		Timeout   duration
		UserAgent string `required:"true"`
	}

	// Cookie is an HTTP cookie.
	Cookie struct {
		Name  string
		Value string
	}

	duration struct {
		time.Duration
	}
)

func loadConfig() *Config {
	m := multiconfig.NewWithPath("fadownload.toml")
	c := new(Config)
	m.MustLoad(c)
	return c
}

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func (c *FA) faAPIConfig() faapi.Config {
	cookies := make([]faapi.Cookie, len(c.Cookies))
	for i, cookie := range c.Cookies {
		cookies[i] = faapi.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		}
	}
	return faapi.Config{
		Cookies:        cookies,
		Proxy:          c.Proxy,
		RateLimit:      c.RateLimit.convert(),
		RequestTimeout: c.RequestTimeout.convert(),
		RetryDelay:     c.RetryDelay.convert(),
		RetryLimit:     c.RetryLimit,
		Timeout:        c.Timeout.convert(),
		UserAgent:      c.UserAgent,
	}
}

func (d duration) convert() time.Duration {
	// this is so dumb
	td, err := time.ParseDuration(d.String())
	if err != nil {
		panic(err)
	}
	return td
}
