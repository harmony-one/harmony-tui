package data

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"

	"github.com/spf13/viper"
)

var (
	BeaconChainEndpoint string

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
	BlockTimestamp  string
	ShardID         float64
	Leader          string
	ViewID          float64
	Epoch           float64

	SizeInt         int64
	NoOfTransaction int
	NoOfStakingTransaction int
	StateRoot       string

	AppVersion      string
	NetworkID       string
	IsArchival      bool
	BLSKeys         []string

	PeerCount       int64

	Elected bool
	EposStatus string
	BootedStatus string
	TotalDelegation *big.Int
	LifetimeAvalibility string
	LifetimeRewards *big.Int
	//CollectableRewards uint64
	APR string
	ValidatorBLSKeys []string
	CurrentAvailibility string
	EffectiveStake string

	Balance         string
	TotalBalance    float64

	EarningRate     float64

	NumKeys int

	Quitter func(string)
)

func RefreshData() {
	// TODO: Get API endpoint using GetShardingStructure
	//shardStructureReply, err := rpc.Request(rpc.Method.GetShardingStructure, viper.GetString("HmyURL"), []interface{}{})
	//if err != nil {
  //	panic(err)
	//}
	//BeaconChainEndpoint, _ = shardStructureReply["results"].([]interface{})[0].(map[string]interface{})["http"].(string)

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
			BlockTimestamp, _ = latestHeader["result"].(map[string]interface{})["timestamp"].(string)
			ShardID, _ = latestHeader["result"].(map[string]interface{})["shardID"].(float64)
			Leader, _ = latestHeader["result"].(map[string]interface{})["leader"].(string)
			ViewID, _ = latestHeader["result"].(map[string]interface{})["viewID"].(float64)
			Epoch, _ = latestHeader["result"].(map[string]interface{})["epoch"].(float64)
			hexaBlockNumber := numToHex(BlockNumber)

			latestBlockReply, err := getBlockByNumber(hexaBlockNumber, false)
			if err != nil {
				panic(err)
			}
			size, _ := latestBlockReply["result"].(map[string]interface{})["size"].(string)
			SizeInt = hexToNum(size)
			transactions, _ := latestBlockReply["result"].(map[string]interface{})["transactions"].([]string)
			NoOfTransaction = len(transactions)
			staking, _ := latestBlockReply["result"].(map[string]interface{})["stakingTransactions"].([]string)
			NoOfStakingTransaction = len(staking)
			StateRoot, _ = latestBlockReply["result"].(map[string]interface{})["stateRoot"].(string)

			metadataReply, err := getNodeMetadata()
			if err != nil {
				panic(err)
			}
			AppVersion, _ = metadataReply["result"].(map[string]interface{})["version"].(string)
			NetworkID, _ = metadataReply["result"].(map[string]interface{})["network"].(string)
			IsArchival, _ = metadataReply["result"].(map[string]interface{})["is-archival"].(bool)
			BLSKeys, _ = metadataReply["result"].(map[string]interface{})["blskey"].([]string)
			NumKeys = len(BLSKeys)

			peerCountReply, err := getPeerCount()
			if err != nil {
				panic(err)
			}
			count, _ := peerCountReply["result"].(string)
			PeerCount = hexToNum(count)

			validatorInformationReply, err := getValidatorInformation()
			if err != nil {
				// TODO: This can fail because of the endpoint going down
				panic(err)
			}
			Elected, _ = validatorInformationReply["result"].(map[string]interface{})["currently-in-committee"].(bool)
			EposStatus, _ = validatorInformationReply["result"].(map[string]interface{})["epos-status"].(string)
			BootedStatus, _ = validatorInformationReply["result"].(map[string]interface{})["booted-status"].(string)
			tempDel, _ := validatorInformationReply["result"].(map[string]interface{})["total-delegation"].(int64)
			TotalDelegation = big.NewInt(tempDel)
			lifetimeBlocksSigned, _ := validatorInformationReply["result"].(map[string]interface{})["lifetime"].(map[string]interface{})["blocks"].(map[string]interface{})["signed"].(float64)
			lifetimeBlocksToSign, _ := validatorInformationReply["result"].(map[string]interface{})["lifetime"].(map[string]interface{})["blocks"].(map[string]interface{})["to-sign"].(float64)
			LifetimeAvalibility = fmt.Sprintf("%.2f%%", lifetimeBlocksSigned / lifetimeBlocksToSign)
			LifetimeRewards, _ = validatorInformationReply["result"].(map[string]interface{})["lifetime"].(map[string]interface{})["reward-accumulated"].(*big.Int)
			//delegations, _ := validatorInformationReply["result"].(map[string]interface{})["validator"].(map[string]interface{})["delegations"].([]interface{})
			//for d := range delegations {
			//	if d.(map[string]interface{})
			//}
			//CollectableRewards, _ =
			APR, _ = validatorInformationReply["result"].(map[string]interface{})["lifetime"].(map[string]interface{})["apr"].(string)
			ValidatorBLSKeys, _ = validatorInformationReply["result"].(map[string]interface{})["validator"].(map[string]interface{})["bls-public-keys"].([]string)
			if Elected {
				CurrentAvailibility, _ = validatorInformationReply["result"].(map[string]interface{})["current-epoch-performance"].(map[string]interface{})["current-epoch-signing-percent"].(map[string]interface{})["current-epoch-signing-percentage"].(string)
				EffectiveStake, _ = validatorInformationReply["result"].(map[string]interface{})["epos-winning-stake"].(string)
			} else {
				CurrentAvailibility = ""
				EffectiveStake = ""
			}

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
	return rpc.Request(rpc.Method.GetLatestBlockHeader, viper.GetString("HmyURL"), []interface{}{})
}

func getBlockByNumber(hexaBlockNumber string, getTransactions bool) (map[string]interface{}, error) {
	return rpc.Request(rpc.Method.GetBlockByNumber, viper.GetString("HmyURL"), []interface{}{hexaBlockNumber, getTransactions})
}

func getPeerCount() (map[string]interface{}, error) {
	return rpc.Request(rpc.Method.PeerCount, viper.GetString("HmyURL"), []interface{}{})
}

func getNodeMetadata() (map[string]interface{}, error) {
	return rpc.Request(rpc.Method.GetNodeMetadata, viper.GetString("HmyURL"), []interface{}{})
}

// Always query BeaconChainEndpoint to get latest validator information on chain
func getValidatorInformation() (map[string]interface{}, error) {
	// TODO: Remove temp code for testing
	return rpc.Request(rpc.Method.GetValidatorInformation, "https://api.s0.os.hmny.io", []interface{}{"one1c0w53749uf70lfzdehhl0t23qdjvha0sf2ug5r"})
	//return rpc.Request(rpc.Method.GetValidatorInformation, BeaconChainEndpoint, []interface{}{viper.GetString("OneAddress")})
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
