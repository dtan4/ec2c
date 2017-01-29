package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel spot instance request",
	RunE:  doCancel,
}

var cancelOpts = struct {
	dryRun     bool
	requestIDs []string
}{}

func doCancel(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	resp, err := svc.CancelSpotInstanceRequests(&ec2.CancelSpotInstanceRequestsInput{
		DryRun:                 aws.Bool(cancelOpts.dryRun),
		SpotInstanceRequestIds: aws.StringSlice(cancelOpts.requestIDs),
	})
	if err != nil {
		return errors.Wrap(err, "failed to execute CancelSpotInstanceRequests")
	}

	fmt.Println(resp)

	return nil
}

func init() {
	RootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().BoolVar(&cancelOpts.dryRun, "dry-run", false, "dry-run")
	cancelCmd.Flags().StringSliceVar(&cancelOpts.requestIDs, "request", []string{}, "Spot instance request IDs")
}
