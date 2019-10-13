package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/harmony-one/harmony-tui/rpc"
	"github.com/harmony-one/harmony/common/denominations"
)

// RPCRoutes reflects the RPC endpoints of the target network across shards
type RPCRoutes struct {
	HTTP    string `json:"http"`
	ShardID int    `json:"shardID"`
	WS      string `json:"ws"`
}

// Structure produces a slice of RPCRoutes for the network across shards
func Structure(node string) ([]RPCRoutes, error) {
	type r struct {
		Result []RPCRoutes `json:"result"`
	}
	p, e := rpc.RawRequest(rpc.Method.GetShardingStructure, node, []interface{}{})
	if e != nil {
		return nil, e
	}
	result := r{}
	json.Unmarshal(p, &result)
	return result.Result, nil
}

func CheckAllShards(node, addr string, noPretty bool) (string, error) {
	var out bytes.Buffer
	out.WriteString("[")
	params := []interface{}{addr, "latest"}
	s, err := Structure(node)
	if err != nil {
		return "", err
	}
	for i, shard := range s {
		balanceRPCReply, _ := rpc.Request(rpc.Method.GetBalance, shard.HTTP, params)
		balance, _ := balanceRPCReply["result"].(string)
		bln, _ := big.NewInt(0).SetString(balance[2:], 16)
		out.WriteString(fmt.Sprintf(`{"shard":%d, "amount":%s}`,
			shard.ShardID,
			ConvertBalanceIntoReadableFormat(bln),
		))
		if i != len(s)-1 {
			out.WriteString(",")
		}
	}
	out.WriteString("]")
	if noPretty {
		return out.String(), nil
	}
	return JSONPrettyFormat(out.String()), nil
}

func JSONPrettyFormat(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

func ConvertBalanceIntoReadableFormat(balance *big.Int) string {
	balance = balance.Div(balance, big.NewInt(denominations.Nano))
	strBalance := fmt.Sprintf("%d", balance.Uint64())

	bytes := []byte(strBalance)
	hasDecimal := false
	for i := 0; i < 11; i++ {
		if len(bytes)-1-i < 0 {
			bytes = append([]byte{'0'}, bytes...)
		}
		if bytes[len(bytes)-1-i] != '0' && i < 9 {
			hasDecimal = true
		}
		if i == 9 {
			newBytes := append([]byte{'.'}, bytes[len(bytes)-i:]...)
			bytes = append(bytes[:len(bytes)-i], newBytes...)
		}
	}
	zerosToRemove := 0
	for i := 0; i < len(bytes); i++ {
		if hasDecimal {
			if bytes[len(bytes)-1-i] == '0' {
				bytes = bytes[:len(bytes)-1-i]
				i--
			} else {
				break
			}
		} else {
			if zerosToRemove < 5 {
				bytes = bytes[:len(bytes)-1-i]
				i--
				zerosToRemove++
			} else {
				break
			}
		}
	}
	return string(bytes)
}
