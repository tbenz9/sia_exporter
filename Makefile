defualt: build

build: 
	GO111MODULE=on go build -o sia_exporter *.go

install:
	GO111MODULE=on go install -o sia_exporter *.go
