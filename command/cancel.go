package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type CancelCommand struct {
	Meta
}

func (c *CancelCommand) Run(args []string) int {
	var (
		dryRun          bool
		requestIdstring string
		requestIds      []*string
	)

	var (
		arguments []string
	)

	svc := ec2.New(session.New(), &aws.Config{})

	flags := flag.NewFlagSet("dtan4", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&dryRun, "dry-run", false, "Dry run (default: false)")
	flags.StringVar(&requestIdstring, "request", "", "Spot Instance request IDs")

	if err := flags.Parse(args[0:]); err != nil {
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	for _, id := range strings.Split(requestIdstring, ",") {
		requestIds = append(requestIds, aws.String(id))
	}

	opts := &ec2.CancelSpotInstanceRequestsInput{
		DryRun:                 aws.Bool(dryRun),
		SpotInstanceRequestIds: requestIds,
	}

	resp, err := svc.CancelSpotInstanceRequests(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	return 0
}

func (c *CancelCommand) Synopsis() string {
	return "Cancel the specified EC2 Spot Instance requests"
}

func (c *CancelCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
