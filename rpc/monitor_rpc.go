/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package rpc

import (
	"math/big"
)

// NodeStats returns the current node info.
func (rpc *MonitorRPC) NodeStats() (nodeStats *NodeStats, err error) {
	err = rpc.call("monitor.NodeStats", nil, &nodeStats)
	return nodeStats, err
}

// NodeInfo returns the current node info.
func (rpc *MonitorRPC) NodeInfo() (nodeInfo *NodeInfo, err error) {
	err = rpc.call("monitor.NodeInfo", nil, &nodeInfo)
	return nodeInfo, err
}

// CurrentBlock returns the current block info.
func (rpc *MonitorRPC) CurrentBlock() (currentBlock *CurrentBlock, err error) {
	request := GetBlockByHeightRequest{
		Height: -1,
		FullTx: true,
	}
	rpcOutputBlock := make(map[string]interface{})
	if err := rpc.call("seele.GetBlockByHeight", request, &rpcOutputBlock); err != nil {
		return nil, err
	}

	timestamp := int64(rpcOutputBlock["timestamp"].(float64))
	difficulty := int64(rpcOutputBlock["difficulty"].(float64))
	height := uint64(rpcOutputBlock["height"].(float64))

	currentBlock = &CurrentBlock{
		HeadHash:   rpcOutputBlock["hash"].(string),
		Height:     height,
		Timestamp:  big.NewInt(timestamp),
		Difficulty: big.NewInt(difficulty),
		Creator:    rpcOutputBlock["creator"].(string),
		TxCount:    len(rpcOutputBlock["transactions"].([]interface{})),
	}
	return currentBlock, err
}

// GetInfo gets the account address that mining rewards will be send to.
func (rpc *MonitorRPC) GetInfo() (result map[string]interface{}, err error) {
	err = rpc.call("seele.GetInfo", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
