// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package browscap_go

import (
	"reflect"
	"strings"
	"sync"
)

type Browser struct {
	parent     string //name of the parent
	built      bool   // has complete data from parents
	buildMutex sync.Mutex

	Browser         string
	BrowserVersion  string
	BrowserMajorVer string
	BrowserMinorVer string
	BrowserBits     string
	BrowserMaker    string
	// Browser, Application, Bot/Crawler, Useragent Anonymizer, Offline Browser,
	// Multimedia Player, Library, Feed Reader, Email Client or unknown
	BrowserType string

	Platform        string
	PlatformShort   string
	PlatformVersion string
	PlatformBits    string
	PlatformMaker   string

	// Mobile Phone, Mobile Device, Tablet, Desktop, TV Device, Console,
	// FonePad, Ebook Reader, Car Entertainment System or unknown
	DeviceType  string
	DeviceName  string
	DeviceBrand string

	Crawler string

	Cookies    string
	JavaScript string

	RenderingEngineName    string
	RenderingEngineVersion string
}

func (browser *Browser) build(browsers map[string]*Browser) {
	if browser.built {
		return
	}

	browser.buildMutex.Lock()
	defer browser.buildMutex.Unlock()
	// Check again after lock if another goroutine built the object while we were waiting
	if browser.built {
		return
	}

	n := reflect.ValueOf(*browser).NumField()

	current := reflect.ValueOf(browser)
	parent := browser.parent
	for parent != "" {
		b, ok := browsers[parent]
		if !ok {
			break
		}

		parentObj := reflect.ValueOf(b)
		for i := 0; i < n; i++ {
			cField := current.Elem().Field(i)
			if cField.String() != "" {
				continue
			}

			pField := parentObj.Elem().Field(i)
			if pField.String() == "" {
				continue
			}

			cField.SetString(pField.String())
		}

		parent = b.parent
	}

	browser.built = true
}

func (browser *Browser) setValue(key, item string) {
	if key == "Parent" {
		browser.parent = item
	} else if key == "Browser" {
		browser.Browser = item
	} else if key == "Browser_Bits" {
		browser.BrowserBits = item
	} else if key == "Browser_Maker" {
		browser.BrowserMaker = item
	} else if key == "Version" {
		browser.BrowserVersion = item
	} else if key == "MajorVer" {
		browser.BrowserMajorVer = item
	} else if key == "MinorVer" {
		browser.BrowserMinorVer = item
	} else if key == "Browser_Type" {
		browser.BrowserType = item
	} else if key == "JavaScript" {
		browser.JavaScript = item
	} else if key == "Cookies" {
		browser.Cookies = item
	} else if key == "Crawler" {
		browser.Crawler = item
	} else if key == "Platform" {
		browser.Platform = item
		browser.PlatformShort = strings.ToLower(item)

		if strings.HasPrefix(browser.PlatformShort, "win") {
			browser.PlatformShort = "win"
		} else if strings.HasPrefix(browser.PlatformShort, "mac") {
			browser.PlatformShort = "mac"
		}
	} else if key == "Platform_Version" {
		browser.PlatformVersion = item
	} else if key == "Platform_Bits" {
		browser.PlatformBits = item
	} else if key == "Platform_Maker" {
		browser.PlatformMaker = item
	} else if key == "RenderingEngine_Name" {
		browser.RenderingEngineName = item
	} else if key == "RenderingEngine_Version" {
		browser.RenderingEngineVersion = item
	} else if key == "Device_Type" {
		browser.DeviceType = item
	} else if key == "Device_Code_Name" {
		browser.DeviceName = item
	} else if key == "Device_Brand_Name" {
		browser.DeviceBrand = item
	}
}

func (browser *Browser) IsCrawler() bool {
	return browser.BrowserType == "Bot/Crawler" || browser.Crawler == "true"
}

func (browser *Browser) IsMobile() bool {
	return browser.DeviceType == "Mobile Phone" || browser.DeviceType == "Mobile Device"
}

func (browser *Browser) IsTablet() bool {
	return browser.DeviceType == "Tablet" || browser.DeviceType == "FonePad" || browser.DeviceType == "Ebook Reader"
}

func (browser *Browser) IsDesktop() bool {
	return browser.DeviceType == "Desktop"
}

func (browser *Browser) IsConsole() bool {
	return browser.DeviceType == "Console"
}

func (browser *Browser) IsTv() bool {
	return browser.DeviceType == "TV Device"
}

func (browser *Browser) IsAndroid() bool {
	return browser.Platform == "Android"
}

func (browser *Browser) IsIPhone() bool {
	return browser.Platform == "iOS" && browser.DeviceName == "iPhone"
}

func (browser *Browser) IsIPad() bool {
	return browser.Platform == "iOS" && browser.DeviceName == "iPad"
}

func (browser *Browser) IsWinPhone() bool {
	return strings.Index(browser.Platform, "WinPhone") != -1 || browser.Platform == "WinMobile"
}
