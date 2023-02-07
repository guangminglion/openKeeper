package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"openKeeper"
	"time"
)

func main() {
	client, err := openKeeper.NewClient([]string{"43.154.157.177:2181"}, "openim", 1, "", "")
	if err != nil {
		panic(err.Error())
	}
	client2, err := openKeeper.NewClient([]string{"43.154.157.177:2181"}, "openim", 1, "", "")
	if err != nil {
		panic(err.Error())
	}

	err = client2.UnRegister()
	if err != nil {

	}
	err = client.UnRegister()
	if err != nil {

	}

	err = client.Register("user", "127.0.0.1", 11000, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	conns, err := client.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	for _, conn := range conns {
		fmt.Println("1:", conn.Target())
	}

	err = client2.Register("msg", "127.0.0.1", 11001, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	conns, err = client2.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	for _, conn := range conns {
		fmt.Println("2:", conn.Target())
	}
	fmt.Println("client:", client.GetNode(), "client2:", client2.GetNode())

	conns, err = client.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	for _, conn := range conns {
		fmt.Println("user1:", conn.Target())
	}

	conns, err = client.GetConns("msg", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	for _, conn := range conns {
		fmt.Println("user2", conn.Target())
	}

	conns, err = client2.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	for _, conn := range conns {
		fmt.Println("GetConns user0:", conn.Target())
	}

	err = client.UnRegister()
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(time.Second * 3)

	conns, err = client2.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	for _, conn := range conns {
		fmt.Println("GetConns user1:", conn.Target())
	}

	err = client2.UnRegister()
	if err != nil {

	}
	err = client.UnRegister()
	if err != nil {

	}
}

func UnRegisterAll() {

}
