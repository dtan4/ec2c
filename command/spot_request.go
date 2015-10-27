package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SpotRequestCommand struct {
	Meta
}

func (c *SpotRequestCommand) Run(args []string) int {
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
		spotPrice                string
		subnetId                 string
		volumeSize               int64
		volumeType               string
	)

	var (
		arguments []string
	)

	svc := ec2.New(&aws.Config{})

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
	flags.StringVar(&spotPrice, "spotPrice", "", "Maximum bidding price")
	flags.StringVar(&subnetId, "subnet", "", "Subnet Id")
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

	launchSpecification := &ec2.RequestSpotLaunchSpecification{
		BlockDeviceMappings: blockDeviceMappings,
		ImageId:             aws.String(amiId),
		KeyName:             aws.String(keyName),
		InstanceType:        aws.String(instanceType),
		Monitoring:          monitoring,
		SecurityGroups:      securityGroups,
		SubnetId:            aws.String(subnetId),
	}

	if subnetId != "" && associatePublicIpAddress {
		launchSpecification.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
			&ec2.InstanceNetworkInterfaceSpecification{
				AssociatePublicIpAddress: aws.Bool(associatePublicIpAddress),
				DeviceIndex:              aws.Int64(int64(0)),
				SubnetId:                 aws.String(subnetId),
				Groups:                   securityGroups,
			},
		}
	}

	opts := &ec2.RequestSpotInstancesInput{
		DryRun:              aws.Bool(dryRun),
		InstanceCount:       aws.Int64(instanceCount),
		LaunchSpecification: launchSpecification,
		SpotPrice:           aws.String(spotPrice),
	}

	resp, err := svc.RequestSpotInstances(opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)

	return 0
}

func (c *SpotRequestCommand) Synopsis() string {
	return "Request new Spot Instances"
}

func (c *SpotRequestCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
