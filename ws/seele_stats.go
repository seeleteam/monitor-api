/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"

	"github.com/seeleteam/monitor-api/config"
	"github.com/seeleteam/monitor-api/core/logs"
	"github.com/seeleteam/monitor-api/rpc"
)

// Service implements an Seele monitor reporting daemon that pushes local
// chain statistics up to a monitoring server.
type Service struct {
	rpc *rpc.MonitorRPC // json rpc

	hostname string // hostname of the node to display on the monitoring page
	node     string // Name of the node to display on the monitoring page
	pass     string // Password to authorize access to the monitoring page
	host     string // Remote address of the monitoring service
	port     int    // monitor api port
	shard    uint   // shard number
	wsRouter string // websocket base path
	wsPath   string // websocket path ex: {host:port}+{wsRouter}

	pongCh                     chan struct{} // Pong notifications are fed into this channel
	fullEventTickerTime        time.Duration
	latestBlockEventTickerTime time.Duration
	delayReConnTime            time.Duration //delay send msg to monitor when web socket server is not be connected
	delaySendTime              time.Duration //delay send msg to monitor when rpc server is not be connected
	latestBlockHeight          uint64        // record the latest block height, if rpc get the same block abort send
	currentBlockHeight         uint64        // record the current block height, if rpc get the same block abort send
	reportErrorAfterTimes      int           // report the error occur times (currentErrorTimes) when error occur over the special times
	currentErrorTimes          int
	currentNetVersion          string // current net version(netWorkId)
}

// New returns a monitoring service ready for stats reporting.
func New(url string, rpc *rpc.MonitorRPC) (*Service, error) {
	// first get RPC NodeInfo and according the Shard choose the ws path
ErrContinue:
	info, err := rpc.NodeInfo()
	if err != nil {
		logs.Error("rpc getNodeInfo error %v", err)
		time.Sleep(5 * time.Second)
		goto ErrContinue

	}
	shard := info.Shard
	websocketURL, _ := config.ShardMap[fmt.Sprintf("%v", shard)]
	if websocketURL == "" {
		logs.Error("shard config error, shard %v exist error web socket url", shard)
		return nil, err
	}
	// Parse the web socket connection url
	if url == "" {
		//addr format should be host:port!
		url = config.SeeleConfig.ServerConfig.Addr
	}

	currentConfig := config.SeeleConfig
	currentWebSocketConfig := currentConfig.ServerConfig.WebSocketConfig
	if currentWebSocketConfig == nil {
		return nil, fmt.Errorf("WebSocketConfig is nil")
	}
	wsRouter := currentWebSocketConfig.WsRouter
	re := regexp.MustCompile("([^:]*):(.+)")
	parts := re.FindStringSubmatch(url)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid websocket url: \"%s\", should be host:port", url)
	}
	port, err := strconv.Atoi(parts[2])
	if err != nil {
		logs.Error("parse url port %v error: %v", port, err)
		return nil, err
	}
	wsPath := fmt.Sprintf("%s%s", websocketURL, wsRouter)
	logs.Debug("init shard %v, wsPath is %v", shard, wsPath)
	host := parts[0]

	// name: INSTANCE_NAME || os.hostname()
	hostname := os.Getenv("INSTANCE_NAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	return &Service{
		rpc:                        rpc,
		hostname:                   hostname,
		node:                       hostname,
		host:                       host,
		port:                       port,
		shard:                      shard,
		wsRouter:                   wsRouter,
		wsPath:                     wsPath,
		pongCh:                     make(chan struct{}),
		fullEventTickerTime:        currentWebSocketConfig.WsFullEventTickerTime,
		latestBlockEventTickerTime: currentWebSocketConfig.WsLatestBlockEventTickerTime,
		delayReConnTime:            currentWebSocketConfig.DelayReConnTime,
		delaySendTime:              currentWebSocketConfig.DelaySendTime,
		reportErrorAfterTimes:      currentWebSocketConfig.ReportErrorAfterTimes,
		currentNetVersion:          info.NetVersion,
	}, nil
}

