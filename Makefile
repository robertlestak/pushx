VERSION=v0.0.8

.PHONY: pushx
pushx: clean bin/pushx_darwin bin/pushx_windows bin/pushx_linux

bin/pushx_darwin:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/pushx_darwin cmd/pushx/*.go
	openssl sha512 bin/pushx_darwin > bin/pushx_darwin.sha512

bin/pushx_linux:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/pushx_linux cmd/pushx/*.go
	openssl sha512 bin/pushx_linux > bin/pushx_linux.sha512

bin/pushx_hostarch:
	mkdir -p bin
	go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/pushx_hostarch cmd/pushx/*.go
	openssl sha512 bin/pushx_hostarch > bin/pushx_hostarch.sha512

bin/pushx_windows:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/pushx_windows cmd/pushx/*.go
	openssl sha512 bin/pushx_windows > bin/pushx_windows.sha512

.PHONY: envvars
envvars:
	egrep -oh --exclude Makefile \
		--exclude-dir bin \
		--exclude-dir scripts \
		-R 'os.Getenv\(.*?\)' . | \
		tr -d ' ' | \
		sort | \
		uniq | \
		sed -e 's,os.Getenv(,,g' -e 's,),,g' \
		-e 's,",,g' \
		-e 's,prefix+,PUSHX_,g'

.PHONY: envvarsyaml
envvarsyaml:
	bash scripts/envvarsyaml.sh

.PHONY: clean
clean:
	rm -rf bin

.PHONY: slim
slim:
	bash scripts/build_drivers.sh build $(drivers)

.PHONY: listdrivers
listdrivers:
	bash scripts/build_drivers.sh list
