#!/usr/bin/env bash

BIN_ROOT="$(cd "$(dirname "$BASH_SOURCE[0]")"; pwd)"

main() {
	cd "$BIN_ROOT/.."
	grep "//go:generate" -l \
		--include "*.go" \
		--recursive . \
		| xargs dirname | sort | uniq \
		| xargs go generate
	
}

main "$@"