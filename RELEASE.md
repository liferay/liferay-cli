# Release steps for liferay cli

* Push a tag to the repo which begins with "v"
* execute
  ```bash
  GITHUB_TOKEN=<your_gh_token> goreleaser release --rm-dist
  ```
