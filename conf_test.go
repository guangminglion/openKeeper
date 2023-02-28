package openKeeper

import (
	"testing"
)

func TestRegisterConf(t *testing.T) {
	client, err := NewClient([]string{"43.154.157.177:2181"}, "openim", 1, "", "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	s := "test"
	key := "conf"
	err = client.RegisterConf2Registry(key, []byte(s))
	if err != nil {
		t.Fatalf(err.Error())
	}
	conf, err := client.GetConfFromRegistry(key)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if string(conf) != s {
		t.Fatalf("not the same result")
	}
}
