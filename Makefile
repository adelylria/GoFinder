run:
	go run cmd/main.go
build-exe:
	go build -ldflags="-H=windowsgui" -o build/goFinder.exe ./cmd/main.go
build-dll:
	gcc -shared -o hotkey/hotkey.dll hotkey/hotkey.c -Wall -mwindows
