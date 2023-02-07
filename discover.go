package openKeeper

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"strings"
)

func (s *ZkClient) watch() {
	for {
		select {
		case event := <-s.eventChan:
			switch event.Type {
			case zk.EventSession:
			case zk.EventNodeCreated:
			case zk.EventNodeChildrenChanged:
				l := strings.Split(event.Path, "/")
				s.lock.Lock()
				delete(s.rpcLocalCache, l[len(l)-1])
				fmt.Println("update", s.rpcLocalCache, s.node)
				s.lock.Unlock()
			case zk.EventNodeDataChanged:
			case zk.EventNodeDeleted:
				fmt.Println("de")
			case zk.EventNotWatching:
			}
		}
	}
}

func (s *ZkClient) GetConnsRemote(serviceName string, opts ...grpc.DialOption) (conns []*grpc.ClientConn, err error) {
	path := s.getPath(serviceName)
	childNodes, _, err := s.conn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, nil
		}
		return nil, err
	}
	for _, child := range childNodes {
		fullPath := path + "/" + child
		data, _, err := s.conn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				continue
			}
			return nil, err
		}
		conn, err := grpc.Dial(string(data), opts...)
		if err == nil {
			conns = append(conns, conn)
		}
	}
	return conns, nil
}

func (s *ZkClient) GetConns(serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	conns, ok := s.rpcLocalCache[serviceName]
	if !ok {
		fmt.Println("remote", s.node)
		var err error
		conns, err = s.GetConnsRemote(serviceName, opts...)
		if err != nil {
			return nil, err
		}
		_, _, _, err = s.conn.ChildrenW(s.getPath(serviceName))
		if err != nil {
			return nil, err
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.rpcLocalCache[serviceName] = conns
	}
	return conns, nil
}

func (s *ZkClient) GetConnStrategy(serviceName string, strategy func(slice []*grpc.ClientConn) int, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	conns, err := s.GetConns(serviceName, opts...)
	if len(conns) > 0 {
		return conns[strategy(conns)], nil
	}
	return nil, err
}

func (s *ZkClient) GetConn(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return s.GetConnStrategy(serviceName, func(slice []*grpc.ClientConn) int { return 0 }, opts...)
}
