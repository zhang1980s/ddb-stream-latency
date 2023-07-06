build:
	if [ -d bin ]; then rm -rf bin/*; fi;
	GOARCH=amd64 GOOS=linux go build -o bin/reader reader.go
	zip bin/reader.zip bin/reader