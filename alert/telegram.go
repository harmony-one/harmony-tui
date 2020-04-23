package alert

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/harmony-one/harmony-tui/config"
	"github.com/harmony-one/harmony-tui/data"
	"github.com/harmony-one/harmony-tui/widgets"
	"github.com/spf13/viper"
)

var (
	bot           *tgbotapi.BotAPI
	lastAlertTime = time.Time{}
	chatId        int64
	lastAlert     string
)

func StartTelegramAlerts() {
	chatId = viper.GetInt64("TelegramChatId")
	b, err := tgbotapi.NewBotAPI(viper.GetString("TelegramToken"))
	if err != nil || viper.GetString("TelegramToken") == "" {
		return
	}
	bot = b
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if viper.GetInt64("TelegramChatId") == 0 {
			setChatId(update.Message.Chat.ID)
		}

		if update.Message.Chat.ID != viper.GetInt64("TelegramChatId") {
			continue
		}

		switch text := strings.ToLower(update.Message.Text); {
		case strings.Contains(text, "bingo"):
			SendTelegramMessage("Last Bingo timestamp = " + data.Bingo)
		case strings.Contains(text, "block"):
			res, _ := data.GetLatestHeader()
			blockNo, _ := res["result"].(map[string]interface{})["blockNumber"].(float64)
			SendTelegramMessage(fmt.Sprintf("Current BlockNo : %0.f", blockNo))
		case strings.Contains(text, "epoch"):
			res, _ := data.GetLatestHeader()
			epoch, _ := res["result"].(map[string]interface{})["epoch"].(float64)
			SendTelegramMessage(fmt.Sprintf("Epoch : %0.f", epoch))
		case strings.Contains(text, "balance"):
			balance, _ := data.GetBalance()
			SendTelegramMessage(balance)
		case strings.Contains(text, "version"):
			SendTelegramMessage(data.Metadata.Version)
		case strings.Contains(text, "shard"):
			res, _ := data.GetLatestHeader()
			shardID, _ := res["result"].(map[string]interface{})["shardID"].(float64)
			SendTelegramMessage(fmt.Sprintf("Shard Id : %0.f", shardID))
		case strings.Contains(text, "system"):
			SendTelegramMessage(fmt.Sprintf("CPU : %d%% \nMemory : %d%% \nDisk : %d%%", widgets.CpuUsage(), widgets.MemoryUsage(), widgets.DiskUsage()))
		case strings.Contains(text, "cpu"):
			SendTelegramMessage(fmt.Sprintf("CPU : %d%%", widgets.CpuUsage(), widgets.MemoryUsage(), widgets.DiskUsage()))
		case strings.Contains(text, "memory"):
			SendTelegramMessage(fmt.Sprintf("Memory : %d%%", widgets.CpuUsage(), widgets.MemoryUsage(), widgets.DiskUsage()))
		case strings.Contains(text, "disk"):
			SendTelegramMessage(fmt.Sprintf("Disk : %d%%", widgets.CpuUsage(), widgets.MemoryUsage(), widgets.DiskUsage()))
		case strings.Contains(text, "help"):
			SendTelegramMessage(`
			Below information can be requested from bot:
			[command] - [description]
			bingo     - Last bingo timestamp
			block     - Current block number
			version   - Version of harmony binary
			shard     - ShardId
			balance   - Balance in one account
			system    - Get system stats
			`)
		default:
			SendTelegramMessage("Command not recognized")
		}
	}
}

func SendTelegramMessage(message string) {

	lastAlertTime = time.Now()
	msg := tgbotapi.NewMessage(viper.GetInt64("TelegramChatId"), message)
	if bot != nil {
		bot.Send(msg)
	}
}

func SendTelegramAlert(alert string) {
	if lastAlert != alert {
		SendTelegramMessage(alert)
		lastAlert = alert
	}
}

func setChatId(id int64) {
	if chatId != id {
		chatId = id
		viper.Set("TelegramChatId", chatId)
		config.WriteConfig()
	}
}
