package shim

import _ "embed"

var (
	//go:embed sqnc-shim-source
	SqncShimSource []byte

	//go:embed sqnc-shim-uses
	SqncShimUses []byte
)
