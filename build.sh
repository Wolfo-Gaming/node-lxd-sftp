cd lib
GOOS=windows && GOARCH=amd64 && go build -o ../bin/sftp-windows-amd64.exe main.go
GOOS=windows && GOARCH=386 && go build -o ../bin/sftp-windows-ia32.exe main.go
GOOS=darwin && GOARCH=amd64 && go build -o ../bin/sftp-darwin-x64 main.go
GOOS=darwin && GOARCH=arm64 && go build -o ../bin/sftp-darwin-arm64 main.go
GOOS=linux && GOARCH=amd64 && go build -o ../bin/sftp-linux-x64 main.go
GOOS=linux && GOARCH=arm64 && go build -o ../bin/sftp-linux-arm64 main.go
GOOS=linux && GOARCH=386 && go build -o ../bin/sftp-linux-ia32 main.go