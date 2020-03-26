# README

## func

```bash
docker build -t func .
docker run -p 7071:7071 -v ${PWD}:/pwd/ -w /pwd/ --rm -it func:latest bash
func start
```
