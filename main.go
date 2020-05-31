package main

import (
	"errors"

	"github.com/pulumi/pulumi-gcp/sdk/v3/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a network
		network, err := compute.NewNetwork(ctx, "network", &compute.NetworkArgs{})
		if err != nil {
			return err
		}

		_, err = compute.NewFirewall(ctx, "firewall", &compute.FirewallArgs{
			Network: network.ID(),
			Allows: compute.FirewallAllowArray{
				compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.StringArray{
						pulumi.String("22"),
						pulumi.String("80"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		const startupScript = `#!/usr/bin/env bash
echo "Hello, World!" > index.html
nohup python -m SimpleHTTPServer 80 &`

		// Create a Virtual Machine Instance
		computeInstance, err := compute.NewInstance(ctx, "instance", &compute.InstanceArgs{
			MachineType:           pulumi.String("f1-micro"),
			Zone:                  pulumi.String("us-east1-d"),
			MetadataStartupScript: pulumi.String(startupScript),
			BootDisk: compute.InstanceBootDiskArgs{
				InitializeParams: compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String("debian-cloud/debian-9"),
				},
			},
			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				compute.InstanceNetworkInterfaceArgs{
					Network: network.ID(),

					// AccessConfigs must include a single empty config to request an ephemeral IP
					AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
						compute.InstanceNetworkInterfaceAccessConfigArgs{},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the name and IP address of the Instance
		ctx.Export("instanceName", computeInstance.Name)
		ctx.Export("instanceIP", computeInstance.NetworkInterfaces.Apply(func(input interface{}) (interface{}, error) {
			xini, ok := input.([]compute.InstanceNetworkInterface)
			if !ok {
				return nil, errors.New(("not OK!"))
			}

			if len(xini) == 0 {
				return nil, errors.New(("no instance network interfaces!"))
			}

			if len(xini[0].AccessConfigs) == 0 {
				return nil, errors.New(("no instance network interface access configs!"))
			}

			return xini[0].AccessConfigs[0].NatIp, nil
		}))

		return nil
	})
}
