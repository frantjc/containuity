#!/bin/sh

# placeholder for github.com/frantjc/sequence/internal/cmd/sqnc-shim
# which overrides this at build time
#
# sqnc-shim and shim.sh should always be copies of one another while
# in source control, the latter only existing the replace the former
# if it gets accidentally overridden

echo '::error::this should not happen'
echo '::error::github.com/frantjc/sequence/internal/shim.Bytes is incorrect--did you build sequence properly?'

if [ -n "$_SQNC_SHIM_SWITCH"]; then
    echo '::error::unable to salvage...'
    exit 1
fi

echo '::error::attempting to salvage...'

if [ -n "$GITHUB_ENV" ]; then
    echo "::info::sourcing $GITHUB_ENV"
    source $GITHUB_ENV
fi

if [ -n "$GITHUB_PATH" ]; then
    echo "::info::adding $GITHUB_PATH to PATH"
    for $path in $(cat $GITHUB_PATH); do
        PATH=$PATH:$path
    done
fi

exec ${@:1}
