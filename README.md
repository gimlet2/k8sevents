# k8sevents

Project to watch on k8s events and log them in json format.

## Build
`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo`
`docker build .`
