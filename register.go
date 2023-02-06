package openKeeper

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func (s *ZkClient) Register(rpcRegisterName, host string, port int, opts ...grpc.DialOption) error {
	if s.isRegister {
		return nil
	}
	if err := s.ensureName(rpcRegisterName); err != nil {
		return err
	}
	addr := s.getAddr(host, port)
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}
	node, err := s.conn.CreateProtectedEphemeralSequential(s.getPath(rpcRegisterName), []byte(addr), zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	s.rpcRegisterName = rpcRegisterName
	s.rpcRegisterAddr = addr
	s.isRegister = true
	s.node = node
	s.rpcLocalCache[rpcRegisterName] = append(s.rpcLocalCache[rpcRegisterName], conn)
	return err
}

func (s *ZkClient) UnRegister() error {
	if !s.isRegister {
		return errors.New(fmt.Sprintf("service %s is not register", s.rpcRegisterName))
	}
	s.lock.Lock()
	addresses, ok := s.rpcLocalCache[s.rpcRegisterName]
	if !ok || len(addresses) == 0 {
		return errors.New(fmt.Sprintf("service %s is not register, could not find conns in Local", s.rpcRegisterName))
	}
	// delete our hash of the service
	for index, node := range s.rpcLocalCache[s.rpcRegisterName] {
		if s.rpcRegisterAddr == node.Target() {
			err := s.conn.Delete(s.getPath(node.Target()), -1)
			if err != nil {
				return err
			}
			if index == len(s.rpcLocalCache[s.rpcRegisterName])-1 {
				s.rpcLocalCache[s.rpcRegisterName] = s.rpcLocalCache[s.rpcRegisterName][:index]
			} else {
				s.rpcLocalCache[s.rpcRegisterName] = append(s.rpcLocalCache[s.rpcRegisterName][:index], s.rpcLocalCache[s.rpcRegisterName][index+1:]...)
			}
		}

	}
	s.lock.Unlock()
	s.isRegister = false
	return nil
}
