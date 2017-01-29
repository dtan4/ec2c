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

// launchCmd represents the launch command
var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch EC2 instance",
	RunE:  doLaunch,
}

var launchOpts = struct {
	amiID                    string
	associatePublicIPAddress bool
	availabilityZone         string
	dryRun                   bool
	instanceCount            int64
	instanceType             string
	keyName                  string
	securityGroups           []string
	subnetID                 string
	userData                 string
	volumeSize               int64
	volumeType               string
	instanceProfileARN       string
	roleName                 string
}{}

func doLaunch(cmd *cobra.Command, args []string) error {
	svc := ec2.New(session.New(), &aws.Config{})

	blockDeviceMappings := []*ec2.BlockDeviceMapping{
		&ec2.BlockDeviceMapping{
			DeviceName: aws.String("/dev/xvda"),
			Ebs: &ec2.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(true),
				VolumeSize:          aws.Int64(launchOpts.volumeSize),
				VolumeType:          aws.String(launchOpts.volumeType),
			},
		},
	}

	monitoring := &ec2.RunInstancesMonitoringEnabled{
		Enabled: aws.Bool(true),
	}

	opts := &ec2.RunInstancesInput{
		BlockDeviceMappings: blockDeviceMappings,
		DryRun:              aws.Bool(launchOpts.dryRun),
		ImageId:             aws.String(launchOpts.amiID),
		KeyName:             aws.String(launchOpts.keyName),
		InstanceType:        aws.String(launchOpts.instanceType),
		MaxCount:            aws.Int64(launchOpts.instanceCount),
		MinCount:            aws.Int64(1),
		Monitoring:          monitoring,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn:  aws.String(launchOpts.instanceProfileARN),
			Name: aws.String(launchOpts.roleName),
		},
	}

	if launchOpts.subnetID != "" && launchOpts.associatePublicIPAddress {
		opts.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			&ec2.InstanceNetworkInterfaceSpecification{
				AssociatePublicIpAddress: aws.Bool(launchOpts.associatePublicIPAddress),
				DeviceIndex:              aws.Int64(int64(0)),
				SubnetId:                 aws.String(launchOpts.subnetID),
				Groups:                   aws.StringSlice(launchOpts.securityGroups),
			},
		}
	}

	if launchOpts.userData != "" {
		buf, err := ioutil.ReadFile(launchOpts.userData)
		if err != nil {
			return errors.Wrapf(err, "failed to read user data file %q", launchOpts.userData)
		}

		opts.UserData = aws.String(base64.StdEncoding.EncodeToString(buf))
	}

	resp, err := svc.RunInstances(opts)
	if err != nil {
		return errors.Wrap(err, "failed to execute RunInstance")
	}

	for _, instance := range resp.Instances {
		fmt.Println(*instance.InstanceId)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(launchCmd)

	launchCmd.Flags().StringVar(&launchOpts.amiID, "ami", "", "AMI ID")
	launchCmd.Flags().BoolVar(&launchOpts.associatePublicIPAddress, "public", false, "Associate public IP address")
	launchCmd.Flags().StringVar(&launchOpts.availabilityZone, "az", "", "Availability zone")
	launchCmd.Flags().BoolVar(&launchOpts.dryRun, "dry-run", false, "Dry run")
	launchCmd.Flags().Int64Var(&launchOpts.instanceCount, "count", 1, "Number of instance")
	launchCmd.Flags().StringVar(&launchOpts.instanceType, "type", "", "Instance type")
	launchCmd.Flags().StringVar(&launchOpts.keyName, "key", "", "SSH key name")
	launchCmd.Flags().StringSliceVar(&launchOpts.securityGroups, "sg", []string{}, "Security Group IDs")
	launchCmd.Flags().StringVar(&launchOpts.subnetID, "subnet", "", "Subnet ID")
	launchCmd.Flags().StringVar(&launchOpts.userData, "user-data", "", "User data filename")
	launchCmd.Flags().Int64Var(&launchOpts.volumeSize, "volume-size", 8, "Volume size")
	launchCmd.Flags().StringVar(&launchOpts.volumeType, "volume-type", "gp2", "Volume type")
	launchCmd.Flags().StringVar(&launchOpts.instanceProfileARN, "instance-profile", "", "Instance Profile ARN")
	launchCmd.Flags().StringVar(&launchOpts.roleName, "iam-role", "", "IAM Role name")
}
