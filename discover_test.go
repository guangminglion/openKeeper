package openKeeper

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestRegisterDiscover(t *testing.T) {
	root := "openim"
	addr := "43.154.157.177:2181"
	client, err := NewClient([]string{addr}, root, 1, "", "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	client2, err := NewClient([]string{addr}, root, 1, "", "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	client3, err := NewClient([]string{addr}, root, 1, "", "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	user := "user"
	host := "127.0.0.1"
	err = client.Register(user, host, 1001, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = client2.Register(user, host, 1002, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = client3.Register(user, host, 1003, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}

	conns, err := client.GetConns(user, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(conns) != 3 {
		t.Fatalf("error len: %d, %v", len(conns), conns)
	}

	err = client2.UnRegister()
	if err != nil {
		t.Fatalf(err.Error())
	}
	time.Sleep(time.Second * 2)
	conns, err = client.GetConns(user, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(conns) != 2 {
		t.Fatalf("error len: %d, %v", len(conns), conns)
	}
	err = client3.Register("msg", host, 9091, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	conns, err = client.GetConns("msg", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(conns) != 1 {
		t.Fatalf("error len: %d, %v", len(conns), conns)
	}
	defer func() {
		client2.UnRegister()
		client.UnRegister()
		client3.UnRegister()
	}()
}
