package command

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type ListRequestsCommand struct {
	Meta
}

func (c *ListRequestsCommand) Run(args []string) int {
	svc := ec2.New(&aws.Config{})

	resp, err := svc.DescribeSpotInstanceRequests(nil)
	if err != nil {
		panic(err)
	}

	var instanceId, requestName string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	fmt.Fprintln(w, strings.Join([]string{"REQUEST ID", "MAX PRICE", "INSTANCE TYPE", "INSTANCE ID", "STATE", "STATUS", "NAME"}, "\t"))
	for _, request := range resp.SpotInstanceRequests {
		requestName = ""

		if request.InstanceId != nil {
			instanceId = *request.InstanceId
		} else {
			instanceId = ""
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
				instanceId,
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
