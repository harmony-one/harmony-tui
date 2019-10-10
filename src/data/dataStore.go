package data

import(
	"strconv"
	
	"text-based-ui/src/rpc"
	"text-based-ui/config"
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

func init() {
	RefreshData()
}

func RefreshData(){
	latestHeader, err := rpc.Request("hmy_latestHeader", config.HmyURL, []interface{}{})

	BlockHash, _ = latestHeader["result"].(map[string]interface{})["blockHash"].(string)
	BlockNumber, _ = latestHeader["result"].(map[string]interface{})["blockNumber"].(float64)
	ShardID, _ = latestHeader["result"].(map[string]interface{})["shardID"].(float64)
	Leader, _ = latestHeader["result"].(map[string]interface{})["leader"].(string)
	ViewID, _ = latestHeader["result"].(map[string]interface{})["viewID"].(float64)
	Epoch, _ = latestHeader["result"].(map[string]interface{})["epoch"].(float64)
	hexaBlockNumber := numToHex(BlockNumber)

	peerCountRply, err := rpc.Request(rpc.Method.PeerCount, config.HmyURL, []interface{}{})
	if err!=nil {
		panic(err)
	}
	tempPeerCount, _ := peerCountRply["result"].(string)
	PeerCount = hexToNum(tempPeerCount)
	latestBlock, err := rpc.Request(rpc.Method.GetBlockByNumber, config.HmyURL, []interface{}{hexaBlockNumber, true})
	if err!=nil {
		panic(err)
	}
	size, _ := latestBlock["result"].(map[string]interface{})["size"].(string)
	SizeInt = hexToNum(size)
	temp, _ := latestBlock["result"].(map[string]interface{})["transactions"].([]string)
	NoOfTransaction = len(temp)
	StateRoot, _ = latestBlock["result"].(map[string]interface{})["stateRoot"].(string)
}

func hexToNum(hex string) int64 {
	rval, _ := strconv.ParseInt(hex[2:], 16, 32)
	return rval
}

func numToHex(num float64) string {
	return "0x" + strconv.FormatInt(int64(num), 16)
}