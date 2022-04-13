#!/bin/sh

# placeholder for github.com/frantjc/sequence/cmd/sqncshim-uses
# or github.com/frantjc/sequence/cmd/sqncshim
# which overrides this at build time
#
# sqncshim and sqncshim-uses are copies of sqncshim.sh in case one of the former
# gets accidentally overwritten in source control :)

echo '::error::this should not happen'
echo '::error::github.com/frantjc/sequence/workflow.Shim is incorrect--did you build sequence properly?'
echo '::error::attempting to salvage...'
exec ${@:1}
