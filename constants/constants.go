package constants

var Const = struct {
	DockerLocaldevServerImage string
	DockerNetwork             string
	RepoDir                   string
	RepoRemote                string
	RepoBranch                string
	RepoSync                  string
}{
	DockerLocaldevServerImage: "docker.localdev.server.image",
	DockerNetwork:             "docker.network",
	RepoDir:                   "repo.dir",
	RepoRemote:                "repo.remote",
	RepoBranch:                "repo.branch",
	RepoSync:                  "repo.sync",
}
