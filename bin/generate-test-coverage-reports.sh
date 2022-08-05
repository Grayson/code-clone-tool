#!/usr/bin/env bash

REPO_ROOT="$(cd "$(dirname "$BASH_SOURCE[0]")"; pwd)"

main() {
	local packages=(
		"./"
		"./lib"
		"./lib/either"
		"./lib/fs"
		"./lib/gitapi"
		"./lib/githubapi"
		"./lib/optional"
		"./lib/shell"
		"./lib/stage"
	)

	local outdir="./coverage_reports"
	mkdir -p "$outdir"

	for x in "${packages[@]}"; do
		local name="$(basename "$x")"
		while [[ "." == "${name:0:1}" ]]; do
			name="${name#.}"
		done
		name="${name:-main}"

		local profile="${outdir}/${name}.coverage.out"
		local output="${outdir}/${name}.coverage.html"
		go test "$x" -coverprofile="$profile"
		go tool cover -html="$profile" -o "${output}"
	done
}

main "$@"
