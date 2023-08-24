package constants

var Const = struct {
	CETypesEtag                   string
	CETypesFile                   string
	CETypesURL                    string
	CliReleasesEtag               string
	CliReleasesFile               string
	CliReleasesURL                string
	CliUpdateCheck                string
	DockerLocaldevServerImage     string
	DockerLocaldevServerPullimage string
	DockerNetwork                 string
	ExtWorkspaceDir               string
	ExtWorkspaceDirSpecified      string
	ReleasesEtag                  string
	ReleasesFile                  string
	ReleasesURL                   string
	RepoDir                       string
	RepoRemote                    string
	RepoBranch                    string
	RepoSync                      string
	TlsLfrdevDomain               string
}{
	CETypesEtag:                   "cetypes.etag",
	CETypesFile:                   "cetypes.file",
	CETypesURL:                    "cetypes.url",
	CliReleasesEtag:               "cli.releases.etag",
	CliReleasesFile:               "cli.releases.file",
	CliReleasesURL:                "cli.releases.url",
	CliUpdateCheck:                "cli.update.check",
	DockerLocaldevServerImage:     "docker.localdev.server.image",
	DockerLocaldevServerPullimage: "docker.localdev.server.pull-image",
	DockerNetwork:                 "docker.network",
	ExtWorkspaceDir:               "extension.workspace.dir",
	ExtWorkspaceDirSpecified:      "extension.workspace.dir-specified",
	ReleasesEtag:                  "releases.etag",
	ReleasesFile:                  "releases.file",
	ReleasesURL:                   "releases.url",
	RepoDir:                       "localdev.resources.dir",
	RepoRemote:                    "localdev.resources.remote",
	RepoBranch:                    "localdev.resources.branch",
	RepoSync:                      "localdev.resources.sync",
	TlsLfrdevDomain:               "tls.lfrdev.domain",
}
