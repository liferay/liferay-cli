package ext

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"liferay.com/liferay/cli/constants"
)

var replacements = regexp.MustCompile(`[^a-zA-Z0-9_.-]`)

func GetExtensionDirKey() string {
	path := viper.GetString(constants.Const.ExtClientExtensionDir)

	dirPath := filepath.Dir(path)
	baseDirName := filepath.Base(path)

	return MakeExtensionDirKey(
		baseDirName, dirPath, string(filepath.Separator))
}

func MakeExtensionDirKey(baseDirName string, dirPath string, pathSeparator string) string {
	parts := strings.Split(dirPath, pathSeparator)

	newParts := make([]string, 0)

	for _, part := range parts {
		part = strings.ToLower(part)
		length := len(part)
		if length > 2 {
			newParts = append(newParts, part[0:2])
		} else if length > 0 && length <= 2 {
			newParts = append(newParts, part)
		}
	}

	morphed := strings.Join(newParts, "-") + "-" + baseDirName

	return replacements.ReplaceAllString(morphed, "")
}
