package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// terminateCmd represents the terminate command
var terminateCmd = &cobra.Command{
	Use:   "terminate",
	Short: "Terminate EC2 instances",
	RunE:  doTerminate,
}

var terminateOpts = struct {
	dryRun      bool
	instanceIDs []string
}{}

func doTerminate(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	opts := &ec2.TerminateInstancesInput{
		DryRun:      aws.Bool(terminateOpts.dryRun),
		InstanceIds: aws.StringSlice(terminateOpts.instanceIDs),
	}

	resp, err := svc.TerminateInstances(opts)
	if err != nil {
		return errors.Wrap(err, "failed to execute TerminateInstances")
	}

	for _, instance := range resp.TerminatingInstances {
		fmt.Println(*instance.InstanceId)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(terminateCmd)

	terminateCmd.Flags().BoolVar(&terminateOpts.dryRun, "dry-run", false, "Dry run")
	terminateCmd.Flags().StringSliceVar(&terminateOpts.instanceIDs, "instances", []string{}, "Instance IDs")
}
