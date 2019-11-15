package alert

import (
	"fmt"
	"time"

	"github.com/harmony-one/harmony-tui/data"
	"github.com/harmony-one/harmony-tui/widgets"
	"github.com/spf13/viper"
)

var (
	nodeSync = true
)

func StartAlerting() {
	go StartTelegramAlerts()
	ticker := time.NewTicker(viper.GetDuration("AlertCheckIntervalInMin"))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkBingoAlert()
			checkDiskFullAlert()
		}
	}
}

func checkBingoAlert() {
	t, parseErr := time.Parse(viper.GetString("TimestampLayout"), data.Bingo)

	if parseErr == nil {
		if time.Since(t).Seconds() > viper.GetDuration("AlertCheckIntervalInMin").Seconds() {
			nodeSync = false
			SendTelegramAlert("=== Alert ===\nHarmony node out of sync\n LastBingo : " + data.Bingo + "\nOneAddress : " + viper.GetString("OneAddress"))
		} else if nodeSync == false {
			nodeSync = true
			SendTelegramAlert("Node recovered\n LastBingo : " + data.Bingo + "\nOneAddress : " + viper.GetString("OneAddress"))
		}
	}
}

func checkDiskFullAlert() {
	usage := widgets.DiskUsage()
	if usage > viper.GetInt("DiskSpaceAlertPerecent") {
		SendTelegramAlert(fmt.Sprintf("=== Alert ===\nDisk space almost full\nDisk space used %d%", usage))
	}
}
