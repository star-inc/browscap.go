// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	TestIniFile             = "/tmp/browscap.ini"
	TestMacOsxChromeAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
	TestPixel6ChromeAgent   = "Mozilla/5.0 (Linux; Android 12; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.74 Mobile Safari/537.36"
	TestIPhoneSafariAgent   = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_2 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8H7 Safari/6533.18.5"
	TestIPhone12SafariAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 15_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Mobile/15E148 Safari/604.1"
	TestIPhone12ChromeAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 15_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/95.0.4638.50 Mobile/15E148 Safari/604.1"
)

func TestMain(m *testing.M) {
	if _, err := os.Stat(TestIniFile); !errors.Is(err, os.ErrNotExist) {
		if err := InitBrowsCap(TestIniFile, false); err != nil {
			log.Panicln(err)
		}
	}

	currentVersion := InitializedVersion()
	latestVersion, err := LatestVersion()
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Curerent Version: %s, Latest Version: %s", currentVersion, latestVersion)
	if currentVersion != latestVersion {
		log.Printf("Start downloading version %q\n", currentVersion)
		if err = DownloadFile(TestIniFile); err != nil {
			log.Fatalf("%v\n", err)
		}
		log.Printf("Initializing with new version")
		if err = InitBrowsCap(TestIniFile, true); err != nil {
			log.Fatalln(err)
		}
		if version != InitializedVersion() {
			log.Fatalln("New file is wrong")
		}
	}

	os.Exit(m.Run())
}

func TestInitBrowsCap(t *testing.T) {
	if err := InitBrowsCap(TestIniFile, false); err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetBrowserOnMacOsxChrome(t *testing.T) {
	if browser, ok := GetBrowser(TestMacOsxChromeAgent); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "Chrome" {
		t.Errorf("Expected Chrome but got %q", browser.Browser)
	} else if browser.Platform != "MacOSX" {
		t.Errorf("Expected MacOSX but got %q", browser.Platform)
	} else if browser.BrowserVersion != "37.0" {
		t.Errorf("Expected 37.0 but got %q", browser.BrowserVersion)
	} else if browser.RenderingEngineName != "Blink" {
		t.Errorf("Expected Blink but got %q", browser.RenderingEngineName)
	} else if browser.Crawler != "false" {
		t.Errorf("Expected false but got %q", browser.Crawler)
	}
}

func TestGetBrowserOnPixel6Chrome(t *testing.T) {
	if _, ok := GetBrowser(TestPixel6ChromeAgent); !ok {
		t.Error("Browser not found")
	}
}

func TestGetBrowserOnIPhoneSafari(t *testing.T) {
	if browser, ok := GetBrowser(TestIPhoneSafariAgent); !ok {
		t.Error("Browser not found")
	} else if browser.DeviceName != "iPhone" {
		t.Errorf("Expected iPhone but got %q", browser.DeviceName)
	} else if browser.Platform != "iOS" {
		t.Errorf("Expected iOS but got %q", browser.Platform)
	} else if browser.PlatformVersion != "4.3" {
		t.Errorf("Expected 4.3 but got %q", browser.PlatformVersion)
	} else if browser.IsMobile() != true {
		t.Errorf("Expected true but got %t", browser.IsMobile())
	}
}

func TestGetBrowserOnIPhone12Safari(t *testing.T) {
	if browser, ok := GetBrowser(TestIPhone12SafariAgent); !ok {
		t.Error("Browser not found")
	} else if browser.DeviceName != "iPhone" {
		t.Errorf("Expected iPhone but got %q", browser.DeviceName)
	} else if browser.Platform != "iOS" {
		t.Errorf("Expected iOS but got %q", browser.Platform)
	} else if browser.PlatformVersion != "15.1" {
		t.Errorf("Expected 15.1 but got %q", browser.PlatformVersion)
	} else if browser.IsMobile() != true {
		t.Errorf("Expected true but got %t", browser.IsMobile())
	}
}

func TestGetBrowserOnIPhone12Chrome(t *testing.T) {
	if browser, ok := GetBrowser(TestIPhone12ChromeAgent); !ok {
		t.Error("Browser not found")
	} else if browser.DeviceName != "iPhone" {
		t.Errorf("Expected iPhone but got %q", browser.DeviceName)
	} else if browser.Platform != "iOS" {
		t.Errorf("Expected iOS but got %q", browser.Platform)
	} else if browser.PlatformVersion != "15.1" {
		t.Errorf("Expected 15.1 but got %q", browser.PlatformVersion)
	} else if browser.IsMobile() != true {
		t.Errorf("Expected true but got %t", browser.IsMobile())
	}
}

func TestGetBrowserYandex(t *testing.T) {
	if browser, ok := GetBrowser("Yandex Browser 1.1"); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "Yandex Browser" {
		t.Errorf("Expected Yandex Browser but got %q", browser.Browser)
	} else if browser.IsCrawler() != false {
		t.Errorf("Expected false but got %t", browser.IsCrawler())
	}
}

func TestGetBrowser360Spider(t *testing.T) {
	if browser, ok := GetBrowser("360Spider"); !ok {
		t.Error("Browser not found")
	} else if browser.Browser != "360Spider" {
		t.Errorf("Expected Chrome but got %q", browser.Browser)
	} else if browser.IsCrawler() != true {
		t.Errorf("Expected true but got %t", browser.IsCrawler())
	}
}

func TestGetBrowserIssues(t *testing.T) {
	// https://github.com/digitalcrab/browscap_go/issues/4
	ua := "Mozilla/5.0 (iPad; CPU OS 5_0_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A405 Safari/7534.48.3"
	if browser, ok := GetBrowser(ua); !ok {
		t.Error("Browser not found")
	} else if browser.DeviceType != "Tablet" {
		t.Errorf("Expected tablet %q", browser.DeviceType)
	}
}

func TestLatestVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	version, err := LatestVersion()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if version == "" {
		t.Fatalf("Version not found")
	}
	t.Logf("Latest version is %q, current version: %q", version, InitializedVersion())
}

func TestDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	version, err := LatestVersion()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if version != InitializedVersion() {
		t.Logf("Start downloading version %q", version)
		if err = DownloadFile(TestIniFile); err != nil {
			t.Fatalf("%v", err)
		}

		t.Logf("Initializing with new version")
		if err = InitBrowsCap(TestIniFile, true); err != nil {
			t.Error(err)
		}

		if version != InitializedVersion() {
			t.Fatalf("New file is wrong")
		}
	}
}

func BenchmarkInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := InitBrowsCap(TestIniFile, true)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGetBrowser(b *testing.B) {
	data, err := ioutil.ReadFile("test-data/user_agents_sample.txt")
	if err != nil {
		b.Error(err)
	}

	uas := strings.Split(strings.TrimSpace(string(data)), "\n")

	if err := InitBrowsCap(TestIniFile, false); err != nil {
		b.Fatalf("%v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx := i % len(uas)

		_, ok := GetBrowser(uas[idx])
		if !ok {
			b.Errorf("User agent not recognized: %s", uas[idx])
		}
	}
}
