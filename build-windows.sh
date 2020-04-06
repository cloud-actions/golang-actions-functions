mkdir -p bin/

echo "GOOS=windows GOARCH=amd64 go build"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/hello-gopher_windows.exe .
