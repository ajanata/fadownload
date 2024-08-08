fadownload - Download all submissions for a user from FurAffinity

Copy fadownload.example.toml to fadownload.toml and fill in the required values for FurAffinity cookies, or you will only get General-rated pieces. Get them from your web browser's dev tools.

Make sure that the account you use with this has the Classic UI template selected, or some things may not work right.

Usage:
go run . -user user-to-download -output directory/to/store/files
