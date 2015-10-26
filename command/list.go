package command

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {
	svc := ec2.New(&aws.Config{})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	var privateIPAddress, publicIPAddress string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"INSTANCE ID", "STATUS", "PRIVATE IP", "PUBLIC IP"}, "\t"))
	for idx, _ := range resp.Reservations {
		for _, instance := range resp.Reservations[idx].Instances {
			if instance.PrivateIPAddress != nil {
				privateIPAddress = *instance.PrivateIPAddress
			} else {
				privateIPAddress = ""
			}

			if instance.PublicIPAddress != nil {
				publicIPAddress = *instance.PublicIPAddress
			} else {
				publicIPAddress = ""
			}

			fmt.Fprintln(w, strings.Join([]string{*instance.InstanceID, *instance.State.Name, privateIPAddress, publicIPAddress}, "\t"))
		}
	}

	w.Flush()

	return 0
}

func (c *ListCommand) Synopsis() string {
	return "List EC2 instances"
}

func (c *ListCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
