package osc_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/osc"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "osc",
		"region":            os.Getenv("OSC_REGION"),
		"tag_key":           os.Getenv("TAG_KEY"),
		"tag_value":         os.Getenv("TAG_VALUE"),
		"addr_type":         os.Getenv("ADDR_TYPE"),
		"access_key_id":     os.Getenv("OSC_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("OSC_ACCESS_KEY_SECRET"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("OSC credentials or region missing")
	}

	p := &osc.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	l.Printf("addrs : %s", addrs)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) == 0 {
		t.Fatalf("bad: %v", addrs)
	}
}
