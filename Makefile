VERSION := $(or $(AppVersion), "v0.0.0")
COMMIT := $(or $(shell git rev-parse --short HEAD), "unknown")
BUILDDATE := $(shell date +%Y-%m-%d)

LDFLAGS := -X 'main.AppVersion=$(VERSION)' -X 'main.CommitHash=$(COMMIT)' -X 'main.BuildDate=$(BUILDDATE)' -s -w

all: build

tidy:
	go mod tidy

test:
	go test -v ./...  -race -coverprofile=coverage.out -covermode=atomic

coverage: test
	go tool cover -func=coverage.out

coverage-html: test
	mkdir -p coverage
	go tool cover -html=coverage.out -o coverage/index.html

coverage-serve: coverage-html
	python3 -m http.server 8080 -d coverage

install: build
	cp eyez /usr/local/bin/eyez
	cp man/eyez.1 /usr/local/share/man/man1/eyez.1

uninstall:
	rm /usr/local/bin/eyez
	rm /usr/local/share/man/man1/eyez.1

build:
	go build -gcflags="all=-N -l" -ldflags="$(LDFLAGS)" -o eyez

dist:
	cp man/eyez.1 man/eyez.old
	sed -e "s|BUILDDATE|$(BUILDDATE)|g" -e "s|VERSION|$(VERSION)|g" man/eyez.old > man/eyez.1

	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/eyez-linux-amd64
	cp build/eyez-linux-amd64 build/eyez
	tar -zcvf build/eyez-linux-amd64.tar.gz build/eyez man/eyez.1

	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o build/eyez-linux-arm64
	cp build/eyez-linux-arm64 build/eyez
	tar -zcvf build/eyez-linux-arm64.tar.gz build/eyez man/eyez.1

	GOOS=linux GOARCH=arm go build -ldflags="$(LDFLAGS)" -o build/eyez-linux-arm
	cp build/eyez-linux-arm build/eyez
	tar -zcvf build/eyez-linux-arm.tar.gz build/eyez man/eyez.1

	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/eyez-darwin-amd64
	cp build/eyez-darwin-amd64 build/eyez
	tar -zcvf build/eyez-darwin-amd64.tar.gz build/eyez man/eyez.1

	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o build/eyez-darwin-arm64
	cp build/eyez-darwin-arm64 build/eyez
	tar -zcvf build/eyez-darwin-arm64.tar.gz build/eyez man/eyez.1
	rm build/eyez

	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o build/eyez-windows-amd64.exe
	GOOS=windows GOARCH=386 go build -ldflags="$(LDFLAGS)" -o build/eyez-windows-i386.exe

	# Generating checksum
	cd build && sha256sum * >> checksum-sha256sum.txt
	cd build && md5sum * >> checksum-md5sum.txt

	# Cleaning
	mv man/eyez.old man/eyez.1

clean:
	rm -rf eyez*
	rm -rf build
