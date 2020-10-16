// Package osc provides node discovery for Outscale API.
package osc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/antihax/optional"
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
	vmId := args["vm_id"]
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
	options := osc.ReadVmsRequest{Filters: osc.FiltersVm{}}
	l.Printf("[INFO] discover-osc: Filter instances with %s=%s", tagKey, tagValue)
	if tagKey != "" && tagValue != "" {
		options = osc.ReadVmsRequest{
			Filters: osc.FiltersVm{
				TagKeys: []string{tagKey},
				TagValues: []string{tagValue},
			},
		};
	} else if vmId != "" {
		options = osc.ReadVmsRequest{
			Filters: osc.FiltersVm{
				VmIds: []string{vmId},
			},
		};
	} else {
		options = osc.ReadVmsRequest{Filters: osc.FiltersVm{}}
	}
	readOpts := osc.ReadVmsOpts{ReadVmsRequest: optional.NewInterface(options)};
	read, _, err := client.VmApi.ReadVms(auth, &readOpts)
	if err != nil {
		return nil, fmt.Errorf("discover-osc: ReadVms failed: %s", err)
	}

	l.Printf("[DEBUG] discover-osc: Found %d reservations", len(read.Vms))
	var addrs []string
	for _, vm := range read.Vms {
		l.Printf("[DEBUG] discover-osc: Vm %s", vm.VmId)
		switch addrType {
		case "private_ip":
			if vm.PrivateIp != "" {
				l.Printf("[INFO] discover-osc: Private Ip %s on Vm %s", vm.PrivateIp, vm.VmId)
				addrs = append(addrs, vm.PrivateIp)
			}
			break;
		case "public_ip":
			if vm.PublicIp != "" {
				l.Printf("[INFO] discover-osc: Vm %s has public ip %s", vm.VmId, vm.PublicIp)
				addrs = append(addrs, vm.PublicIp)
			}
			break;
		default:
			if vm.PrivateIp != "" {
				l.Printf("[INFO] discover-osc: Private Ip %s on Vm %s", vm.PrivateIp, vm.VmId)
				addrs = append(addrs, vm.PrivateIp)
			}
			if vm.PublicIp != "" {
				l.Printf("[INFO] discover-osc: Vm %s has public ip %s", vm.VmId, vm.PublicIp)
				addrs = append(addrs, vm.PublicIp)
			}
			break;
		}
	}
	l.Printf("[DEBUG] discover-osc: Found ip addresses: %v", addrs)
	return addrs, nil
}
