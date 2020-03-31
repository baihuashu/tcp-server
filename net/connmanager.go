package net

import (
	"errors"
	"fmt"
	"sync"
	"github.com/baihuashu/tcp-server/iface"
)

type ConnManager struct {
	conns    map[uint32]iface.IConnection
	connLock sync.RWMutex //保护连接集合的读写锁

}

func NewConnMgr() *ConnManager {
	return &ConnManager{
		conns:    make(map[uint32]iface.IConnection),
		connLock: sync.RWMutex{},
	}
}

func (connMgr *ConnManager) Add(conn iface.IConnection) {
	//保护共享资源map，写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//添加
	connMgr.conns[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num =", connMgr.Len())
}
func (connMgr *ConnManager) Remove(conn iface.IConnection) {
	//保护共享资源map，写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除
	delete(connMgr.conns, conn.GetConnID())
	fmt.Println("remove connection successfully: conn num =", connMgr.Len())

}
func (connMgr *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	//加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.conns[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found!")
	}

}
func (connMgr *ConnManager) Len() int {
	return len(connMgr.conns)
}
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源map，写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止conn工作并删除

	for connID,conn := range connMgr.conns{
		//停止
		conn.Stop()
		//删除
		delete(connMgr.conns,connID)
	}
	fmt.Println("clear all conn successfully! conn num = ",connMgr.Len())
}
