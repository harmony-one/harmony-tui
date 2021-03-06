package data

import (
	"encoding/json"
	"regexp"
	"strconv"
	"time"

	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/harmony/numeric"
	"github.com/spf13/viper"
)

var (
	VersionData map[string]interface{}
	Announce    string
	OnAnnounce  string
	OnPrepared  string
	Bingo       string
	OnCommitted string

	LatestHeader  LatestHeaderReply
	LatestBlock   BlockByNumberReply
	Metadata      NodeMetadataReply
	ValidatorInfo ValidatorInformationReply
	LifetimeAvail numeric.Dec
	PeerCount     int64
	Balance       string
	TotalBalance  float64

	Quitter func(string)

	oneAddressPattern   = regexp.MustCompile("one1[0-9a-z]+")
	EarningRate         = numeric.NewDec(0)
	BeaconChainEndpoint = ""
)

func RefreshData() {

	for range time.Tick(viper.GetDuration("RPCRefreshInterval")) {

		// Only need to successfully get field once
		if BeaconChainEndpoint == "" {
			if shardingReply, err := getShardingStructure(); err == nil {
				for _, s := range shardingReply {
					if s.ShardID == uint32(0) {
						BeaconChainEndpoint = s.HTTP
						break
					}
				}
			}
		}

		if latestHeader, err := getLatestHeader(); err != nil {
			//If latestHeader fails, do not update anything
			continue
		} else {
			LatestHeader = latestHeader
			hexaBlockNumber := numToHex(LatestHeader.BlockNumber)
			// Only get block data if latest header request succeeds
			if latestBlockReply, err := getBlockByNumber(hexaBlockNumber); err != nil {
				LatestBlock = latestBlockReply
			} else {
				LatestBlock = latestBlockReply
				LatestBlock.BlockSizeInt = int(hexToNum(LatestBlock.BlockSize))
				LatestBlock.NumTransactions = len(LatestBlock.Transactions)
				LatestBlock.NumStakingTransactions = len(LatestBlock.StakingTransactions)
			}
		}

		if metadataReply, err := getNodeMetadata(); err == nil {
			Metadata = metadataReply
		}

		if peerCountReply, err := getPeerCount(); err == nil {
			PeerCount = hexToNum(peerCountReply)
		}

		Balance, TotalBalance = GetBalance()

		if BeaconChainEndpoint != "" {
			if validatorReply, err := getValidatorInformation(); err == nil {
				ValidatorInfo = validatorReply
				if lifetime := ValidatorInfo.Lifetime; lifetime != nil {
					lifetimeSigned := numeric.NewDecFromBigInt(ValidatorInfo.Lifetime.Signing.NumBlocksSigned)
					lifetimeToSign := numeric.NewDecFromBigInt(ValidatorInfo.Lifetime.Signing.NumBlocksToSign)
					if lifetimeToSign.GT(numeric.NewDec(0)) {
						LifetimeAvail = lifetimeSigned.Quo(lifetimeToSign)
					} else {
						LifetimeAvail = numeric.NewDec(0)
					}
				} else {
					LifetimeAvail = numeric.NewDec(0)
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
	if err = json.Unmarshal(r, &temp); err != nil {
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
	if err = json.Unmarshal(r, &temp); err != nil {
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
	if err = json.Unmarshal(r, &temp); err != nil {
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
	if err = json.Unmarshal(r, &temp); err != nil {
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
	if err = json.Unmarshal(r, &temp); err != nil {
		return NodeMetadataReply{}, err
	}
	return temp.Result, nil
}

// Always query BeaconChainEndpoint to get latest validator information on chain
func getValidatorInformation() (ValidatorInformationReply, error) {
	type reply struct {
		Result ValidatorInformationReply `json:"result"`
	}

	r, err := rpc.RawRequest(rpc.Method.GetValidatorInformation, BeaconChainEndpoint, []interface{}{viper.GetString("OneAddress")})
	if err != nil {
		return ValidatorInformationReply{}, err
	}

	temp := reply{}
	if err = json.Unmarshal(r, &temp); err != nil {
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
		if err := json.Unmarshal([]byte(balance), &temp); err != nil {
			balance = "No data"
			return balance, tempBal
		}
		balance = "Address: " + viper.GetString("OneAddress")

		for _, b := range temp {
			balance += "\n Balance in Shard " + strconv.FormatFloat(b["shard"].(float64), 'f', 0, 64) + ":  " + strconv.FormatFloat(b["amount"].(float64), 'f', 4, 64)
			tempBal += b["amount"].(float64)
		}
	}
	return balance, tempBal
}
