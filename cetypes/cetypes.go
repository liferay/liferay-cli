package cetypes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"liferay.com/liferay/cli/constants"
	"liferay.com/liferay/cli/http"
	"liferay.com/liferay/cli/workspace"
)

func init() {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault(constants.Const.CETypesFile, filepath.Join(dirname, ".liferay", "cli", "client-extension-types.json"))
	viper.SetDefault(constants.Const.CETypesURL, "https://raw.githubusercontent.com/liferay/liferay-portal/%s/modules/apps/client-extension/client-extension-type-api/src/main/resources/com/liferay/client/extension/type/dependencies/client-extension-types.json")
}

func ClientExtensionTypeKeys(verbose bool) ([]string, error) {
	dat, err := ClientExtensionTypesJSON(verbose)

	if err != nil {
		return nil, err
	}

	keys := make([]string, len(dat))

	i := 0
	for k := range dat {
		entry := dat[k]
		keys[i] = entry["name"].(string)
		i++
	}

	return keys, nil
}

func ClientExtensionTypesJSON(verbose bool) ([]map[string]interface{}, error) {
	var dat []map[string]interface{}

	tag, err := workspace.GetProductVersionAsTag(verbose)

	if err != nil {
		return nil, err
	}

	bytes, err := http.GetOrFetchBytes(http.GetOrFetchBytesOptions{
		EtagKey: constants.Const.CETypesEtag,
		FileKey: constants.Const.CETypesFile,
		URL:     fmt.Sprintf(viper.GetString(constants.Const.CETypesURL), tag),
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
