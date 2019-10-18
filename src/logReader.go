package src

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/harmony-one/harmony-tui/config"
	"github.com/harmony-one/harmony-tui/data"
	"github.com/hpcloud/tail"
)

var previousJSONString = ""

func TailZeroLogFile() {
	fname, err := GetLogFilePath("zerolog")

	if err != nil {
		return
	}

	t, _ := tail.TailFile(fname, tail.Config{Follow: true, MustExist: false, Logger: log.New(ioutil.Discard, "", 0), Location: &tail.SeekInfo{Offset: 1, Whence: 2}})

	for line := range t.Lines {
		var temp map[string]interface{}
		json.Unmarshal([]byte(line.Text), &temp)

		if temp == nil {
			continue
		}

		if strings.Contains(line.Text, "Signers") {
			data.BlockData = temp
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
			case strings.Contains(message, "[Block Reward]"):
				data.BlockReward = time
			case strings.Contains(message, "[OnCommitted]"):
				data.OnCommitted = time
			case strings.Contains(message, "HOORAY") || strings.Contains(message, "BINGO"):
				data.Bingo = time
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
