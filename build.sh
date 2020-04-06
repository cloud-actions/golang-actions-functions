mkdir -p _/

echo "GOOS=linux GOARCH=amd64 go build"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o _/hello-gopher_linux .
echo "GOOS=darwin GOARCH=amd64 go build"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o _/hello-gopher_darwin .
echo "GOOS=windows GOARCH=amd64 go build"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o _/hello-gopher_windows.exe .
