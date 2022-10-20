package prereq

import (
	"liferay.com/liferay/cli/docker"
	"liferay.com/liferay/cli/git"
)

func Prereq(verbose bool) {
	git.SyncGit(verbose)
	docker.BuildImages(verbose)
}
