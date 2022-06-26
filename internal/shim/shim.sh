#!/bin/sh

# placeholder for github.com/frantjc/sequence/internal/cmd/sqnc-shim
# which overrides this at build time
#
# sqnc-shim and shim.sh should always be copies of one another while
# in source control, the latter only existing the replace the former
# if it gets accidentally overridden

echo '::error::this should not happen'
echo '::error::github.com/frantjc/sequence/internal/shim.Bytes is incorrect--did you build sequence properly?'
echo '::error::attempting to salvage...'
exec ${@:1}
