#!/bin/sh

# placeholder for github.com/frantjc/sequence/internal/cmd/sqnc-shim-uses
# or github.com/frantjc/sequence/internal/cmd/sqnc-shim-source
# which overrides this at build time
#
# source and uses are copies of shim.sh in case one of the former
# gets accidentally overwritten in source control :)

echo '::error::this should not happen'
echo '::error::github.com/frantjc/sequence/workflow.Shim is incorrect--did you build sequence properly?'
echo '::error::attempting to salvage...'
exec ${@:1}
