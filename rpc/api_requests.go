/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package rpc

// GetBlockByHeightRequest request param for GetBlockByHeight api
type GetBlockByHeightRequest struct {
	Height int64
	FullTx bool
}
