package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// untagCmd represents the untag command
var untagCmd = &cobra.Command{
	Use:   "untag",
	Short: "Remove instance tags",
	RunE:  doUntag,
}

var untagOpts = struct {
	dryRun      bool
	instanceIDs []string
	tagStrings  []string
}{}

func doUntag(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	tags := []*ec2.Tag{}

	for _, tagString := range untagOpts.tagStrings {
		tags = append(tags, &ec2.Tag{
			Key: aws.String(tagString),
		})
	}

	opts := &ec2.DeleteTagsInput{
		DryRun:    aws.Bool(untagOpts.dryRun),
		Resources: aws.StringSlice(untagOpts.instanceIDs),
		Tags:      tags,
	}

	_, err := svc.DeleteTags(opts)
	if err != nil {
		return errors.Wrap(err, "failed to execute DeleteTags")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(untagCmd)

	tagCmd.Flags().BoolVar(&untagOpts.dryRun, "dry-run", false, "Dry run")
	tagCmd.Flags().StringSliceVar(&untagOpts.instanceIDs, "instances", []string{}, "Instance IDs")
	tagCmd.Flags().StringSliceVar(&untagOpts.tagStrings, "tags", []string{}, "KEY=value tags")
}
