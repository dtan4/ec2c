package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type TagCommand struct {
	Meta
}

func (c *TagCommand) Run(args []string) int {
	var (
		dryRun           bool
		instanceIdString string
		instanceIds      []*string
		tagString        string
		tags             []*ec2.Tag
	)

	var (
		arguments []string
	)

	svc := ec2.New(&aws.Config{})

	flags := flag.NewFlagSet("dtan4", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&dryRun, "dry-run", false, "Dry run (default: false)")
	flags.StringVar(&instanceIdString, "instance", "", "Instance Ids")
	flags.StringVar(&tagString, "tags", "", "KEY=VALUE tags")

	if err := flags.Parse(args[0:]); err != nil {
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	for _, id := range strings.Split(instanceIdString, ",") {
		instanceIds = append(instanceIds, aws.String(id))
	}

	for _, tag := range strings.Split(tagString, ",") {
		keyValue := strings.Split(tag, "=")

		if len(keyValue) >= 2 {
			tags = append(tags, &ec2.Tag{
				Key:   aws.String(keyValue[0]),
				Value: aws.String(strings.Join(keyValue[1:], "=")),
			})
		}
	}

	opts := &ec2.CreateTagsInput{
		DryRun:    aws.Bool(dryRun),
		Resources: instanceIds,
		Tags:      tags,
	}

	resp, err := svc.CreateTags(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	return 0
}

func (c *TagCommand) Synopsis() string {
	return "Tag EC2 instances"
}

func (c *TagCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
