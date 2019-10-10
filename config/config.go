package config

var (
	LogPath = "./"
	HmyURL = "http://localhost:9500/"
	HarmonyPath = "./"
)

func SetConfig(env string) {
	if env=="local" {
		LogPath = "/Users/manish/go/src/github.com/harmony-one/harmony/tmp_log/"
		HarmonyPath = "/Users/manish/go/src/github.com/harmony-one/harmony/bin/"
	} else if env=="ec2" {
		LogPath = "./latest/"
		HarmonyPath = "./"
	}
}