package openKeeper

import (
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func (s *ZkClient) createRpcNode(rpcRegisterName string) error {
	isExist, _, err := s.conn.Exists(s.getPath(rpcRegisterName))
	if err != nil {
		return err
	}
	if !isExist {
		_, err = s.conn.Create(s.getPath(rpcRegisterName), []byte(rpcRegisterName), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *ZkClient) Register(rpcRegisterName, host string, port int, opts ...grpc.DialOption) error {
	if err := s.createRpcNode(rpcRegisterName); err != nil {
		return err
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if err := s.ensureName(rpcRegisterName); err != nil {
		return err
	}
	addr := s.getAddr(host, port)
	_, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}

	node, err := s.conn.CreateProtectedEphemeralSequential(s.getPath(rpcRegisterName)+"/"+addr+"_", []byte(addr), zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	go s.watch()
	s.node = node
	return nil
}

func (s *ZkClient) UnRegister() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	err := s.conn.Delete(s.node, -1)
	if err != nil {
		return err
	}
	s.rpcLocalCache = map[string][]*grpc.ClientConn{}
	s.node = ""
	return nil
}
