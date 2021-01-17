.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/test test/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/create create/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/read read/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/update update/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete delete/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
