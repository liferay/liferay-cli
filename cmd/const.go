package cmd

var Const = struct {
	dockerLocaldevServerImage string
	dockerNetwork             string
	repoDir                   string
	repoRemote                string
	repoBranch                string
	repoSync                  string
}{
	dockerLocaldevServerImage: "docker.localdev.server.image",
	dockerNetwork:             "docker.network",
	repoDir:                   "repo.dir",
	repoRemote:                "repo.remote",
	repoBranch:                "repo.branch",
	repoSync:                  "repo.sync",
}
