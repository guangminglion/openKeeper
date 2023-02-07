package openKeeper

import "github.com/go-zookeeper/zk"

const ConfName = "openIMConf"

func (s *ZkClient) RegisterConf(conf []byte) error {
	_, err := s.conn.Create(s.getPath(ConfName), conf, 0, zk.WorldACL(zk.PermAll))
	return err
}

func (s *ZkClient) LoadConf() ([]byte, error) {
	bytes, _, err := s.conn.Get(s.getPath(ConfName))
	return bytes, err
}
