package src

import (
	"time"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"encoding/json"

	"text-based-ui/src/data"
	"text-based-ui/config"
)

var previousJSONString = ""

func TailZeroLogFile() {
	fname := GetLogFilePath("zerolog")
	
	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop()
	
	file, err := os.Open(fname)
	if err!=nil {
		panic(err)
	}
	defer file.Close()
	buf := make([]byte, 10240)
	
	for {
		select {
		case <- ticker.C:
			stat, err := os.Stat(fname)
			start := stat.Size() - 10240
			_, err = file.ReadAt(buf, start)
			if err == nil {
				jsonArray := strings.Split(fmt.Sprintf("%s\n", buf), "{\"level\":")
				foundSigner := false
				for i:=0; i<len(jsonArray); i++ {
					if !foundSigner && strings.Contains(jsonArray[i], "Signers") {
						json.Unmarshal([]byte("{\"level\":" + jsonArray[i]),&data.BlockData)
						foundSigner = false
					}
					if strings.Contains(jsonArray[i], "\"message\":\"[") {
						var temp map[string]interface{}
						json.Unmarshal([]byte("{\"level\":" + jsonArray[i]),&temp)
						if temp == nil {
							continue
						}
						message := temp["message"].(string)
						if temp["time"]!=nil {

							time := temp["time"].(string)
							switch {
							case strings.Contains(message, "[OnAnnounce]") :
								data.OnAnnounce = time
							case strings.Contains(message, "[Announce]") :
								data.Announce = time
							case strings.Contains(message, "[OnPrepared]") :
								data.OnPrepared = time
							case strings.Contains(message, "[Block Reward]") :
								data.BlockReward = time
							case strings.Contains(message, "BINGO") :
								data.Bingo = time
							case strings.Contains(message, "[OnCommitted]") :
								data.OnCommitted = time
							}
						}
					}
				}
			}else{panic(err)}
		}
	}
}

func GetLogFilePath(prefix string) string {
	root := config.LogPath
	lastModified := time.Time{}
	var file string
	
	err := filepath.Walk(root, func(path string, info os.FileInfo,err error) error {
		
		if (strings.HasPrefix(info.Name(), prefix) && strings.HasSuffix(info.Name(),".log")){
			if lastModified.Before(info.ModTime()) {
				file = path
				lastModified = info.ModTime()
			}
		}
		return nil
	})
	if err!=nil {
		panic(err)
	}
	return file
}