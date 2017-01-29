package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// requestCmd represents the request command
var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Request spot instance",
	RunE:  doRequest,
}

var requestOpts = struct {
	amiID                    string
	associatePublicIPAddress bool
	availabilityZone         string
	dryRun                   bool
	instanceCount            int64
	instanceType             string
	keyName                  string
	securityGroups           []string
	spotPrice                string
	subnetID                 string
	userData                 string
	volumeSize               int64
	volumeType               string
}{}

func doRequest(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	blockDeviceMappings := []*ec2.BlockDeviceMapping{
		&ec2.BlockDeviceMapping{
			DeviceName: aws.String("/dev/xvda"),
			Ebs: &ec2.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(true),
				VolumeSize:          aws.Int64(requestOpts.volumeSize),
				VolumeType:          aws.String(requestOpts.volumeType),
			},
		},
	}

	launchSpecification := &ec2.RequestSpotLaunchSpecification{
		BlockDeviceMappings: blockDeviceMappings,
		ImageId:             aws.String(requestOpts.amiID),
		KeyName:             aws.String(requestOpts.keyName),
		InstanceType:        aws.String(requestOpts.instanceType),
		Monitoring: &ec2.RunInstancesMonitoringEnabled{
			Enabled: aws.Bool(true),
		},
		SecurityGroups: aws.StringSlice(requestOpts.securityGroups),
		SubnetId:       aws.String(requestOpts.subnetID),
	}

	if requestOpts.subnetID != "" && requestOpts.associatePublicIPAddress {
		launchSpecification.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			&ec2.InstanceNetworkInterfaceSpecification{
				AssociatePublicIpAddress: aws.Bool(requestOpts.associatePublicIPAddress),
				DeviceIndex:              aws.Int64(int64(0)),
				SubnetId:                 aws.String(requestOpts.subnetID),
				Groups:                   aws.StringSlice(requestOpts.securityGroups),
			},
		}
	}

	if requestOpts.userData != "" {
		buf, err := ioutil.ReadFile(requestOpts.userData)
		if err != nil {
			return errors.Wrapf(err, "failed to open user data file %q", requestOpts.userData)
		}

		launchSpecification.UserData = aws.String(base64.StdEncoding.EncodeToString(buf))
	}

	opts := &ec2.RequestSpotInstancesInput{
		DryRun:              aws.Bool(requestOpts.dryRun),
		InstanceCount:       aws.Int64(requestOpts.instanceCount),
		LaunchSpecification: launchSpecification,
		SpotPrice:           aws.String(requestOpts.spotPrice),
	}

	resp, err := svc.RequestSpotInstances(opts)
	if err != nil {
		return errors.Wrap(err, "failed to execute RequestSpotInstances")
	}

	fmt.Println(resp)

	return nil
}

func init() {
	RootCmd.AddCommand(requestCmd)

	requestCmd.Flags().StringVar(&requestOpts.amiID, "ami", "", "AMI ID")
	requestCmd.Flags().BoolVar(&requestOpts.associatePublicIPAddress, "public", false, "Associate public IP address")
	requestCmd.Flags().StringVar(&requestOpts.availabilityZone, "az", "", "Availability zone")
	requestCmd.Flags().BoolVar(&requestOpts.dryRun, "dry-run", false, "Dry run")
	requestCmd.Flags().Int64Var(&requestOpts.instanceCount, "count", 1, "Number of instance")
	requestCmd.Flags().StringVar(&requestOpts.instanceType, "type", "", "Instance type")
	requestCmd.Flags().StringVar(&requestOpts.keyName, "key", "", "SSH key name")
	requestCmd.Flags().StringSliceVar(&requestOpts.securityGroups, "sg", []string{}, "Security Group IDs")
	requestCmd.Flags().StringVar(&requestOpts.spotPrice, "spot-price", "", "Maximum bidding price")
	requestCmd.Flags().StringVar(&requestOpts.subnetID, "subnet", "", "Subnet ID")
	requestCmd.Flags().StringVar(&requestOpts.userData, "user-data", "", "User data filename")
	requestCmd.Flags().Int64Var(&requestOpts.volumeSize, "volume-size", 8, "Volume size")
	requestCmd.Flags().StringVar(&requestOpts.volumeType, "volume-type", "gp2", "Volume type")
}
