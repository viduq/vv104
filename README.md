# vv104
IEC 60870-5-104 (IEC 104) server and client implementation with test features.

Work in progress.

- Run as client with `go run cmd/main.go`
- Use `-s` to run as server.
- Use `-h 192.168.0.1` to connect to or listen on specific hostname or IP address. Default is localhost.
- Use `-p 2405` to connect or listen on specific port. Default is `2404`.
- For more cli flags use `--help`

This repository is mainly used in my other repository [https://github.com/viduq/brez104](https://github.com/viduq/brez104) which is a GUI program for IEC 104.
