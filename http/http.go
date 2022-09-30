package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

type GetOrFetchBytesOptions struct {
	// The key used to store the etag in configuration
	EtagKey string
	// The key used to store the file path in configuration
	FileKey string
	// The actual URL of the file to retrieve
	URL string
	// Set the verbosity of output
	Verbose bool
}

func GetOrFetchBytes(options GetOrFetchBytesOptions) ([]byte, error) {
	req, err := http.NewRequest("GET", options.URL, nil)

	releasesEtag := viper.GetString(options.EtagKey)

	if releasesEtag != "" {
		if options.Verbose {
			fmt.Printf("Using stored Etag for %s=%s\n", options.EtagKey, releasesEtag)
		}
		req.Header.Add("If-None-Match", releasesEtag)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode == http.StatusNotModified {
		f := viper.GetString(options.FileKey)
		if options.Verbose {
			fmt.Printf("Not mofified use local %s=%s\n", options.FileKey, f)
		}
		return os.ReadFile(f)
	} else {
		releasesEtag = resp.Header.Get("ETag")
		viper.Set(options.EtagKey, releasesEtag)
		viper.WriteConfig()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, os.WriteFile(
			viper.GetString(options.FileKey), body, 0644)
	}
}
