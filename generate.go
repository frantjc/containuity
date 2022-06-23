package sequence

import _ "github.com/bufbuild/buf/cmd/buf" //nolint:typecheck

//go:generate go run -tags generate github.com/bufbuild/buf/cmd/buf format -w

//go:generate go run -tags generate github.com/bufbuild/buf/cmd/buf generate .

//go:generate go fmt ./...
