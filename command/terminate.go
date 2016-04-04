package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type TerminateCommand struct {
	Meta
}

func (c *TerminateCommand) Run(args []string) int {
	var (
		dryRun           bool
		instanceIdstring string
		instanceIds      []*string
	)

	var (
		arguments []string
	)

	svc := ec2.New(session.New(), &aws.Config{})

	flags := flag.NewFlagSet("dtan4", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&dryRun, "dry-run", false, "Dry run (default: false)")
	flags.StringVar(&instanceIdstring, "instance", "", "Instance IDs")

	if err := flags.Parse(args[0:]); err != nil {
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	for _, id := range strings.Split(instanceIdstring, ",") {
		instanceIds = append(instanceIds, aws.String(id))
	}

	opts := &ec2.TerminateInstancesInput{
		DryRun:      aws.Bool(dryRun),
		InstanceIds: instanceIds,
	}

	resp, err := svc.TerminateInstances(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	return 0
}

func (c *TerminateCommand) Synopsis() string {
	return "Terminate the specified EC2 instance"
}

func (c *TerminateCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
