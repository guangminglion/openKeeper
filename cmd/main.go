package cmd

import "openKeeper"

func main() {
	client, err := openKeeper.NewClient([]string{"43.154.157.177:9092"}, "openim", 10, "", "")
	if err != nil {
		panic(err.Error())
	}
	if err := client.Register("user", "127.0.0.1", 11000); err != nil {
		panic("")
	}
	client.GetConns("msg")
	client.GetConn("msg")
}
