package workspace

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/magiconair/properties"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
	"liferay.com/lcectl/releases"
)

func GetProductVersion() (string, error) {
	repoDir := viper.GetString(constants.Const.RepoDir)

	p, err := properties.LoadFile(
		path.Join(repoDir, "docker", "images", "localdev-server", "workspace", "gradle.properties"), properties.UTF8)

	if err != nil {
		return "", err
	}

	return p.GetString("liferay.workspace.product", ""), nil
}

func GetProductVersionAsTag(verbose bool) (string, error) {
	// dxp-7.4-u42
	version, err := GetProductVersion()

	if err != nil {
		return "", err
	}

	release, err := releases.GetReleaseByVersion(version, verbose)

	if err != nil {
		return "", err
	}

	// 7.4.3.u44
	targetPlatformVersion := release.(map[string]interface{})["targetPlatformVersion"].(string)

	parts := strings.Split(targetPlatformVersion, ".")

	if len(parts) != 4 {
		return "", errors.New(fmt.Sprintf("version '%s' doesn't have correct syntax", version))
	}

	lastPart := parts[3][1:]

	// 7.4.3.44-ga44
	// TODO this might not be safe... double check the algorithm
	return fmt.Sprintf("%s.%s.%s.%s-ga%s", parts[0], parts[1], parts[2][1:], lastPart, lastPart), nil
}
