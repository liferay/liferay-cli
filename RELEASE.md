# Release steps for liferay cli

* Push a tag to the repo which begins with "v"
* execute
  ```bash
  git tag "v${VERSION}"
  git push upstream --tags
  rm -rf dist
  GITHUB_TOKEN=<your_gh_token> goreleaser release --rm-dist
  ```
