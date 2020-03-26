# README

## development
```bash
# init go module
go mod init github.com/asw101/hello-gopher

go run .

go build .
```

## build & copy
```bash
# build locally for linux, mac, windows
source build-local.sh

# build with golang:1.13.2 container (for windows/app service)
source build-container.sh

# build w/o copy
source build.sh
```
