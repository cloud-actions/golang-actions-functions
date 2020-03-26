docker run --rm -v ${PWD}:/pwd/ -w /pwd/ -i golang:1.13.2 bash build.sh

# copy to ../hello-serverless-go/
mkdir -p ../hello-serverless-go/bin/
cp '_/hello-gopher_linux' ../hello-serverless-go/bin/
cp '_/hello-gopher_windows.exe' ../hello-serverless-go/bin/
