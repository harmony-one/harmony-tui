package config

import(
	"time"
)

var (
	LogPath = "./"
	HmyURL = "http://localhost:9500/"
	HarmonyPath = "./"
	BlockchainInterval = 3000*time.Millisecond
	SystemStatsInterval = 500*time.Millisecond
)

func SetConfig(env string) {
	if env=="local" {
		LogPath = "/Users/manish/go/src/github.com/harmony-one/harmony/tmp_log/"
		HarmonyPath = "/Users/manish/go/src/github.com/harmony-one/harmony/bin/"
		BlockchainInterval = 3000*time.Millisecond
		SystemStatsInterval = 250*time.Millisecond
	} else if env=="ec2" {
		LogPath = "./latest/"
		HarmonyPath = "./"
		BlockchainInterval = 5000*time.Millisecond
		SystemStatsInterval = 500*time.Millisecond
	}
}