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
	err = rpc.call("monitor_nodeStats", nil, &nodeStats)
	if err != nil {
		return
	}

	var hashrate uint64
	err = rpc.call("miner_getHashrate", nil, &hashrate)
	if err != nil {
		return
	}
	nodeStats.Hashrate = hashrate
	return nodeStats, err
}

// NodeInfo returns the current node info.
func (rpc *MonitorRPC) NodeInfo() (nodeInfo *NodeInfo, err error) {
	err = rpc.call("monitor_nodeInfo", nil, &nodeInfo)
	return nodeInfo, err
}

// CurrentBlock returns the current block info.
func (rpc *MonitorRPC) CurrentBlock(h int64, fullTx bool) (currentBlock *CurrentBlock, err error) {
	request := GetBlockByHeightRequest{
		Height: h,
		FullTx: fullTx,
	}
	var req []interface{}
	req = append(req, request.Height)
	req = append(req, request.FullTx)
	rpcOutputBlock := make(map[string]interface{})
	if err := rpc.call("seele_getBlockByHeight", req, &rpcOutputBlock); err != nil {
		return nil, err
	}

	return getBlockByHeight(rpcOutputBlock, fullTx), err
}
func getBlockByHeight(rpcOutputBlock map[string]interface{}, fullTx bool) *CurrentBlock {
	headerMp := rpcOutputBlock["header"].(map[string]interface{})
	timestamp := int64(headerMp["CreateTimestamp"].(float64))
	difficulty := int64(headerMp["Difficulty"].(float64))
	height := uint64(headerMp["Height"].(float64))
	creator := headerMp["Creator"].(string)

	return &CurrentBlock{
		HeadHash:   rpcOutputBlock["hash"].(string),
		Height:     height,
		Timestamp:  big.NewInt(timestamp),
		Difficulty: big.NewInt(difficulty),
		Creator:    creator,
		TxCount:    len(rpcOutputBlock["transactions"].([]interface{})),
	}
}

// GetInfo gets the account address that mining rewards will be send to.
func (rpc *MonitorRPC) GetInfo() (result map[string]interface{}, err error) {
	err = rpc.call("seele_getInfo", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
