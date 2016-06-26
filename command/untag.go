package command

import (
	"flag"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type UntagCommand struct {
	Meta
}

func (c *UntagCommand) Run(args []string) int {
	var (
		dryRun           bool
		instanceIDString string
		instanceIDs      []*string
		tagString        string
		tags             []*ec2.Tag
	)

	var (
		arguments []string
	)

	svc := ec2.New(session.New(), &aws.Config{})

	flags := flag.NewFlagSet("dtan4", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&dryRun, "dry-run", false, "Dry run (default: false)")
	flags.StringVar(&instanceIDString, "instance", "", "Instance Ids")
	flags.StringVar(&tagString, "tags", "", "KEY tags")

	if err := flags.Parse(args[0:]); err != nil {
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	for _, id := range strings.Split(instanceIDString, ",") {
		instanceIDs = append(instanceIDs, aws.String(id))
	}

	for _, tag := range strings.Split(tagString, ",") {
		tags = append(tags, &ec2.Tag{
			Key: aws.String(tag),
		})
	}

	opts := &ec2.DeleteTagsInput{
		DryRun:    aws.Bool(dryRun),
		Resources: instanceIDs,
		Tags:      tags,
	}

	_, err := svc.DeleteTags(opts)
	if err != nil {
		panic(err)
	}

	return 0
}

func (c *UntagCommand) Synopsis() string {
	return "Tag EC2 instances"
}

func (c *UntagCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
