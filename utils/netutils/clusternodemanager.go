// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package netutils

import (
	"fmt"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidCluster = fmt.Errorf("not found cluster addr info")
	ErrNoValidNode    = fmt.Errorf("not found valid node")
)

//
// ClusterNodeManager
//
// it manager cluster node
//
type ClusterNodeManager struct {
	allMutex        sync.RWMutex  // sync lock: allNodes, brokenNodes, availableNodes
	allNodes        []string      // cluster node addr list: ip:port
	brokenMutex     sync.RWMutex  // sync lock: allNodes, brokenNodes, availableNodes
	brokenNodes     []string      // current unavailable cluster node: ip:port
	availMutex      sync.RWMutex  // sync lock: allNodes, brokenNodes, availableNodes
	availableNodes  []string      // current available cluster node: ip:port
	healthyInterval uint64        // healthy check interval
	idx             uint64        // robin parameter, atomic operate
	stopCh          chan struct{} // Close signal
	logger          logging.Logger
}

//
// @param servers cluster address list, ip:port
//
func NewClusterNodeManager(interval uint64, logger logging.Logger, servers ...string) (*ClusterNodeManager, error) {
	if len(servers) <= 0 {
		return nil, ErrInvalidCluster
	}

	cli := &ClusterNodeManager{
		stopCh:          make(chan struct{}),
		allNodes:        make([]string, 0, len(servers)),
		healthyInterval: interval,
		logger:          logger,
	}

	for _, server := range servers {
		ipPort := strings.Split(server, ":")
		if len(ipPort) != 2 {
			logger.Log("file", "clusternodemanager.go",
				"function", "NewClusterNodeManager",
				"action", "check servers",
				"server", server,
				"error", "failed format check")
			continue
		}

		portStr := ipPort[1]
		_, err := strconv.ParseUint(portStr, 10, 16)
		if err != nil {
			logger.Log("file", "clusternodemanager.go",
				"function", "NewClusterNodeManager",
				"action", "check port",
				"port", portStr,
				"error", err)
			return nil, err
		}
		cli.allNodes = append(cli.allNodes, server)
	}

	if len(cli.allNodes) <= 0 {
		logger.Log("file", "clusternodemanager.go",
			"function", "NewClusterNodeManager",
			"action", "check available node",
			"error", "not found available node")
		return nil, nil
	}

	cli.availableNodes = make([]string, len(cli.allNodes))
	copy(cli.availableNodes, cli.allNodes)
	cli.checkAvailableNodes()
	go cli.heart()
	return cli, nil
}

//
// retrieve random one cluster info
//
func (cnm *ClusterNodeManager) Random() (string, error) {
	if len(cnm.availableNodes) <= 0 {
		return "", ErrNoValidNode
	}

	res := atomic.AddUint64(&cnm.idx, 1)

	cnm.availMutex.RLock()
	defer cnm.availMutex.RUnlock()
	robinIdx := res % uint64(len(cnm.availableNodes))
	return cnm.availableNodes[robinIdx], nil
}

//
// all cluster nodes
//
func (cnm *ClusterNodeManager) All() []string {
	cnm.allMutex.RLock()
	defer cnm.allMutex.RUnlock()
	l := make([]string, len(cnm.allNodes))
	copy(l, cnm.allNodes)
	return l
}

//
// all available cluster nodes
//
func (cnm *ClusterNodeManager) Available() []string {
	cnm.availMutex.RLock()
	defer cnm.availMutex.RUnlock()
	l := make([]string, len(cnm.availableNodes))
	copy(l, cnm.availableNodes)
	return l
}

//
// close manager
//
func (cnm *ClusterNodeManager) Close() {
	cnm.stopCh <- struct{}{}
}

//
// listen healthy signal and close signal
//
func (cnm *ClusterNodeManager) heart() {
	interval := time.Duration(cnm.healthyInterval)
	checkChan := time.Tick(interval * time.Second)
	for {
		select {
		case <-checkChan:
			cnm.checkAvailableNodes()
			cnm.checkBrokenNodes()
			cnm.logger.Log("file", "clusternodemanager.go",
				"function", "heart",
				"action", "time check node healthy",
				"healthy", len(cnm.availableNodes),
				"ill", len(cnm.brokenNodes))

		case <-cnm.stopCh:
			cnm.logger.Log("file", "clusternodemanager.go",
				"function", "heart",
				"action", "close")
			break
		}
	}
}

//
// check available node
//
func (cnm *ClusterNodeManager) checkAvailableNodes() {
	cnm.availMutex.RLock()
	l := make([]string, len(cnm.availableNodes))
	copy(l, cnm.availableNodes)
	cnm.availMutex.RUnlock()

	for _, tmp := range l {
		if !cnm.checkSocket(tmp) {
			cnm.setBrokenNodes(tmp)
		}
	}
}

func (cnm *ClusterNodeManager) checkBrokenNodes() {
	cnm.brokenMutex.RLock()
	l := make([]string, len(cnm.brokenNodes))
	copy(l, cnm.brokenNodes)
	cnm.brokenMutex.RUnlock()

	for _, tmp := range l {
		if cnm.checkSocket(tmp) {
			cnm.setAvailableNodes(tmp)
		}
	}
}

func (cnm *ClusterNodeManager) setBrokenNodes(node string) {
	cnm.availMutex.Lock()
	defer cnm.availMutex.Unlock()
	for i, e := range cnm.availableNodes {
		if e == node {
			l := cnm.availableNodes[0:i]
			r := cnm.availableNodes[i+1:]
			cnm.availableNodes = append(l, r...)
			break
		}
	}
	cnm.brokenMutex.Lock()
	defer cnm.brokenMutex.Unlock()
	cnm.brokenNodes = append(cnm.brokenNodes, node)
}

func (cnm *ClusterNodeManager) setAvailableNodes(node string) {
	cnm.brokenMutex.Lock()
	defer cnm.brokenMutex.Unlock()
	for i, e := range cnm.brokenNodes {
		if e == node {
			l := cnm.brokenNodes[0:i]
			r := cnm.brokenNodes[i+1:]
			cnm.brokenNodes = append(l, r...)
			break
		}
	}
	cnm.availMutex.Lock()
	defer cnm.availMutex.Unlock()
	cnm.availableNodes = append(cnm.availableNodes, node)
}

func (cnm *ClusterNodeManager) checkSocket(hostPort string) bool {
	addr, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		cnm.logger.Log("file", "clusternodemanager.go",
			"function", "checkSocket",
			"action", "resolver ip port",
			"error", err)
		return false
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		cnm.logger.Log("file", "clusternodemanager.go",
			"function", "checkSocket",
			"action", "dial ip port",
			"error", err)
		return false
	}
	defer conn.Close()

	return true
}
