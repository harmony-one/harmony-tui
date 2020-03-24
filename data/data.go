package data

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"

	"github.com/spf13/viper"
)

var (
	BlockData   map[string]interface{}
	VersionData map[string]interface{}
	Announce    string
	OnAnnounce  string
	OnPrepared  string
	BlockReward string
	Bingo       string
	OnCommitted string

	BlockHash       string
	BlockNumber     float64
	ShardID         float64
	Leader          string
	ViewID          float64
	Epoch           float64
	SizeInt         int64
	NoOfTransaction int
	StateRoot       string
	PeerCount       int64
	Balance         string
	TotalBalance    float64
	AppVersion      string
	EarningRate     float64

	Quitter func(string)
)

func RefreshData() {
	ticker := time.NewTicker(viper.GetDuration("RPCRefreshInterval"))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			latestHeader, err := GetLatestHeader()
			if err != nil {
				return
			}
			BlockHash, _ = latestHeader["result"].(map[string]interface{})["blockHash"].(string)
			BlockNumber, _ = latestHeader["result"].(map[string]interface{})["blockNumber"].(float64)
			ShardID, _ = latestHeader["result"].(map[string]interface{})["shardID"].(float64)
			Leader, _ = latestHeader["result"].(map[string]interface{})["leader"].(string)
			ViewID, _ = latestHeader["result"].(map[string]interface{})["viewID"].(float64)
			Epoch, _ = latestHeader["result"].(map[string]interface{})["epoch"].(float64)
			hexaBlockNumber := numToHex(BlockNumber)

			peerCountRply, err := rpc.Request(rpc.Method.PeerCount, viper.GetString("HmyURL"), []interface{}{})
			if err != nil {
				panic(err)
			}
			tempPeerCount, _ := peerCountRply["result"].(string)
			PeerCount = hexToNum(tempPeerCount)
			latestBlock, err := rpc.Request(rpc.Method.GetBlockByNumber, viper.GetString("HmyURL"), []interface{}{hexaBlockNumber, true})
			if err != nil {
				panic(err)
			}
			size, _ := latestBlock["result"].(map[string]interface{})["size"].(string)
			SizeInt = hexToNum(size)
			temp, _ := latestBlock["result"].(map[string]interface{})["transactions"].([]string)
			NoOfTransaction = len(temp)
			StateRoot, _ = latestBlock["result"].(map[string]interface{})["stateRoot"].(string)
			Balance, TotalBalance = GetBalance()
		}
	}
}

func hexToNum(hex string) int64 {
	rval, _ := strconv.ParseInt(hex[2:], 16, 32)
	return rval
}

func numToHex(num float64) string {
	return "0x" + strconv.FormatInt(int64(num), 16)
}

func GetLatestHeader() (map[string]interface{}, error) {
	return rpc.Request("hmy_latestHeader", viper.GetString("HmyURL"), []interface{}{})
}

func GetBalance() (string, float64) {
	tempBal := 0.00
	balance, err := sharding.CheckAllShards(viper.GetString("HmyURL"), viper.GetString("OneAddress"), true)
	if err != nil {
		balance = "No data"
	} else {
		var temp []map[string]interface{}
		err := json.Unmarshal([]byte(balance), &temp)
		if err != nil {
			panic(err)
		}
		balance = "Address: " + viper.GetString("OneAddress")

		for _, b := range temp {
			balance += "\n Balance in Shard " + strconv.FormatFloat(b["shard"].(float64), 'f', 0, 64) + ":  " + strconv.FormatFloat(b["amount"].(float64), 'f', 4, 64)
			tempBal += b["amount"].(float64)
		}
	}
	return balance, tempBal
}
