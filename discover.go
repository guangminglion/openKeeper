package openKeeper

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func (s *ZkClient) watch() {
	for {
		select {
		case event := <-s.eventChan:
			switch event.Type {
			case zk.EventNodeCreated:
				fmt.Println("EventNodeCreated", event.Path, event.Err, event.Server, event.State)
			case zk.EventNodeChildrenChanged:
				fmt.Println("EventNodeChildrenChanged", event.Path, event.Err, event.Server, event.State)

			case zk.EventNodeDataChanged:
				fmt.Println("EventNodeDataChanged", event.Path, event.Err, event.Server, event.State)

			case zk.EventNodeDeleted:
				fmt.Println("EventNodeDeleted", event.Path, event.Err, event.Server, event.State)

			}
		}
	}
}

func (s *ZkClient) GetConnsRemote(serviceName string, opts ...grpc.DialOption) (conns []*grpc.ClientConn, err error) {
	path := s.zkRoot + "/" + serviceName
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
		conn, err := grpc.Dial(string(data))
		if err == nil {
			conns = append(conns, conn)
		}
	}
	return conns, nil
}

func (s *ZkClient) GetConns(serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	conns, ok := s.rpcLocalCache[serviceName]
	if !ok {
		conns, err := s.GetConnsRemote(serviceName)
		if err != nil {
			return nil, err
		}
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
