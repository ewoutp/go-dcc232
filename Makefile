PROJECT=dcc232test

all: binaries

clean: 
	rm -Rf bin

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm" \
		-ldflags="-X main.projectVersion=${VERSION} -X main.projectBuild=${COMMIT}" \
		-output="bin/{{.OS}}/{{.Arch}}/$(PROJECT)" \
		-tags="netgo" \
		github.com/ewoutp/go-dcc232/cmd/dcc232test

test:
	go test github.com/ewoutp/go-dcc232

