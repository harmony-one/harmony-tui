package data

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/harmony-one/harmony-tui/rpc"
	"github.com/harmony-one/harmony-tui/config"
)

var BlockData map[string]interface{}
var VersionData map[string]interface{}
var Announce string = ""
var OnAnnounce string = ""
var OnPrepared string = ""
var BlockReward string = ""
var Bingo string = ""
var OnCommitted string = ""

var BlockHash string
var BlockNumber float64
var ShardID float64
var Leader string
var ViewID float64
var Epoch float64
var SizeInt int64
var NoOfTransaction int
var StateRoot string
var PeerCount int64
var OneAddress string
var Balance string
var AppVersion string

var Quitter func(string)

func init() {
	go refreshData()
}

func refreshData() {

	ticker := time.NewTicker(config.BlockchainInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			latestHeader, err := rpc.Request("hmy_latestHeader", config.HmyURL, []interface{}{})
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

			peerCountRply, err := rpc.Request(rpc.Method.PeerCount, config.HmyURL, []interface{}{})
			if err != nil {
				panic(err)
			}
			tempPeerCount, _ := peerCountRply["result"].(string)
			PeerCount = hexToNum(tempPeerCount)
			latestBlock, err := rpc.Request(rpc.Method.GetBlockByNumber, config.HmyURL, []interface{}{hexaBlockNumber, true})
			if err != nil {
				panic(err)
			}
			size, _ := latestBlock["result"].(map[string]interface{})["size"].(string)
			SizeInt = hexToNum(size)
			temp, _ := latestBlock["result"].(map[string]interface{})["transactions"].([]string)
			NoOfTransaction = len(temp)
			StateRoot, _ = latestBlock["result"].(map[string]interface{})["stateRoot"].(string)
			Balance, err = CheckAllShards(config.HmyURL, OneAddress, true)
			if err != nil {
				Balance = "No data"
			} else {
				var temp []map[string]interface{}
				err := json.Unmarshal([]byte(Balance), &temp)
				if err != nil {
					panic(err)
				}
				Balance = "Address: " + OneAddress
				for _, b := range temp {
					Balance += "\n Balance in Shard " + strconv.FormatFloat(b["shard"].(float64), 'f', 0, 64) + ":  " + strconv.FormatFloat(b["amount"].(float64), 'f', 4, 64)
				}
			}
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

func SetOneAddress(addr string) {
   OneAddress = addr
}
