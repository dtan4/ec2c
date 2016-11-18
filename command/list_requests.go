package command

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dtan4/ec2c/msg"
)

type ListRequestsCommand struct {
	Meta
}

func (c *ListRequestsCommand) Run(args []string) int {
	svc := ec2.New(session.New(), &aws.Config{})

	resp, err := svc.DescribeSpotInstanceRequests(nil)
	if err != nil {
		msg.Errorf("Failed to retrieve SpotRequest list. error: %s\n", err)
		return 1
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

	return 0
}

func (c *ListRequestsCommand) Synopsis() string {
	return "List Spot Instance requests"
}

func (c *ListRequestsCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
