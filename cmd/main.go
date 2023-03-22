package main

import (
	"fmt"
	"time"

	"github.com/OpenIMSDK/openKeeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, err := openKeeper.NewClient([]string{"43.154.157.177:2181"}, "openim", 10, "", "")
	if err != nil {
		panic(err.Error())
	}
	//conns, err := client.GetConns("User", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(conns)
	//for _, v := range conns {
	//	fmt.Println(v.Target())
	//}

	conns, err := client.GetConns("MessageGateway", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(conns[0].Target())

	//client2, err := openKeeper.NewClient([]string{"43.154.157.177:2181"}, "openim", 1, "", "")
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//client3, err := openKeeper.NewClient([]string{"43.154.157.177:2181"}, "openim", 1, "", "")
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//err = client.UnRegister()
	//if err != nil {
	//
	//}
	//err = client2.UnRegister()
	//if err != nil {
	//
	//}
	//err = client3.UnRegister()
	//if err != nil {
	//
	//}
	//

	//for _, conn := range conns {
	//	fmt.Println("1:", conn.Target())
	//}
	//
	//err = client2.Register("msg", "127.0.0.1", 11001, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//conns, err = client2.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//for _, conn := range conns {
	//	fmt.Println("2:", conn.Target())
	//}
	//fmt.Println("client:", client.GetNode(), "client2:", client2.GetNode())
	//
	//err = client3.Register("msg", "127.0.0.1", 11001, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//conns, err = client.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//for _, conn := range conns {
	//	fmt.Println("user:", conn.Target())
	//}
	//
	//conns, err = client.GetConns("msg", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//for _, conn := range conns {
	//	fmt.Println("msg", conn.Target())
	//}
	//
	//err = client.UnRegister()
	//if err != nil {
	//	panic(err.Error())
	//}
	//time.Sleep(time.Second * 3)
	//
	//conns, err = client2.GetConns("user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//for _, conn := range conns {
	//	fmt.Println("after unregister:", conn.Target())
	//}
	//
	//conns, err = client2.GetConns("msg", grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	panic(err.Error())
	//}
	//for _, conn := range conns {
	//	fmt.Println("GetConns user1:", conn.Target())
	//}
	//
	//err = client2.UnRegister()
	//if err != nil {
	//
	//}
	//err = client.UnRegister()
	//if err != nil {
	//
	//}

	for {
		time.Sleep(time.Second * 10)
	}
}

func UnRegisterAll() {

}
