# lib-sftp-go

lib for sending file over sftp

## install

clone the repo

```
go mod tidy
cd cmd
```

run as sudo
`go run main.go upload --host=123.209.119.15 --port=2022 --user=pi --pass=mypass`

```
go build ctl.go && sudo ./ctl  --help
```

## docs

[CLI](docs/cmd.md)

