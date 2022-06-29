package shim

import _ "embed"

var (
	//go:embed sqnc-shim
	Bytes []byte
)
