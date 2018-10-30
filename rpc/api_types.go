/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package rpc

import (
	"math/big"

	"github.com/seeleteam/go-seele/common"
)

// NodeInfo is the collection of meta information about a node that is displayed
// on the monitoring page.
type NodeInfo struct {
	Name       string `json:"name"`
	Node       string `json:"node"`
	Port       int    `json:"port"`
	NetVersion uint64 `json:"netVersion"`
	Protocol   string `json:"protocol"`
	API        string `json:"api"`
	Os         string `json:"os"`
	OsVer      string `json:"os_v"`
	Client     string `json:"client"`
	History    bool   `json:"canUpdateHistory"`
	Shard      uint   `json:"shard"`
}

// NodeStats is the information about the local node.
type NodeStats struct {
	Active   bool   `json:"active"`
	Syncing  bool   `json:"syncing"`
	Mining   bool   `json:"mining"`
	Hashrate uint64 `json:"hashrate"`
	Peers    int    `json:"peers"`
}

// CurrentBlock is the informations about the best block
type CurrentBlock struct {
	HeadHash   string   `json:"headHash"`
	Height     uint64   `json:"height"`
	Timestamp  *big.Int `json:"timestamp"`
	Difficulty *big.Int `json:"difficulty"`
	Creator    string   `json:"creator"`
	TxCount    int      `json:"txcount"`
}

// MinerInfo miner simple info
type MinerInfo struct {
	Coinbase           common.Address
	CurrentBlockHeight uint64
	HeaderHash         common.Hash
}
