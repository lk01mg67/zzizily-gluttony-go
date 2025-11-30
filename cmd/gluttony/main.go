package main

import (
	"fmt"
	"os"
	"time"

	"github.com/deuxksy/zzizily-gluttony-go/internal/configuration"
	"github.com/deuxksy/zzizily-gluttony-go/internal/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/input"
	"github.com/spf13/viper"
)

var yyMMddHHmm = time.Now().Local().Format("0601021504")

func init() {
	profile := initProfile()
	setRuntimeConfig(profile)
}

func setRuntimeConfig(profile string) {
	viper.AddConfigPath("configs")
	viper.SetConfigName(profile)
	viper.SetConfigType("yaml")
	viper.Set("Verbose", true)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&configuration.RuntimeConf)
	if err != nil {
		panic(err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Warn("Config file changed: %s", e.Name)
		var err error
		err = viper.ReadInConfig()
		if err != nil {
			logger.Error(err.Error())
			return
		}
		err = viper.Unmarshal(&configuration.RuntimeConf)
		if err != nil {
			logger.Error(err.Error())
			return
		}
	})
	viper.WatchConfig()
}

func initProfile() string {
	var profile string
	profile = os.Getenv("GO_PROFILE")
	if len(profile) <= 0 {
		profile = "local"
	}
	return profile
}

func initRod() *rod.Browser {
	browser := rod.New().MustConnect()
	browser.DefaultDevice(devices.IPadMini)
	return browser
}

func runLogin(browser *rod.Browser, scene configuration.Scene, idx int) *rod.Page {
	page := browser.MustPage(scene.Url)
	logger.Debug(page.MustInfo().URL)

	userID := os.Getenv("USERID")
	userPW := os.Getenv("USERPW")

	if userID == "" || userPW == "" {
		logger.Error("USERID or USERPW environment variables are not set")
		os.Exit(1)
	}

	page.MustWaitLoad().MustElement("input[name=userId]").MustWaitVisible().MustInput(userID)
	page.MustElement("input[name=userPwd]").MustWaitVisible().MustInput(userPW)

	page.MustScreenshotFullPage(fmt.Sprintf("screenshot/%s/%s-%02d.png", yyMMddHHmm, scene.Name, idx))
	page.MustElement("input[name=userPwd]").MustType(input.Enter)
	return page
}

func runBooking(page *rod.Page, scene configuration.Scene, idx int) {
	page.MustWaitLoad().MustNavigate(scene.Url)
	logger.Debug(page.MustInfo().URL)
	page.MustWaitLoad().MustScreenshotFullPage(fmt.Sprintf("screenshot/%s/%s-%02d-1.png", yyMMddHHmm, scene.Name, idx))

	if page.MustWaitLoad().MustHas(".bg-color-blue") {
		page.MustWaitLoad().MustScreenshotFullPage(fmt.Sprintf("screenshot/%s/%s-%02d-2.png", yyMMddHHmm, scene.Name, idx))
		elements := page.MustElements(`div[class="fc-event fc-event-hori fc-event-start fc-event-end bg-color-blue"]`)
		elements.Last().MustClick()
		page.MustWaitLoad().MustElement(`a[class="btn btn-info btn-sm"]`).MustClick()
		logger.Info("Complete %s", scene.Name)
	} else {
		logger.Warn("Not Found %s", scene.Name)
	}
	page.MustScreenshotFullPage(fmt.Sprintf("screenshot/%s/%s-%02d-3.png", yyMMddHHmm, scene.Name, idx))
}

func main() {
	logger.Info("%s", "Gluttony")
	browser := initRod()
	defer browser.MustClose()

	var page *rod.Page

	for i, scene := range configuration.RuntimeConf.Scenario {
		logger.Info("Running scene: %s (Type: %s)", scene.Name, scene.Type)
		switch scene.Type {
		case "login":
			page = runLogin(browser, scene, i)
			// Disable alerts after login
			page.Eval(`window.alert = () => {}`)
		case "booking":
			if page == nil {
				logger.Error("Cannot run booking scene without a page context (Login first)")
				continue
			}
			runBooking(page, scene, i)
		default:
			logger.Warn("Unknown scene type: %s", scene.Type)
		}
	}
}
