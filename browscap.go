// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"unicode"
)

const (
	DownloadUrl     = "https://browscap.org/stream?q=Full_PHP_BrowsCapINI"
	CheckVersionUrl = "https://browscap.org/version-number"
)

var (
	dict        *dictionary
	initialized bool
	version     string
)

func InitBrowsCap(path string, force bool) error {
	if initialized && !force {
		return nil
	}
	var err error

	// Load ini file
	if dict, err = loadFromIniFile(path); err != nil {
		return fmt.Errorf("browscap: An error occurred while reading file, %v ", err)
	}

	initialized = true
	return nil
}

func InitializedVersion() string {
	return version
}

func LatestVersion() (string, error) {
	response, err := http.Get(CheckVersionUrl)
	if err != nil {
		return "", fmt.Errorf("browscap: error sending request, %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Panicln(err)
		}
	}(response.Body)

	// Get body of response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("browscap: error reading the response data of request, %v", err)
	}

	// Check 200 status
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("browscap: error unexpected status code %d", response.StatusCode)
	}

	return string(body), nil
}

func DownloadFile(saveAs string) error {
	response, err := http.Get(DownloadUrl)
	if err != nil {
		return fmt.Errorf("browscap: error sending request, %v", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Panicln(err)
		}
	}(response.Body)

	// Get body of response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("browscap: error reading the response data of request, %v", err)
	}

	// Check 200 status
	if response.StatusCode != 200 {
		return fmt.Errorf("browscap: error unexpected status code %d", response.StatusCode)
	}

	if err = ioutil.WriteFile(saveAs, body, os.ModePerm); err != nil {
		return fmt.Errorf("browscap: error saving file, %v", err)
	}

	return nil
}

func GetBrowser(userAgent string) (browser *Browser, ok bool) {
	if !initialized {
		return
	}

	agent := mapToBytes(unicode.ToLower, userAgent)
	defer bytesPool.Put(agent)

	name := dict.tree.Find(agent)
	if name == "" {
		return
	}

	browser = dict.getBrowser(name)
	if browser != nil {
		ok = true
	}

	return
}
