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

type RequestCommand struct {
	Meta
}

func (c *RequestCommand) Run(args []string) int {
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
		spotPrice                string
		subnetID                 string
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

	flags.StringVar(&amiID, "ami", "", "AMI Id")
	flags.BoolVar(&associatePublicIPAddress, "publicip", false, "Associate Public Ip (default: false)")
	flags.StringVar(&availabilityZone, "az", "", "Availability zone")
	flags.BoolVar(&dryRun, "dry-run", false, "Dry run (default: false)")
	flags.Int64Var(&instanceCount, "count", 1, "Number of instances (default: 1)")
	flags.StringVar(&instanceType, "type", "", "Instance type")
	flags.StringVar(&keyName, "key", "", "SSH key name")
	flags.StringVar(&securityGroupIDs, "sg", "", "Security group Ids")
	flags.StringVar(&spotPrice, "spotPrice", "", "Maximum bidding price")
	flags.StringVar(&subnetID, "subnet", "", "Subnet Id")
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

	launchSpecification := &ec2.RequestSpotLaunchSpecification{
		BlockDeviceMappings: blockDeviceMappings,
		ImageId:             aws.String(amiID),
		KeyName:             aws.String(keyName),
		InstanceType:        aws.String(instanceType),
		Monitoring:          monitoring,
		SecurityGroups:      securityGroups,
		SubnetId:            aws.String(subnetID),
	}

	if subnetID != "" && associatePublicIPAddress {
		launchSpecification.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
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

		launchSpecification.UserData = aws.String(base64.StdEncoding.EncodeToString(buf))
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

func (c *RequestCommand) Synopsis() string {
	return "Request new Spot Instances"
}

func (c *RequestCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
