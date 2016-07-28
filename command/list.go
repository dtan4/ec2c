package command

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {
	svc := ec2.New(session.New(), &aws.Config{})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	var privateIPAddress, publicIPAddress, instanceName string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"INSTANCE ID", "STATUS", "INSTANCE TYPE", "PRIVATE IP", "PUBLIC IP", "NAME", "TAG"}, "\t"))
	for idx, _ := range resp.Reservations {
		for _, instance := range resp.Reservations[idx].Instances {
			if instance.PrivateIpAddress != nil {
				privateIPAddress = *instance.PrivateIpAddress
			} else {
				privateIPAddress = ""
			}

			if instance.PublicIpAddress != nil {
				publicIPAddress = *instance.PublicIpAddress
			} else {
				publicIPAddress = ""
			}

			instanceName = ""
			tagKeyValue := []string{}

			for _, tag := range instance.Tags {
				keyValue := *tag.Key
				if *tag.Key == "Name" {
					instanceName = *tag.Value
				}
				if len(*tag.Value) > 0 {
					keyValue += "=" + *tag.Value
				}
				tagKeyValue = append(tagKeyValue, keyValue)
			}

			fmt.Fprintln(w, strings.Join(
				[]string{*instance.InstanceId, *instance.State.Name, *instance.InstanceType, privateIPAddress, publicIPAddress, instanceName, strings.Join(tagKeyValue, ",")}, "\t",
			))
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
