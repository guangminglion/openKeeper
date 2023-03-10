package openKeeper

import (
	"errors"
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"strings"
)

var ErrConnIsNil = errors.New("conn is nil")

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
				if len(l) > 1 {
					delete(s.localConns, l[len(l)-1])
				}
				s.lock.Unlock()
			case zk.EventNodeDataChanged:
			case zk.EventNodeDeleted:
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
	s.lock.RLock()
	conns, ok := s.localConns[serviceName]
	if !ok {
		s.lock.RUnlock()
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
		s.localConns[serviceName] = conns
	} else {
		s.lock.RUnlock()
	}
	if len(conns) == 0 {
		return nil, ErrConnIsNil
	}
	return conns, nil
}

func (s *ZkClient) GetConn(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	robin := Robin{}
	return s.GetConnStrategy(serviceName, robin.Robin, append(s.options, opts...)...)
}

func (s *ZkClient) GetConnStrategy(serviceName string, strategy func(slice []*grpc.ClientConn) int, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	conns, err := s.GetConns(serviceName, opts...)
	if len(conns) > 0 {
		return conns[strategy(conns)], nil
	}
	return nil, err
}

type Robin struct {
	next int
}

func (r *Robin) Robin(slice []*grpc.ClientConn) int {
	index := r.next
	r.next += 1
	if r.next > len(slice)-1 {
		r.next = 0
	}
	return index
}
