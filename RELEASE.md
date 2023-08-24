# Release steps for liferay cli

# Install at least goreleaser version 1.20.0

* Push a tag to the repo which begins with "v"
* execute
  ```bash
  git push upstream
  git tag "v${VERSION}"
  git push upstream --tags
  rm -rf dist
  GITHUB_TOKEN=<your_gh_token> goreleaser release --clean
  ```
