// Package osc provides node discovery for Amazon AWS.
package osc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/outscale/osc-sdk-go/osc"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Outscale OSC:

    provider:          "osc"
    region:            The OSC region. Default to region of instance.
    tag_key:           The tag key to filter on
    tag_value:         The tag value to filter on
    access_key_id:     The OSC access key to use
    secret_access_key: The OSC secret access key to use
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "osc" {
		return nil, fmt.Errorf("discover-osc: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	addrType := args["addr_type"]
	accessKey := args["access_key_id"]
	secretKey := args["secret_access_key"]

	l.Printf("[DEBUG] discover-osc: Using region=%s tag_key=%s tag_value=%s", region, tagKey, tagValue)
	if accessKey == "" && secretKey == "" {
		l.Printf("[DEBUG] discover-osc: No static credentials")
		l.Printf("[DEBUG] discover-osc: Using environment variables, shared credentials or instance role")
	} else {
		l.Printf("[DEBUG] discover-osc: Static credentials provided")
	}

	l.Printf("[INFO] discover-osc: Region is %s", region)
	client := osc.NewAPIClient(osc.NewConfiguration())
	l.Printf("[DEBUG] discover-osc: Creating session...")
	auth := context.WithValue(context.Background(), osc.ContextAWSv4, osc.AWSv4{
		AccessKey: accessKey,
		SecretKey: secretKey,
	})

	l.Printf("[INFO] discover-aws: Filter instances with %s=%s", tagKey, tagValue)
	// NEED TO REPLACE THAT WITH AN OSC API CALL "/ReadVms/"
	read, _, err := client.VmApi.ReadVms(auth, nil)
	if err != nil {
		return nil, fmt.Errorf("discover-osc: ReadVms failed: %s", err)
	}

	l.Printf("[DEBUG] discover-aws: Found %d reservations", len(read.Vms))
	var addrs []string
	for _, vm := range read.Vms {
		l.Printf("------------------------\nVmId: %s", vm.VmId)
		// l.Printf("[DEBUG] discover-osc: Vm %s", vm.VmId)
		if addrType == addrType {

			l.Printf("PrivateIp : %v", vm.PrivateIp)
			l.Printf("PublicIp : %v", vm.PrivateIp)
			l.Printf("AddrType: %s", addrType)
		}

		// switch addrType {
		// case "private_ip":
		if vm.PrivateIp != "" {
			l.Printf("[INFO] discover-osc: Private Ip %s on Vm %s", vm.PrivateIp, vm.VmId)
			addrs = append(addrs, vm.PrivateIp)
		}
		// case "public_ip":
		if vm.PublicIp != "" {
			l.Printf("[INFO] discover-osc: Vm %s has public ip %s", vm.VmId, vm.PublicIp)
			addrs = append(addrs, vm.PublicIp)
		}
		// default:
		// // EC2-Classic don't have the PrivateIpAddress field
		// 	if vm.PrivateIpAddress == nil {
		// 		l.Printf("[DEBUG] discover-aws: Instance %s has no private ip", id)
		// 		continue
		// }
		// l.Printf("[INFO] discover-aws: Instance %s has private ip %s", id, *inst.PrivateIpAddress)
		// addrs = append(addrs, vinst.PrivateIpAddress)
		l.Printf("------------------------\n")

	}
	l.Printf("[DEBUG] discover-osc: Found ip addresses: %v", addrs)
	return addrs, nil
}
