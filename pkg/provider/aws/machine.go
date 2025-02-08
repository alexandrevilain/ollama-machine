// Licensed to Alexandre VILAIN under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Alexandre VILAIN licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alexandrevilain/ollama-machine/pkg/ollama"
	"github.com/alexandrevilain/ollama-machine/pkg/provider"
	"github.com/alexandrevilain/ollama-machine/pkg/ssh"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

var waitInstanceTerminated = 10 * time.Minute

type MachineManager struct {
	client *ec2.Client
}

func newMachineManager(client *ec2.Client) *MachineManager {
	return &MachineManager{client: client}
}

func (m *MachineManager) MachineKind() provider.MachineKind {
	return provider.MachineKindVM
}

func (m *MachineManager) createSecurityGroup(ctx context.Context, _ *provider.CreateMachineRequest) (string, error) {
	securityGroupName := fmt.Sprintf("ollama-machine-%s", uuid.New().String()) // Create a unique security group name per machine.

	createSgInput := &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(securityGroupName),
		Description: aws.String("Security group for SSH and Ollama access"),
	}

	sgResult, err := m.client.CreateSecurityGroup(ctx, createSgInput)
	if err != nil {
		return "", fmt.Errorf("unable to create security group: %w", err)
	}

	ingressRules := []types.IpPermission{
		{
			IpProtocol: aws.String("tcp"),
			FromPort:   aws.Int32(ssh.DefaultPort),
			ToPort:     aws.Int32(ssh.DefaultPort),
			IpRanges: []types.IpRange{
				{
					CidrIp:      aws.String("0.0.0.0/0"),
					Description: aws.String("Allow SSH access from anywhere"),
				},
			},
		},
		// This permission is too permissive and should be created only if the user asks for a public instance.
		// But we don't have this information in the request for now.
		{
			IpProtocol: aws.String("tcp"),
			FromPort:   aws.Int32(ollama.DefaultPort),
			ToPort:     aws.Int32(ollama.DefaultPort),
			IpRanges: []types.IpRange{
				{
					CidrIp:      aws.String("0.0.0.0/0"),
					Description: aws.String("Allow Ollama access"),
				},
			},
		},
	}

	authorizeInput := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       sgResult.GroupId,
		IpPermissions: ingressRules,
	}

	_, err = m.client.AuthorizeSecurityGroupIngress(ctx, authorizeInput)
	if err != nil {
		return "", fmt.Errorf("unable to set security group ingress rules: %w", err)
	}

	return *sgResult.GroupId, nil
}

func (m *MachineManager) Create(ctx context.Context, req *provider.CreateMachineRequest) (*provider.Machine, error) { //nolint:funlen
	securityGroupID, err := m.createSecurityGroup(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create security group: %w", err)
	}

	if req.InstanceType == "" {
		req.InstanceType = "t3.micro" // If instance type has been missed, we use a cheap default one.
	}

	var imageID string
	if strings.HasPrefix(req.Image, "ami-") {
		imageID = req.Image
	} else {
		images, err := m.client.DescribeImages(ctx, &ec2.DescribeImagesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("name"),
					Values: []string{req.Image},
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe images: %w", err)
		}

		if len(images.Images) == 0 {
			return nil, errors.New("image not found")
		}

		imageID = *images.Images[0].ImageId
	}

	// Define instance parameters
	input := &ec2.RunInstancesInput{
		ImageId:      aws.String(imageID),
		InstanceType: types.InstanceType(req.InstanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		UserData:     aws.String(base64.StdEncoding.EncodeToString(req.UserData)),
		// Enable public IP address.
		NetworkInterfaces: []types.InstanceNetworkInterfaceSpecification{
			{
				DeviceIndex:              aws.Int32(0),
				AssociatePublicIpAddress: aws.Bool(true),
				DeleteOnTermination:      aws.Bool(true),
				Groups:                   []string{securityGroupID},
			},
		},
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(req.Name),
					},
					{
						Key:   aws.String("created_by"),
						Value: aws.String("ollama-machine"),
					},
				},
			},
		},
	}

	if req.Zone != "" {
		input.Placement = &types.Placement{
			AvailabilityZone: aws.String(req.Zone),
		}
	}

	result, err := m.client.RunInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	machine := m.machineFromInstance(result.Instances[0])
	machine.Name = req.Name

	return m.machineFromInstance(result.Instances[0]), nil
}

func (m *MachineManager) Delete(ctx context.Context, id string) error {
	instance, err := m.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		return fmt.Errorf("failed to describe instance: %w", err)
	}

	if len(instance.Reservations) == 0 || len(instance.Reservations[0].Instances) == 0 {
		return errors.New("instance not found")
	}

	securityGroups := instance.Reservations[0].Instances[0].SecurityGroups

	_, err = m.client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		return fmt.Errorf("failed to terminate instance: %w", err)
	}

	waiter := ec2.NewInstanceTerminatedWaiter(m.client)
	waitInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	}

	if err := waiter.Wait(ctx, waitInput, waitInstanceTerminated); err != nil {
		return fmt.Errorf("failed to wait for instance to be terminated: %w", err)
	}

	for _, sg := range securityGroups {
		_, err = m.client.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{
			GroupId: sg.GroupId,
		})
		if err != nil {
			return fmt.Errorf("failed to delete security group: %w", err)
		}
	}

	return nil
}

func (m *MachineManager) Start(ctx context.Context, id string) error {
	_, err := m.client.StartInstances(ctx, &ec2.StartInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		return fmt.Errorf("failed to start instance: %w", err)
	}

	return nil
}

func (m *MachineManager) Stop(ctx context.Context, id string) error {
	_, err := m.client.StopInstances(ctx, &ec2.StopInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		return fmt.Errorf("failed to stop instance: %w", err)
	}

	return nil
}

func (m *MachineManager) Get(ctx context.Context, id string) (*provider.Machine, error) {
	instance, err := m.client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe instance: %w", err)
	}

	if len(instance.Reservations) == 0 || len(instance.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance not found")
	}

	return m.machineFromInstance(instance.Reservations[0].Instances[0]), nil
}

func (m *MachineManager) machineFromInstance(instance types.Instance) *provider.Machine {
	var ip string
	if instance.PublicIpAddress != nil {
		ip = *instance.PublicIpAddress
	}

	state := provider.MachineStatePending
	if instance.State != nil {
		switch instance.State.Name {
		case types.InstanceStateNameRunning:
			state = provider.MachineStateRunning
		case types.InstanceStateNamePending, types.InstanceStateNameShuttingDown, types.InstanceStateNameStopping:
			state = provider.MachineStatePending
		case types.InstanceStateNameTerminated:
			state = provider.MachineStateTerminated
		case types.InstanceStateNameStopped:
			state = provider.MachineStateStopped
		}
	}

	// Get instance name from tags
	var instanceName string
	for _, tag := range instance.Tags {
		if *tag.Key == "Name" {
			instanceName = *tag.Value

			break
		}
	}

	// Fallback to instance ID if name is not set.
	if instanceName == "" {
		instanceName = *instance.InstanceId
	}

	return &provider.Machine{
		ID:     *instance.InstanceId,
		Name:   instanceName,
		IP:     ip,
		Region: m.client.Options().Region,
		State:  state,
	}
}
