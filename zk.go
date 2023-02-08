package openKeeper

import (
	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
	"time"
)

type ZkClient struct {
	zkServers []string
	zkRoot    string

	conn      *zk.Conn
	eventChan <-chan zk.Event
	node      string

	lock          sync.Mutex
	rpcLocalCache map[string][]*grpc.ClientConn
}

func NewClient(zkServers []string, zkRoot string, timeout int, userName, password string) (*ZkClient, error) {
	client := &ZkClient{
		zkServers:     zkServers,
		zkRoot:        "/",
		rpcLocalCache: make(map[string][]*grpc.ClientConn, 0),
	}
	conn, eventChan, err := zk.Connect(zkServers, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	client.zkRoot += zkRoot
	client.eventChan = eventChan
	client.conn = conn
	if err := client.ensureRoot(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

func (s *ZkClient) Close() {
	s.conn.Close()
}

func (s *ZkClient) ensureAndCreate(node string) error {
	exists, _, err := s.conn.Exists(s.zkRoot)
	if err != nil {
		return err
	}
	if !exists {
		_, err := s.conn.Create(node, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *ZkClient) GetZkConn() *zk.Conn {
	return s.conn
}

func (s *ZkClient) GetRootPath() string {
	return s.zkRoot
}

func (s *ZkClient) GetNode() string {
	return s.node
}

func (s *ZkClient) ensureRoot() error {
	return s.ensureAndCreate(s.zkRoot)
}

func (s *ZkClient) ensureName(name string) error {
	path := s.zkRoot + "/" + name
	return s.ensureAndCreate(path)
}

func (s *ZkClient) getPath(rpcRegisterName string) string {
	return s.zkRoot + "/" + rpcRegisterName
}

func (s *ZkClient) getAddr(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
