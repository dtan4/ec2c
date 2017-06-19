package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List EC2 instances",
	RunE:  doRun,
}

var listOpts = struct {
	launchTime bool
	tags       bool
}{}

func doRun(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		return errors.Wrap(err, "failed to execute DescribeInstances")
	}

	var privateIPAddress, publicIPAddress, instanceName string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	headers := []string{
		"INSTANCE ID",
		"STATUS",
		"INSTANCE TYPE",
		"AVAILABILITY ZONE",
		"PRIVATE IP",
		"PUBLIC IP",
	}

	if listOpts.launchTime {
		headers = append(headers, "LAUNCH TIME")
	}

	headers = append(headers, "NAME")

	if listOpts.tags {
		headers = append(headers, "TAGS")
	}

	fmt.Fprintln(w, strings.Join(headers, "\t"))

	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
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

			fields := []string{
				*instance.InstanceId,
				*instance.State.Name,
				*instance.InstanceType,
				*instance.Placement.AvailabilityZone,
				privateIPAddress,
				publicIPAddress,
			}

			if listOpts.launchTime {
				fields = append(fields, (*instance.LaunchTime).Local().String())
			}

			fields = append(fields, instanceName)

			if listOpts.tags {
				fields = append(fields, strings.Join(tagKeyValue, ","))
			}

			fmt.Fprintln(w, strings.Join(fields, "\t"))
		}
	}

	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVar(&listOpts.launchTime, "launch-time", false, "Print instance launch time")
	listCmd.Flags().BoolVar(&listOpts.tags, "tags", false, "Print instance tags")
}
