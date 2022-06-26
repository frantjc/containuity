package sequence

// dummy import so that go mod doesn't remove
// the github.com/bufbuild/buf module and break
// the go:generate comments in this file
import _ "github.com/bufbuild/buf/private/buf/cmd/buf" //

//go:generate go run -tags generate github.com/bufbuild/buf/cmd/buf format -w

//go:generate go run -tags generate github.com/bufbuild/buf/cmd/buf generate .

//go:generate go fmt ./...
