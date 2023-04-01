cd lib
set GOOS=windows&& set GOARCH=amd64&& go build -o ../bin/sftp-win32-amd64.exe main.go
set GOOS=windows&& set GOARCH=386&& go build -o ../bin/sftp-win32-ia32.exe main.go
set GOOS=darwin&& set GOARCH=amd64&& go build -o ../bin/sftp-darwin-x64 main.go
set GOOS=darwin&& set GOARCH=arm64&& go build -o ../bin/sftp-darwin-arm64 main.go
set GOOS=linux&& set GOARCH=amd64&& go build -o ../bin/sftp-linux-x64 main.go
set GOOS=linux&& set GOARCH=arm64&& go build -o ../bin/sftp-linux-arm64 main.go
set GOOS=linux&& set GOARCH=386&& go build -o ../bin/sftp-linux-ia32 main.go