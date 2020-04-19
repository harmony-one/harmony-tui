package data

import (
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/harmony/numeric"
	//"github.com/harmony-one/harmony/staking/types"

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

	LatestHeader  LatestHeaderReply
	LatestBlock   BlockByNumberReply
	Metadata      NodeMetadataReply
	ValidatorInfo ValidatorInformationReply
	LifetimeAvail string
	PeerCount     int64
	Balance       string
	TotalBalance  float64
	EarningRate   float64

	Quitter func(string)

	oneAddressPattern = regexp.MustCompile("one1[0-9a-z]+")
)

func RefreshData() {
	shardingReply, err := getShardingStructure()
	if err != nil {
  	panic(err)
	}
	for _, s := range shardingReply {
		if s.ShardID == uint32(0) {
			BeaconChainEndpoint = s.HTTP
		}
	}

	ticker := time.NewTicker(viper.GetDuration("RPCRefreshInterval"))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			latestHeader, err := getLatestHeader()
			if err != nil {
				return
			}
			LatestHeader = latestHeader
			hexaBlockNumber := numToHex(LatestHeader.BlockNumber)

			latestBlockReply, err := getBlockByNumber(hexaBlockNumber)
			if err != nil {
				panic(err)
			}
			LatestBlock = latestBlockReply
			LatestBlock.BlockSizeInt = int(hexToNum(LatestBlock.BlockSize))
			LatestBlock.NumTransactions = len(LatestBlock.Transactions)
			LatestBlock.NumStakingTransactions = len(LatestBlock.StakingTransactions)

			metadataReply, err := getNodeMetadata()
			if err != nil {
				panic(err)
			}
			Metadata = metadataReply

			peerCountReply, err := getPeerCount()
			if err != nil {
				panic(err)
			}
			PeerCount = hexToNum(peerCountReply)

			validatorReply, err := getValidatorInformation()
			// Possible to get bad response due to rate limiting on endpoint or endpoint going down
			if err == nil {
				ValidatorInfo = validatorReply
				lifetimeSigned := numeric.NewDecFromBigInt(ValidatorInfo.Lifetime.Signing.NumBlocksSigned)
				lifetimeToSign := numeric.NewDecFromBigInt(ValidatorInfo.Lifetime.Signing.NumBlocksToSign)
				LifetimeAvail = lifetimeSigned.Quo(lifetimeToSign).String()
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

func getLatestHeader() (LatestHeaderReply, error) {
	type reply struct {
		Result LatestHeaderReply `json:"result"`
	}

	r, err := rpc.RawRequest(rpc.Method.GetLatestBlockHeader, viper.GetString("HmyURL"), []interface{}{})
	if err != nil {
		return LatestHeaderReply{}, err
	}

	temp := reply{}
	err = json.Unmarshal(r, &temp)
	if err != nil {
		return LatestHeaderReply{}, err
	}
	return temp.Result, nil
}

func getBlockByNumber(hexaBlockNumber string) (BlockByNumberReply, error) {
	type reply struct {
		Result BlockByNumberReply `json:"result"`
	}

	r, err := rpc.RawRequest(rpc.Method.GetBlockByNumber, viper.GetString("HmyURL"), []interface{}{hexaBlockNumber, false})
	if err != nil {
		return BlockByNumberReply{}, err
	}

	temp := reply{}
	err = json.Unmarshal(r, &temp)
	if err != nil {
		return BlockByNumberReply{}, err
	}
	return temp.Result, nil
}

func getPeerCount() (string, error) {
	type reply struct {
		Result string `json:"result"`
	}

	r, err := rpc.RawRequest(rpc.Method.PeerCount, viper.GetString("HmyURL"), []interface{}{})
	if err != nil {
		return "", err
	}

	temp := reply{}
	err = json.Unmarshal(r, &temp)
	if err != nil {
		return "", err
	}
	return temp.Result, nil
}

func getShardingStructure() ([]StructureReply, error) {
	type reply struct {
		Result []StructureReply `json:"result"`
	}

	r, err := rpc.RawRequest(rpc.Method.GetShardingStructure, viper.GetString("HmyURL"), []interface{}{})
	if err != nil {
		return []StructureReply{}, err
	}

	temp := reply{}
	err = json.Unmarshal(r, &temp)
	if err != nil {
		return []StructureReply{}, err
	}
	return temp.Result, nil
}

func getNodeMetadata() (NodeMetadataReply, error) {
	type reply struct {
		Result NodeMetadataReply `json:"result"`
	}

	r, err := rpc.RawRequest(rpc.Method.GetNodeMetadata, viper.GetString("HmyURL"), []interface{}{})
	if err != nil {
		return NodeMetadataReply{}, err
	}

	temp := reply{}
	err = json.Unmarshal(r, &temp)
	if err != nil {
		return NodeMetadataReply{}, err
	}
	return temp.Result, nil
}

// Always query BeaconChainEndpoint to get latest validator information on chain
func getValidatorInformation() (ValidatorInformationReply, error) {
	type reply struct {
		Result ValidatorInformationReply `json:"result"`
	}

	// TODO: Remove temp code for testing
	//r, err := rpc.RawRequest(rpc.Method.GetValidatorInformation, "https://api.s0.os.hmny.io", []interface{}{"one1c0w53749uf70lfzdehhl0t23qdjvha0sf2ug5r"})
	//r, err := rpc.RawRequest(rpc.Method.GetValidatorInformation, "https://api.s0.os.hmny.io", []interface{}{"one1rhpfn58kvmmdmqfnw4uuzgedkvcfk7h67zsrc8"})
	//r, err := rpc.RawRequest(rpc.Method.GetValidatorInformation, BeaconChainEndpoint, []interface{}{"one1rhpfn58kvmmdmqfnw4uuzgedkvcfk7h67zsrc8"})
	r, err := rpc.RawRequest(rpc.Method.GetValidatorInformation, BeaconChainEndpoint, []interface{}{viper.GetString("OneAddress")})
	if err != nil {
		return ValidatorInformationReply{}, err
	}

	temp := reply{}
	err = json.Unmarshal(r, &temp)
	if err != nil {
		return ValidatorInformationReply{}, err
	}
	return temp.Result, nil
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
