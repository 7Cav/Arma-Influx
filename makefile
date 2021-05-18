build:
	export GOARCH="amd64"
	export CGO_ENABLED=1
	go build -o armago_x64.so -buildmode=c-shared .
