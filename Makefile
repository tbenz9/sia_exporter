PLATFORMS := linux/amd64 windows/amd64 darwin/amd64 linux/arm64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

defualt: build

build: 
	GO111MODULE=on go build -o sia_exporter *.go

install:
	GO111MODULE=on go install -o sia_exporter *.go

release: $(PLATFORMS)

$(PLATFORMS):
	GO111MODULE=on GOOS=$(os) GOARCH=$(arch) go build -o 'sia_exporter-$(os)-$(arch)' *.go
