package src

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/harmony-one/harmony-tui/data"
	"github.com/spf13/viper"
)

var previousJSONString = ""

func TailZeroLogFile() {
	for range time.Tick(time.Second * 5) {
		readLogs()
	}
}

func GetLogFilePath(prefix string) (string, error) {
	root := viper.GetString("LogPath")
	lastModified := time.Time{}
	var file string
	if check, _ := exists(root); !check {
		return "", errors.New("log path does not exist")
	}

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), prefix) && strings.HasSuffix(info.Name(), ".log") {
			if lastModified.Before(info.ModTime()) {
				file = path
				lastModified = info.ModTime()
			}
		}
		return nil
	}); err != nil {
		return "", errors.New("log path does not exist")
	}
	return file, nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func readLogs() {
	fname, err := GetLogFilePath("zerolog")
	if err != nil {
		return
	}

	file, err := os.Open(fname)
	if err != nil {
		return
	}
	defer file.Close()

	buf := make([]byte, 20480)
	stat, err := os.Stat(fname)
	start := stat.Size() - 20480
	_, err = file.ReadAt(buf, start)
	s := strings.Split(string(buf), "{\"level\":\"")
	for _, line := range s {
		var temp map[string]interface{}
		json.Unmarshal([]byte("{\"level\":\""+line), &temp)

		if temp == nil {
			continue
		}

		if temp["time"] != nil && temp["message"] != nil {
			message := temp["message"].(string)
			time := temp["time"].(string)
			switch {
			case strings.Contains(message, "[OnAnnounce]"):
				data.OnAnnounce = time
			case strings.Contains(message, "[Announce]"):
				data.Announce = time
			case strings.Contains(message, "[OnPrepared]"):
				data.OnPrepared = time
			case strings.Contains(message, "[OnCommitted]"):
				data.OnCommitted = time
			case strings.Contains(message, "HOORAY") || strings.Contains(message, "BINGO"):
				data.Bingo = time
			}
		}
	}
}
