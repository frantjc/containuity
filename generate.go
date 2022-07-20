package sequence

//go:generate buf format -w

//go:generate buf generate .

//go:generate env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./internal/shim ./internal/cmd/sqnc-shim

//go:generate upx ./internal/shim/sqnc-shim
