darwin:
	go build -o ./code2azw3-darwin

linux:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o ./code2azw3-linux64

win:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "-s -w" -a -installsuffix cgo -o ./code2azw3-win.exe
