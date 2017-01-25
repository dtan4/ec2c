package command

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dtan4/ec2c/msg"
)

type ListCommand struct {
	Meta
}

func (c *ListCommand) Run(args []string) int {
	var (
		tags bool
	)

	arguments := []string{}

	svc := ec2.New(session.New(), &aws.Config{})

	flags := flag.NewFlagSet("dtan4", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&tags, "tags", false, "Print instance tags")

	if err := flags.Parse(args[0:]); err != nil {
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		msg.Errorf("Failed to retrieve instance list. error: %s\n", err)
		return 1
	}

	var privateIPAddress, publicIPAddress, instanceName string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	if tags {
		fmt.Fprintln(w, strings.Join([]string{
			"INSTANCE ID",
			"STATUS",
			"INSTANCE TYPE",
			"AVAILABILITY ZONE",
			"PRIVATE IP",
			"PUBLIC IP",
			"NAME",
			"TAG",
		}, "\t"))
	} else {
		fmt.Fprintln(w, strings.Join([]string{
			"INSTANCE ID",
			"STATUS",
			"INSTANCE TYPE",
			"AVAILABILITY ZONE",
			"PRIVATE IP",
			"PUBLIC IP",
			"NAME",
		}, "\t"))
	}

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

			if tags {
				fmt.Fprintln(w, strings.Join(
					[]string{
						*instance.InstanceId,
						*instance.State.Name,
						*instance.InstanceType,
						*instance.Placement.AvailabilityZone,
						privateIPAddress,
						publicIPAddress,
						instanceName,
						strings.Join(tagKeyValue, ","),
					}, "\t",
				))
			} else {
				fmt.Fprintln(w, strings.Join(
					[]string{
						*instance.InstanceId,
						*instance.State.Name,
						*instance.InstanceType,
						*instance.Placement.AvailabilityZone,
						privateIPAddress,
						publicIPAddress,
						instanceName,
					}, "\t",
				))
			}
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
