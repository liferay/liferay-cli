package constants

var Const = struct {
	DockerLocaldevServerImage string
	DockerNetwork             string
	ReleasesEtag              string
	ReleasesFile              string
	ReleasesURL               string
	RepoDir                   string
	RepoRemote                string
	RepoBranch                string
	RepoSync                  string
}{
	DockerLocaldevServerImage: "docker.localdev.server.image",
	DockerNetwork:             "docker.network",
	ReleasesEtag:              "releases.etag",
	ReleasesFile:              "releases.file",
	ReleasesURL:               "releases.url",
	RepoDir:                   "repo.dir",
	RepoRemote:                "repo.remote",
	RepoBranch:                "repo.branch",
	RepoSync:                  "repo.sync",
}
