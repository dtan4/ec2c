package command

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
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
		amiId                    string
		associatePublicIpAddress bool
		availabilityZone         string
		dryRun                   bool
		instanceCount            int64
		instanceType             string
		keyName                  string
		securityGroupIds         string
		securityGroups           []*string
		subnetId                 string
		userData                 string
		volumeSize               int64
		volumeType               string
	)

	var (
		arguments []string
	)

	svc := ec2.New(session.New(), &aws.Config{})

	flags := flag.NewFlagSet("dtan4", flag.ContinueOnError)
	flags.Usage = func() {}

	flags.StringVar(&amiId, "ami", "", "AMI Id")
	flags.BoolVar(&associatePublicIpAddress, "publicip", false, "Associate Public Ip (default: false)")
	flags.StringVar(&availabilityZone, "az", "", "Availability zone")
	flags.BoolVar(&dryRun, "dry-run", false, "Dry run (default: false)")
	flags.Int64Var(&instanceCount, "count", 1, "Number of instances (default: 1)")
	flags.StringVar(&instanceType, "type", "", "Instance type")
	flags.StringVar(&keyName, "key", "", "SSH key name")
	flags.StringVar(&securityGroupIds, "sg", "", "Security group Ids")
	flags.StringVar(&subnetId, "subnet", "", "Subnet Id")
	flags.StringVar(&userData, "userData", "", "User data")
	flags.Int64Var(&volumeSize, "volumeSize", 8, "Volume size (default: 8)")
	flags.StringVar(&volumeType, "volumeType", "gp2", "Volume type (default: gp2)")

	if err := flags.Parse(args[0:]); err != nil {
		return 1
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	for _, id := range strings.Split(securityGroupIds, ",") {
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
		ImageId:             aws.String(amiId),
		KeyName:             aws.String(keyName),
		InstanceType:        aws.String(instanceType),
		MaxCount:            aws.Int64(instanceCount),
		MinCount:            aws.Int64(1),
		Monitoring:          monitoring,
	}

	if subnetId != "" && associatePublicIpAddress {
		opts.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			&ec2.InstanceNetworkInterfaceSpecification{
				AssociatePublicIpAddress: aws.Bool(associatePublicIpAddress),
				DeviceIndex:              aws.Int64(int64(0)),
				SubnetId:                 aws.String(subnetId),
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
