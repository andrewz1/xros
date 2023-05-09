package xros_test

import (
	"testing"

	"github.com/andrewz1/xnet"

	"github.com/andrewz1/xros"
)

func getClient() (cl *xros.Client, err error) {
	cn, err := xnet.Dial("tcp", "10.0.11.229:8728")
	if err != nil {
		return nil, err
	}
	cl = xros.NewClient(cn)
	defer func() {
		if err != nil {
			cl.Close()
		}
	}()
	err = cl.Login("admin", "1q2w3e4r")
	return
}

func TestConn(t *testing.T) {
	cl, err := getClient()
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()
	ss, err := cl.Run("/interface/print", "?name=wlan2")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range ss {
		t.Logf("%s\n", s)
	}
}

func TestDeAuth(t *testing.T) {
	cl, err := getClient()
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()
	//, "=.proplist=.id"
	// , "?mac-address=04:4F:4C:66:A5:14"
	ss, err := cl.Run("/ip/hotspot/host/getall", "?mac-address=04:4F:4C:66:A5:14", "=.proplist=.id")
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range ss {
		t.Logf("%s\n", s)
	}
	if len(ss) != 1 {
		return
	}
	id, ok := ss[0].Data[".id"]
	if !ok {
		return
	}
	_, err = cl.Run("/ip/hotspot/host/remove", "=.id="+id)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuth(t *testing.T) {
	cl, err := getClient()
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()

	if _, err = cl.Run("/ip/hotspot/active/login", "=ip=10.0.12.251", "=user=1C:91:80:C6:85:10", "=mac-address=1C:91:80:C6:85:10", "=password=mikrotik"); err != nil {
		t.Fatal(err)
	}

}
