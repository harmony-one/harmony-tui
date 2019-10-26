package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"time"

	"github.com/spf13/viper"
)

var (
	logPath                 = "./"
	hmyURL                  = "http://localhost:9500/"
	harmonyPath             = "./"
	rpcRefreshInterval      = 3 * time.Second
	systemStatsInterval     = 500 * time.Millisecond
	widgetInterval          = 1000 * time.Millisecond
	timestampEC2            = "2006-01-02T15:04:05.000000000Z"
	timestampMac            = "2006-01-02T15:04:05.000000-07:00"
	timestampLayout         = timestampEC2
	earningRateInterval     = 20 * time.Second
	outOfSyncTimeInMin      = 5.00
	oneAddress              string
	cfgFile                 string
	alertCheckIntervalInMin = 10 * time.Minute
	diskSpaceAlertPerecent  = 90
)

func init() {
	viper.SetDefault("LogPath", logPath)
	viper.SetDefault("HmyURL", hmyURL)
	viper.SetDefault("HarmonyPath", harmonyPath)
	viper.SetDefault("RPCRefreshInterval", rpcRefreshInterval)
	viper.SetDefault("SystemStatsInterval", systemStatsInterval)
	viper.SetDefault("WidgetInterval", widgetInterval)
	viper.SetDefault("TimestampLayout", timestampLayout)
	viper.SetDefault("EarningRateInterval", earningRateInterval)
	viper.SetDefault("OutOfSyncTimeInMin", outOfSyncTimeInMin)
	viper.SetDefault("OneAddress", oneAddress)
	viper.SetDefault("env", "ec2")
	viper.SetDefault("AlertCheckIntervalInMin", alertCheckIntervalInMin)
	viper.SetDefault("DiskSpaceAlertPerecent", diskSpaceAlertPerecent)
}

func SetConfig() {
	env := flag.String("env", "", "environment of system binary is running on option 1- \"local\" option 2- \"ec2\"")
	flag.StringVar(&oneAddress, "address", "Not Provided", "address of your one account")
	flag.StringVar(&cfgFile, "config", "", "path to configuration file")

	log := flag.String("logPath", "", "path to harmony log folder \"latest\"")
	url := flag.String("hmyUrl", "", "harmony instance url")
	binaryPath := flag.String("hmyPath", "", "path to harmony binary (default is current dir)")
	refreshInterval := flag.String("refreshInterval", "", "Refresh interval of TUI in seconds")
	earningInterval := flag.String("earningInterval", "", "Earning interval of TUI in seconds")
	telegramToken := flag.String("telegramToken", "", "telegram token of your telegram bot")
	flag.Parse()

	initConfig()

	if *env != "" {
		viper.Set("env", env)
	}

	if viper.GetString("env") == "local" {
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = build.Default.GOPATH
		}

		viper.Set("LogPath", gopath+"/src/github.com/harmony-one/harmony/tmp_log/")
		viper.Set("HarmonyPath", gopath+"/src/github.com/harmony-one/harmony/bin/")
		viper.Set("TimestampLayout", timestampMac)
	} else if viper.GetString("env") == "ec2" {
		viper.Set("LogPath", "./latest/")
		viper.Set("HarmonyPath", "./")
	}

	if *log != "" {
		viper.Set("LogPath", *log)
	}
	if *url != "" {
		viper.Set("HmyURL", url)
	}
	if *binaryPath != "" {
		viper.Set("HarmonyPath", *binaryPath)
	}
	if *refreshInterval != "" {
		interval, err := time.ParseDuration(*refreshInterval + "s")
		if err == nil {
			viper.Set("RPCRefreshInterval", interval)
		}
	}
	if *earningInterval != "" {
		interval, err := time.ParseDuration(*earningInterval + "s")
		if err != nil || interval.Seconds() < 10 {
			fmt.Println("Earning duration should be greater than 10 seconds")
			os.Exit(1)
		} else {
			viper.Set("EarningRateInterval", interval)
		}
	}

	if viper.GetString("OneAddress") == "" || oneAddress != "Not Provided" {
		viper.Set("OneAddress", oneAddress)
	}

	if viper.GetString("TelegramToken") == "" || *telegramToken != "" {
		viper.Set("TelegramToken", *telegramToken)
	}

	validateConfig()
	WriteConfig()
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "config-tui" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigName("config-tui")
	}

	viper.ReadInConfig()
}

func validateConfig() {

	if viper.GetDuration("rpcRefreshInterval").Seconds() < 1 {
		fmt.Println("rpcRefreshInterval duration should be between 1 and 20 seconds")
		os.Exit(1)
	}

	if viper.GetDuration("SystemStatsInterval").Seconds() < 0.4 {
		fmt.Println("SystemStatsInterval duration should be more than 400 milli seconds")
		os.Exit(1)
	}

	if viper.GetDuration("WidgetInterval").Seconds() < 1 {
		fmt.Println("WidgetInterval duration should be more than 1 seconds")
		os.Exit(1)
	}

	if viper.GetDuration("EarningRateInterval").Seconds() < 10 {
		fmt.Println("EarningRateInterval duration should be more than 10 seconds")
		os.Exit(1)
	}
}

func WriteConfig() {
	s, _ := json.Marshal(viper.AllSettings())
	ioutil.WriteFile("config-tui.json", s, os.ModePerm)
}
