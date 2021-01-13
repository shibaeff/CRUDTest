GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath handler.go controllers.go  db.go
