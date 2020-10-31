litsen_linux: always
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o litsen_linux -mod=vendor ./cmd/litsen

always:

clean:
	rm -f litsen_linux

.PHONY: always clean
