package src

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/harmony-one/harmony-tui/config"
	"github.com/harmony-one/harmony-tui/data"
)

var previousJSONString = ""

func TailZeroLogFile() {
	fname, err := GetLogFilePath("zerolog")

	if err != nil {
		return
	}

	ticker := time.NewTicker(config.SystemStatsInterval)
	defer ticker.Stop()

	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	buf := make([]byte, 20480)

	for {
		select {
		case <-ticker.C:
			stat, err := os.Stat(fname)
			start := stat.Size() - 20480
			_, err = file.ReadAt(buf, start)
			if err == nil {
				jsonArray := strings.Split(fmt.Sprintf("%s\n", buf), "{\"level\":")
				for i := 0; i < len(jsonArray); i++ {
					if strings.Contains(jsonArray[i], "Signers") {
						var temp map[string]interface{}
						json.Unmarshal([]byte("{\"level\":"+jsonArray[i]), &temp)
						data.BlockData = temp
					}
					if strings.Contains(jsonArray[i], "\"message\":\"[") {
						var temp map[string]interface{}
						json.Unmarshal([]byte("{\"level\":"+jsonArray[i]), &temp)
						if temp == nil {
							continue
						}
						message := temp["message"].(string)
						if temp["time"] != nil {

							time := temp["time"].(string)
							switch {
							case strings.Contains(message, "[OnAnnounce]"):
								data.OnAnnounce = time
							case strings.Contains(message, "[Announce]"):
								data.Announce = time
							case strings.Contains(message, "[OnPrepared]"):
								data.OnPrepared = time
							case strings.Contains(message, "[Block Reward]"):
								data.BlockReward = time
							case strings.Contains(message, "[OnCommitted]"):
								data.OnCommitted = time
							}
						}
					}
				}
			} else {
				panic(err)
			}
		}
	}
}

func GetLogFilePath(prefix string) (string, error) {
	root := config.LogPath
	lastModified := time.Time{}
	var file string
	check, err := exists(root)

	if err != nil {
		panic(err)
	}

	if !check {
		return "", errors.New("Not Exists")
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), prefix) && strings.HasSuffix(info.Name(), ".log") {
			if lastModified.Before(info.ModTime()) {
				file = path
				lastModified = info.ModTime()
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
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
