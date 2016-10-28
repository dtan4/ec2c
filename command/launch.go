package command

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type LaunchCommand struct {
	Meta
}

func (c *LaunchCommand) Run(args []string) int {
	var (
		amiID                    string
		associatePublicIPAddress bool
		availabilityZone         string
		dryRun                   bool
		instanceCount            int64
		instanceType             string
		keyName                  string
		securityGroupIDs         string
		securityGroups           []*string
		subnetID                 string
		userData                 string
		volumeSize               int64
		volumeType               string
		arn                      string
		roleName                 string
	)

	var (
		arguments []string
	)

	svc := ec2.New(session.New(), &aws.Config{})

	flags := flag.NewFlagSet("dtan4", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.StringVar(&amiID, "ami", "", "AMI Id")
	flags.BoolVar(&associatePublicIPAddress, "publicip", false, "Associate Public Ip (default: false)")
	flags.StringVar(&availabilityZone, "az", "", "Availability zone")
	flags.BoolVar(&dryRun, "dry-run", false, "Dry run (default: false)")
	flags.Int64Var(&instanceCount, "count", 1, "Number of instances (default: 1)")
	flags.StringVar(&instanceType, "type", "", "Instance type")
	flags.StringVar(&keyName, "key", "", "SSH key name")
	flags.StringVar(&securityGroupIDs, "sg", "", "Security group Ids")
	flags.StringVar(&subnetID, "subnet", "", "Subnet Id")
	flags.StringVar(&userData, "userData", "", "User data")
	flags.Int64Var(&volumeSize, "volumeSize", 8, "Volume size (default: 8)")
	flags.StringVar(&volumeType, "volumeType", "gp2", "Volume type (default: gp2)")
	flags.StringVar(&arn, "arn", "", "ARN")
	flags.StringVar(&roleName, "roleName", "", "IAM Role Name")

	if err := flags.Parse(args[0:]); err != nil {
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	if arn != "" && roleName != "" {
		fmt.Fprintln(os.Stderr, "InvalidParameterCombination: The parameter 'roleName' and 'arn' can not be used in combination")
		return 1
	}

	for _, id := range strings.Split(securityGroupIDs, ",") {
		securityGroups = append(securityGroups, aws.String(id))
	}

	blockDeviceMappings := []*ec2.BlockDeviceMapping{
		&ec2.BlockDeviceMapping{
			DeviceName: aws.String("/dev/xvda"),
			Ebs: &ec2.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(true),
				VolumeSize:          aws.Int64(volumeSize),
				VolumeType:          aws.String(volumeType),
			},
		},
	}

	monitoring := &ec2.RunInstancesMonitoringEnabled{
		Enabled: aws.Bool(true),
	}

	opts := &ec2.RunInstancesInput{
		BlockDeviceMappings: blockDeviceMappings,
		DryRun:              aws.Bool(dryRun),
		ImageId:             aws.String(amiID),
		KeyName:             aws.String(keyName),
		InstanceType:        aws.String(instanceType),
		MaxCount:            aws.Int64(instanceCount),
		MinCount:            aws.Int64(1),
		Monitoring:          monitoring,
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Arn:  aws.String(arn),
			Name: aws.String(roleName),
		},
	}

	if subnetID != "" && associatePublicIPAddress {
		opts.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			&ec2.InstanceNetworkInterfaceSpecification{
				AssociatePublicIpAddress: aws.Bool(associatePublicIPAddress),
				DeviceIndex:              aws.Int64(int64(0)),
				SubnetId:                 aws.String(subnetID),
				Groups:                   securityGroups,
			},
		}
	}

	if userData != "" {
		buf, err := ioutil.ReadFile(userData)
		if err != nil {
			panic(err)
		}

		opts.UserData = aws.String(base64.StdEncoding.EncodeToString(buf))
	}

	resp, err := svc.RunInstances(opts)
	if err != nil {
		panic(err)
	}

	for _, instance := range resp.Instances {
		fmt.Println(*instance.InstanceId)
	}

	return 0
}

func (c *LaunchCommand) Synopsis() string {
	return "Launch new EC2 instance"
}

func (c *LaunchCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
