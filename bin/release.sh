#!/usr/bin/env bash

REPO_ROOT="$(cd "$(dirname "$BASH_SOURCE[0]")"; pwd)"

check_goreleaser() {
	if [[ "" == "$(which goreleaser)" ]]; then
		echo "goreleaser was not found."
		echo "Please install using instructions from https://goreleaser.com/install/"
		exit 1
	fi
}

ensure_github_token() {
	if [[ "" == $GITHUB_TOKEN ]]; then
		read -t 10 -p "Enter your Github Token (minimum \`write:packages\`): " token
		export GITHUB_TOKEN="$token"
	fi
	if [[ "" == $GITHUB_TOKEN ]]; then
		echo "No valid github token provided"
		exit 1
	fi
}

ensure_tag() {
	if [[ "" == $RELEASE_TAG ]]; then
		read -t 10 -p "Enter your release tag (using semver): " token
		export RELEASE_TAG="$token"
	fi
	if [[ "" == $RELEASE_TAG ]]; then
		echo "No valid tag provided"
		exit 1
	fi
}

main() {
	check_goreleaser
	ensure_github_token
	ensure_tag

	git tag \
		--annotate "$RELEASE_TAG" \
		--message "Release $RELEASE_TAG"
	git push origin "$RELEASE_TAG"

	goreleaser release --rm-dist
}

main "$@"