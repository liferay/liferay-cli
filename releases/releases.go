package releases

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/http"
)

type Releases struct {
}

func init() {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault(constants.Const.ReleasesFile, filepath.Join(dirname, ".liferay", "cli", "releases.json"))
	viper.SetDefault(constants.Const.ReleasesURL, "https://releases-cdn.liferay.com/tools/workspace/.product_info.json")
}

func ReleaseVersions(verbose bool) ([]string, error) {
	dat, err := ReleasesJSON(verbose)

	if err != nil {
		return nil, err
	}

	keys := make([]string, len(dat))

	i := 0
	for k := range dat {
		keys[i] = k
		i++
	}

	return keys, nil
}

func GetReleaseByVersion(version string, verbose bool) (interface{}, error) {
	dat, err := ReleasesJSON(verbose)

	if err != nil {
		return nil, err
	}

	return dat[version], nil
}

func ReleasesJSON(verbose bool) (map[string]interface{}, error) {
	var dat map[string]interface{}

	bytes, err := http.GetOrFetchBytes(http.GetOrFetchBytesOptions{
		EtagKey: constants.Const.ReleasesEtag,
		FileKey: constants.Const.ReleasesFile,
		URL:     viper.GetString(constants.Const.ReleasesURL),
		Verbose: verbose,
	})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &dat); err != nil {
		return nil, err
	}

	return dat, nil
}
