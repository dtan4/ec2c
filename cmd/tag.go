package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Add tags to EC2 instance",
	RunE:  doTag,
}

var tagOpts = struct {
	dryRun      bool
	instanceIDs []string
	tagStrings  []string
}{}

func doTag(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	tags := []*ec2.Tag{}

	for _, tagString := range tagOpts.tagStrings {
		keyValue := strings.Split(tagString, "=")

		if len(keyValue) == 1 {
			tags = append(tags, &ec2.Tag{
				Key:   aws.String(keyValue[0]),
				Value: aws.String(""),
			})
		} else {
			tags = append(tags, &ec2.Tag{
				Key:   aws.String(keyValue[0]),
				Value: aws.String(strings.Join(keyValue[1:], "=")),
			})
		}
	}

	opts := &ec2.CreateTagsInput{
		DryRun:    aws.Bool(tagOpts.dryRun),
		Resources: aws.StringSlice(tagOpts.instanceIDs),
		Tags:      tags,
	}

	_, err := svc.CreateTags(opts)
	if err != nil {
		return errors.Wrap(err, "failed to execute CreateTags")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(tagCmd)

	tagCmd.Flags().BoolVar(&tagOpts.dryRun, "dry-run", false, "Dry run")
	tagCmd.Flags().StringSliceVar(&tagOpts.instanceIDs, "instances", []string{}, "Instance IDs")
	tagCmd.Flags().StringSliceVar(&tagOpts.tagStrings, "tags", []string{}, "KEY=value tags")
}
