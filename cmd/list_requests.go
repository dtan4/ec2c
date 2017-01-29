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

// listRequestsCmd represents the list_requests command
var listRequestsCmd = &cobra.Command{
	Use:   "list-requests",
	Short: "List spot instance requests",
	RunE:  doListRequests,
}

func doListRequests(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	resp, err := svc.DescribeSpotInstanceRequests(nil)
	if err != nil {
		return errors.Wrap(err, "failed to execute DescribeSpotInstanceRequests")
	}

	var instanceID, requestName string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"REQUEST ID", "MAX PRICE", "INSTANCE TYPE", "INSTANCE ID", "STATE", "STATUS", "NAME"}, "\t"))
	for _, request := range resp.SpotInstanceRequests {
		requestName = ""

		if request.InstanceId != nil {
			instanceID = *request.InstanceId
		} else {
			instanceID = ""
		}

		for _, tag := range request.Tags {
			if *tag.Key == "Name" {
				requestName = *tag.Value
				break
			}
		}

		fmt.Fprintln(w, strings.Join(
			[]string{
				*request.SpotInstanceRequestId,
				*request.SpotPrice,
				*request.LaunchSpecification.InstanceType,
				instanceID,
				*request.State,
				*request.Status.Code,
				requestName,
			}, "\t",
		))
	}

	w.Flush()

	return nil
}

func init() {
	RootCmd.AddCommand(listRequestsCmd)
}
