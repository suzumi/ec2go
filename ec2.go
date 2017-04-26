package ec2go

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func newSession() (*session.Session, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func getInstances(sess *session.Session, isAll bool) [][]string {

	client := ec2.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})

	var resp *ec2.DescribeInstancesOutput
	var respErr error
	if isAll {
		resp, respErr = client.DescribeInstances(nil)
	} else {
		states := []*string{}
		states = append(states, toPtr("running"))

		filters := []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("instance-state-name"),
				Values: states,
			},
		}

		req := ec2.DescribeInstancesInput{Filters: filters}
		resp, respErr = client.DescribeInstances(&req)
	}

	if respErr != nil {
		panic(respErr)
	}

	var instances [][]string
	rows := []map[string]string{}
	for idx := range resp.Reservations {

		row := map[string]string{}
		for _, inst := range resp.Reservations[idx].Instances {
			for _, v := range inst.Tags {
				if *v.Key == "Name" {
					row["Name"] = *v.Value
				}
			}
			row["InstanceId"] = *inst.InstanceId
			row["Type"] = *inst.InstanceType
			if inst.PrivateIpAddress != nil {
				row["PrivateIP"] = *inst.PrivateIpAddress
			}
			row["State"] = *inst.State.Name
		}
		rows = append(rows, row)
	}

	for _, v := range rows {
		instance := []string{
			v["InstanceId"],
			v["Name"],
			v["Type"],
			v["PrivateIP"],
			v["State"],
		}
		instances = append(instances, instance)
	}
	return instances
}
