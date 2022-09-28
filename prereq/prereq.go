package prereq

import (
	"liferay.com/lcectl/docker"
	"liferay.com/lcectl/git"
)

func Prereq(verbose bool) {
	git.SyncGit(verbose)
	docker.BuildImages(verbose)
}