// Start start the loop for sending statics data to monitor server with web socket
func (s *Service) Start() {
	s.loop()
}

// loop keeps trying to connect to the monitor server, reporting chain events
// until termination.
func (s *Service) loop() {
	// Loop reporting until termination
	for {
		info, err := s.rpc.NodeInfo()
		if err != nil {
			logs.Error("rpc getNodeInfo error %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		shard := info.Shard
		websocketURL, _ := config.ShardMap[fmt.Sprintf("%v", shard)]

		wsPath := fmt.Sprintf("%s%s", websocketURL, s.wsRouter)
		s.wsPath = wsPath
		logs.Debug("now shard %v, wsPath %v", shard, wsPath)

		// Resolve the URL, defaulting to TLS, but falling back to none too
		path := fmt.Sprintf("%s", s.wsPath)
		urls := []string{path}
		if !strings.Contains(path, "://") {
			urls = []string{"wss://" + path, "ws://" + path}
			urls = append(urls, "https://"+path, "http://"+path)
		}

		// Establish a web socket connection to the server on any supported URL
		var (
			conf *websocket.Config
			conn *websocket.Conn
		)
		for _, url := range urls {
			if conf, err = websocket.NewConfig(url, "http://localhost/"); err != nil {
				continue
			}
			conf.Dialer = &net.Dialer{Timeout: 5 * time.Second}
			if conn, err = websocket.DialConfig(conf); err == nil {
				break
			}
		}
		if err != nil {
			logs.Warn("Stats server unreachable(resend after %v), err %v", s.delayReConnTime, err)
			time.Sleep(s.delayReConnTime)
			continue
		}

		go s.readLoop(conn)

		// first get the node base and append s.node
		coinBase, err := s.getCoinBase(conn)
		if err != nil {
			logs.Warn("Initial get coinBase failed(reconnect after %v), err %v", s.delaySendTime, err)
			if conn != nil {
				conn.Close()
			}
			time.Sleep(s.delaySendTime)
			continue
		}
		s.node = s.hostname + "_" + coinBase

		//Send the initial stats so our node looks decent from the get go
		if err = s.reportAllNodeInfo(conn); err != nil {
			logs.Warn("Initial stats report failed(reconnect after %v), err %v", s.delaySendTime, err)
			if conn != nil {
				conn.Close()
			}
			time.Sleep(s.delaySendTime)
			continue
		}

		fullReport := time.NewTicker(s.fullEventTickerTime)

		blockReport := time.NewTicker(s.latestBlockEventTickerTime)

		for err == nil {
			select {
			case <-fullReport.C:
				if err = s.report(conn); err != nil {
					logs.Warn("Full stats report failed", "err", err)
				}

			case <-blockReport.C:
				if err = s.reportCurrentBlock(conn); err != nil {
					logs.Warn("Current block report failed", "err", err)
				}
			}
		}
		// Make sure the connection is closed
		conn.Close()
	}
}

// readLoop loops as long as the connection is alive and retrieves data packets
// from the network socket. If any of them match an active request, it forwards
// it, if they themselves are requests it initiates a reply, and lastly it drops
// unknown packets.
func (s *Service) readLoop(conn *websocket.Conn) {
	// If the read loop exists, close the connection
	defer conn.Close()

	for {
		// Retrieve the next generic network packet and bail out on error
		var msg map[string][]interface{}
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			logs.Warn("Failed to decode stats server message", "err", err)
			return
		}
		logs.Debug("Received message from stats server", "msg", msg)
		if len(msg["emit"]) == 0 {
			logs.Warn("Stats server sent non-broadcast", "msg", msg)
			return
		}
		command, ok := msg["emit"][0].(string)
		if !ok {
			logs.Warn("Invalid stats server message type", "type", msg["emit"][0])
			return
		}
		// If the message is a ping reply, deliver (someone must be listening!)
		if len(msg["emit"]) == 2 && command == "node-pong" {
			select {
			case s.pongCh <- struct{}{}:
				// Pong delivered, continue listening
				continue
			default:
				// Ping routine dead, abort
				logs.Warn("Stats server pinger seems to have died")
				return
			}
		}
		// Report anything else and continue
		logs.Info("stats message", "msg", msg)
	}
}

// nodeInfo is the collection of metainformation about a node that is displayed
// on the monitoring page.
type nodeInfo struct {
	Name        string `json:"name"`
	Node        string `json:"node"`
	Port        int    `json:"port"` // the monitor api client port, can overwrite use monitor api client api port
	Protocol    string `json:"protocol"`
	API         string `json:"api"`
	Os          string `json:"os"`
	OsVer       string `json:"os_v"`
	Client      string `json:"client"`
	NodeVersion string `json:"nodeVersion"` // the monitor api client version
	NetVersion  string `json:"netVersion"`
	Shard       uint   `json:"shard"`
}

// nodeStats is the information about the local node.
type nodeStats struct {
	Active   bool `json:"active"`
	Syncing  bool `json:"syncing"`
	Mining   bool `json:"mining"`
	Hashrate int  `json:"hashrate"`
	Peers    int  `json:"peers"`
	GasPrice int  `json:"gasPrice"`
	Uptime   int  `json:"uptime"`
}

type apiCurrentBlock struct {
	HeadHash   string   `json:"headHash"`
	Height     uint64   `json:"height"`
	Timestamp  *big.Int `json:"timestamp"`
	Difficulty *big.Int `json:"difficulty"`
	Creater    string   `json:"miner"`
	TxCount    int      `json:"txcount"`
}

// report collects all possible data to report and send it to the stats server.
// This should only be used on reconnects or rarely to avoid overloading the
// server. Use the individual methods for reporting subscribed events.
func (s *Service) report(conn *websocket.Conn) error {
	if err := s.reportLatency(conn); err != nil {
		return err
	}

	if err := s.reportNodeStats(conn); err != nil {
		return err
	}

	return nil
}

// reportLatency sends a ping request to the server, measures the RTT time and
// finally sends a latency update.
func (s *Service) reportLatency(conn *websocket.Conn) error {
	// Send the current time to the monitor server
	latency, err := s.getLatency(conn)
	if err != nil {
		return errors.New(err.Error())
	}

	report := map[string][]interface{}{
		"emit": {"latency", map[string]interface{}{
			"id":         s.node,
			"latency":    latency,
			"netVersion": s.currentNetVersion,
			"shard":      s.shard,
		}},
	}
	// Send back the measured latency
	logs.Debug("Sending measured latency to seele monitor", "latency", latency)
	jsonReport, _ := json.Marshal(report)
	logs.Debug("Sending node latency to monitor\n %v", string(jsonReport))
	return websocket.JSON.Send(conn, report)
}

// reportNodeInfo retrieves various stats about the node at the networking and
// mining layer and reports it to the stats server.
func (s *Service) reportNodeInfo(conn *websocket.Conn) error {
	nodeInfo, err := s.getNodeInfo(conn)

	if err != nil {
		return errors.New(err.Error())
	}
	report := map[string][]interface{}{
		"emit": {"nodeInfo", nodeInfo},
	}
	jsonReport, _ := json.Marshal(report)
	logs.Debug("Sending node info to monitor\n %v", string(jsonReport))
	return websocket.JSON.Send(conn, report)
}

// reportNodeStats retrieves various stats about the node at the networking and
// mining layer and reports it to the stats server.
func (s *Service) reportNodeStats(conn *websocket.Conn) error {
	nodeStats, err := s.getNodeStats(conn)
	if err != nil {
		logs.Error("rpc reportNodeStats error %v", err)
		return err
	}
	report := map[string][]interface{}{
		"emit": {"stats", nodeStats},
	}
	jsonReport, _ := json.Marshal(report)
	logs.Debug("Sending node stats to monitor\n %v", string(jsonReport))
	return websocket.JSON.Send(conn, report)
}

func (s *Service) getLatency(conn *websocket.Conn) (string, error) {
	// Send the current time to the monitor server
	start := time.Now()

	ping := map[string][]interface{}{
		"emit": {"node-ping", map[string]interface{}{
			"id":         s.node,
			"clientTime": start.UnixNano() / 1000000,
			"netVersion": s.currentNetVersion,
			"shard":      s.shard,
		}},
	}
	jsonReport, _ := json.Marshal(ping)
	logs.Debug("Sending node ping to monitor\n %v", string(jsonReport))
	if err := websocket.JSON.Send(conn, ping); err != nil {
		logs.Error("rpc reportLatency error %v", err)
		return "-1", err

	}

	// Wait for the pong request to arrive back
	select {
	case <-s.pongCh:
		// Pong delivered, report the latency
	case <-time.After(5 * time.Second):
		// Ping timeout, abort
		return "-1", errors.New("ping timed out")
	}
	latencyFloat := float32(int((time.Since(start)/time.Duration(2)).Nanoseconds()*10)) / 10000000
	latency := fmt.Sprintf("%.1f", latencyFloat)
	logs.Debug("latency is %vms", latency)
	return latency, nil
}

func (s *Service) getNodeInfo(conn *websocket.Conn) (map[string]interface{}, error) {
	info, err := s.rpc.NodeInfo()
	if err != nil {
		logs.Error("rpc getNodeInfo error %v", err)
		s.detectErrorAndReport(conn)
		return nil, err
	}

	// update netVersion
	s.currentNetVersion = info.NetVersion
	s.shard = info.Shard

	nodeInfoData := nodeInfo{
		Name:        config.APPName,
		NodeVersion: config.VERSION,
		Node:        info.Node,
		Port:        s.port,
		Protocol:    info.Protocol,
		Os:          info.Os,
		OsVer:       info.OsVer,
		Client:      info.Client,
		API:         info.Protocol,
		NetVersion:  s.currentNetVersion,
		Shard:       s.shard,
	}

	nodeInfo := map[string]interface{}{
		"id":   s.node,
		"info": nodeInfoData,
	}

	return nodeInfo, nil
}

func (s *Service) getNodeStats(conn *websocket.Conn) (map[string]interface{}, error) {
	stats, err := s.rpc.NodeStats()
	if err != nil {
		logs.Error("rpc getNodeStats error %v", err)
		s.detectErrorAndReport(conn)
		return nil, err
	}
	nodeStats := map[string]interface{}{
		"id":         s.node,
		"stats":      stats,
		"netVersion": s.currentNetVersion,
		"shard":      s.shard,
	}

	return nodeStats, nil
}

func (s *Service) reportCurrentBlock(conn *websocket.Conn) error {
	if err := s.reportCurrentBlockInfo(conn); err != nil {
		return err
	}
	return nil
}

func (s *Service) getCurrentBlockInfo(conn *websocket.Conn) (map[string]interface{}, error) {
	block, err := s.rpc.CurrentBlock(-1, true)
	if err != nil {
		logs.Error("rpc getCurrentBlockInfo error %v", err)
		s.detectErrorAndReport(conn)
		return nil, err
	}

	blockInfo := map[string]interface{}{
		"id": s.node,
		"block": &apiCurrentBlock{
			HeadHash:   block.HeadHash,
			Height:     block.Height,
			Timestamp:  block.Timestamp,
			Difficulty: block.Difficulty,
			Creater:    block.Creator,
			TxCount:    block.TxCount,
		},
		"netVersion": s.currentNetVersion,
		"shard":      s.shard,
	}
	s.currentBlockHeight = block.Height
	return blockInfo, nil
}

// reportCurrentBlockInfo retrieves various stats about the node at the networking and
// mining layer and reports it to the stats server.
func (s *Service) reportCurrentBlockInfo(conn *websocket.Conn) error {
	blockInfo, err := s.getCurrentBlockInfo(conn)
	if err != nil {
		logs.Error("rpc reportCurrentBlockInfo error %v", err)
		return err
	}

	// if current block height gt the prev block height send the block info
	if s.currentBlockHeight > s.latestBlockHeight {
		s.latestBlockHeight = s.currentBlockHeight
		report := map[string][]interface{}{
			"emit": {"block", blockInfo},
		}
		jsonReport, _ := json.Marshal(report)
		logs.Debug("Sending node current block to monitor\n %v", string(jsonReport))
		return websocket.JSON.Send(conn, report)
	} else {
		logs.Debug("no Sending node current block to monitor, currentBlockHeight: %v, latestBlockHeight: %v", s.currentBlockHeight, s.latestBlockHeight)
	}
	return nil
}

// reportAllNodeInfo send this info to monitor, the first start conn or reconnect
func (s *Service) reportAllNodeInfo(conn *websocket.Conn) error {
	// nodeInfo must come first
	info, err := s.getNodeInfo(conn)
	if err != nil {
		logs.Error("reportAllNodeInfo %v", err)
		return err
	}

	block, err := s.getCurrentBlockInfo(conn)
	if err != nil {
		logs.Error("reportAllNodeInfo %v", err)
		return err
	}

	stats, err := s.getNodeStats(conn)
	if err != nil {
		logs.Error("reportAllNodeInfo %v", err)
		return err
	}
	latency, err := s.getLatency(conn)
	if err != nil {
		logs.Error("reportAllNodeInfo %v", err)
		return err
	}

	allNodeInfo := map[string]interface{}{
		"id":         s.node,
		"info":       info["info"],
		"block":      block["block"],
		"stats":      stats["stats"],
		"latency":    latency,
		"netVersion": s.currentNetVersion,
		"shard":      s.shard,
	}
	report := map[string][]interface{}{
		"emit": {"hello", allNodeInfo},
	}
	jsonReport, _ := json.Marshal(report)
	logs.Debug("Sending node all info to monitor\n %v", string(jsonReport))
	return websocket.JSON.Send(conn, report)
}

func (s *Service) getCoinBase(conn *websocket.Conn) (string, error) {
	info, err := s.rpc.GetInfo()
	if err != nil {
		logs.Error("rpc getCoinBase error %v", err)
		s.detectErrorAndReport(conn)
		return "", err
	}
	coinBase := info["Coinbase"]
	v, ok := coinBase.(string)
	if !ok {
		return "", errors.New("coinBase error")
	}

	return v, nil
}

// detectErrorAndReport detect the error and report to monitor
func (s *Service) detectErrorAndReport(conn *websocket.Conn) error {
	s.currentErrorTimes++
	if s.currentErrorTimes >= s.reportErrorAfterTimes {
		logs.Error("conn error occur times: %v >= %v, will report error", s.currentErrorTimes, s.reportErrorAfterTimes)
		s.currentErrorTimes = 0
		return s.reportServerError(conn)
	}
	logs.Debug("conn error occur times: %v < %v", s.currentErrorTimes, s.reportErrorAfterTimes)
	return nil
}

// reportServerError report the error to monitor server
func (s *Service) reportServerError(conn *websocket.Conn) error {
	nodeStats := map[string]interface{}{
		"id": s.node,
		"stats": map[string]interface{}{
			"active":  false,
			"syncing": false,
		},
		"netVersion": s.currentNetVersion,
		"shard":      s.shard,
	}
	report := map[string][]interface{}{
		"emit": {"stats", nodeStats},
	}
	jsonReport, _ := json.Marshal(report)
	logs.Debug("Sending node error info to monitor\n %v", string(jsonReport))
	return websocket.JSON.Send(conn, report)
}
