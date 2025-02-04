/*
 * Copyright 2021 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package aws

import (
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

const (
	mockVM1            = "mockVM1"
	mockVM1ID          = "mockVM1ID"
	blockVolumeId      = "blockVolumeID"
	networkInterfaceId = "networkInterfaceId"
	mockVMCreationTime = "2012-11-01T22:08:41+00:00"

	mockFunction1ID           = "arn:aws:lambda:eu-central-1:123456789:function:mock-function:1"
	mockFunction1             = "MockFunction1"
	mockFunction1Region       = "eu-central-1"
	mockFunction1CreationTime = "2012-11-01T22:08:41.0+00:00"
)

// mockEC2API implements the EC2API interface for mock testing
type mockEC2API struct {
}

// mockEC2APIWithErrors implements the EC2API interface (API call returning error) for mock testing
type mockEC2APIWithErrors struct {
}

// mockLambdaAPI implements the LambdaAPI interface for mock testing
type mockLambdaAPI struct {
}

// mockLambdaAPI implements the LambdaAPI interface for mock testing if >50 Lambda functions are discovered (not only 50)
type mockLambdaAPI51LambdaFunctions struct {
}

// mockLambdaAPIWithErrors implements the LambdaAPI interface (API call returning error) for mock testing
type mockLambdaAPIWithErrors struct {
}

// ListFunctions is the method implementation of the LambdaAPI interface
func (m mockLambdaAPI) ListFunctions(_ context.Context, _ *lambda.ListFunctionsInput, _ ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	return &lambda.ListFunctionsOutput{
		Functions: []lambdaTypes.FunctionConfiguration{
			{
				FunctionArn:  aws.String(mockFunction1ID),
				FunctionName: aws.String(mockFunction1),
				LastModified: aws.String(mockFunction1CreationTime),
			},
		},
		NextMarker:     nil,
		ResultMetadata: middleware.Metadata{},
	}, nil
}

func (m mockLambdaAPI51LambdaFunctions) ListFunctions(_ context.Context, input *lambda.ListFunctionsInput, _ ...func(*lambda.Options)) (output *lambda.ListFunctionsOutput, err error) {
	var lambdaFunctions []lambdaTypes.FunctionConfiguration
	nextMarker := "ShowNext"
	if input.Marker == nil {
		for i := 0; i < 50; i++ {
			lambdaFunctions = append(lambdaFunctions, lambdaTypes.FunctionConfiguration{
				// We have to set a time in a right format, otherwise the discoverer fails (parse error)
				LastModified: aws.String(mockFunction1CreationTime),
			})
		}
		output = &lambda.ListFunctionsOutput{
			Functions:  lambdaFunctions,
			NextMarker: aws.String(nextMarker),
		}
	} else if *input.Marker == nextMarker {
		for i := 0; i < 5; i++ {
			lambdaFunctions = append(lambdaFunctions, lambdaTypes.FunctionConfiguration{
				// We have to set a time in a right format, otherwise the discoverer fails (parse error)
				LastModified: aws.String(mockFunction1CreationTime),
			})
		}
		output = &lambda.ListFunctionsOutput{
			Functions:  lambdaFunctions,
			NextMarker: nil,
		}
	}
	return
}

// ListFunctions is the method implementation of the LambdaAPI interface
func (m mockLambdaAPIWithErrors) ListFunctions(_ context.Context, _ *lambda.ListFunctionsInput, _ ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	err := &smithy.GenericAPIError{
		Code:    "500",
		Message: "Internal Server Error",
	}
	return nil, err
}

// DescribeInstances is the method implementation of the EC2API interface
func (m mockEC2API) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	// block device mappings for output struct
	blockDeviceMappings := []types.InstanceBlockDeviceMapping{
		{
			DeviceName: aws.String("/dev/sdh"),
			Ebs: &types.EbsInstanceBlockDevice{
				AttachTime:          nil,
				DeleteOnTermination: nil,
				Status:              "",
				VolumeId:            aws.String(blockVolumeId),
			},
		},
	}
	// tags for output struct
	tags := []types.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(mockVM1),
		},
	}
	networkInterfaces := []types.InstanceNetworkInterface{
		{
			NetworkInterfaceId: aws.String(networkInterfaceId),
		},
	}
	// launch time
	launchTime, err := time.Parse(time.RFC3339, mockVMCreationTime)
	if err != nil {
		log.Error(err)
	}

	// output struct containing all necessary information
	output := &ec2.DescribeInstancesOutput{
		NextToken: nil,
		Reservations: []types.Reservation{{
			Groups: nil,
			Instances: []types.Instance{{
				BlockDeviceMappings: blockDeviceMappings,
				InstanceId:          aws.String(mockVM1ID),
				NetworkInterfaces:   networkInterfaces,
				Tags:                tags,
				LaunchTime:          aws.Time(launchTime),
			}},
			OwnerId:       nil,
			RequesterId:   nil,
			ReservationId: nil,
		}},
		ResultMetadata: middleware.Metadata{},
	}
	return output, nil
}

// DescribeInstances is the method implementation of the EC2API interface
func (m mockEC2APIWithErrors) DescribeInstances(_ context.Context, _ *ec2.DescribeInstancesInput, _ ...func(options *ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	err := &smithy.GenericAPIError{
		Code:    "ConnectionError",
		Message: "Couldn't resolve host. Bad connection?",
	}
	return nil, err
}

func TestComputeDiscovery_List(t *testing.T) {
	d := computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		functionAPI:       mockLambdaAPI{},
		isDiscovering:     true,
		awsConfig: &Client{
			cfg: aws.Config{
				Region: "eu-central-1",
			},
			accountID: aws.String("MockAccountID1234"),
		},
	}
	list, err := d.List()
	assert.Nil(t, err)
	assert.NotEmpty(t, list)

	d = computeDiscovery{
		virtualMachineAPI: mockEC2APIWithErrors{},
	}
	_, err = d.List()
	assert.NotNil(t, err)

	d = computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		functionAPI:       mockLambdaAPIWithErrors{},
		isDiscovering:     true,
		awsConfig: &Client{
			cfg: aws.Config{
				Region: "eu-central-1",
			},
			accountID: aws.String("MockAccountID1234"),
		},
	}
	_, err = d.List()
	assert.NotNil(t, err)
}

func TestComputeDiscovery_discoverVirtualMachines(t *testing.T) {
	d := computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		isDiscovering:     true,
		awsConfig: &Client{
			cfg: aws.Config{
				Region: "eu-central-1",
			},
			accountID: aws.String("MockAccountID1234"),
		},
	}
	machines, err := d.discoverVirtualMachines()
	assert.Nil(t, err)
	testMachine := machines[0]
	assert.Equal(t, mockVM1, testMachine.Name)
	assert.Equal(t, voc.ResourceID("arn:aws:ec2:eu-central-1:MockAccountID1234:instance/mockVM1ID"), testMachine.ID)
	assert.NotEmpty(t, testMachine.BlockStorage)
	assert.False(t, testMachine.BootLog.Enabled)
	assert.False(t, testMachine.OSLog.Enabled)
	assert.Equal(t, int64(0), testMachine.CreationTime)
	assert.Equal(t, mockFunction1Region, testMachine.GeoLocation.Region)

	d = computeDiscovery{
		virtualMachineAPI: mockEC2APIWithErrors{},
	}
	_, err = d.discoverVirtualMachines()
	assert.NotNil(t, err)

}

func TestComputeDiscover_Name(t *testing.T) {
	d := computeDiscovery{
		virtualMachineAPI: mockEC2API{},
		isDiscovering:     true,
		awsConfig:         &Client{},
	}
	assert.Equal(t, "AWS Compute", d.Name())
}

func TestComputeDiscovery_getNameOfVM(t *testing.T) {
	type fields struct {
		api           EC2API
		isDiscovering bool
		awsConfig     *Client
	}
	type args struct {
		vm types.Instance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"First Test without tag",
			fields{},
			args{vm: types.Instance{InstanceId: aws.String(mockVM1ID)}},
			mockVM1ID,
		},
		{
			"Second test with tag",
			fields{},
			args{vm: types.Instance{InstanceId: aws.String(mockVM1ID),
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(mockVM1),
					},
				}}},
			mockVM1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &computeDiscovery{
				virtualMachineAPI: tt.fields.api,
				isDiscovering:     tt.fields.isDiscovering,
				awsConfig:         tt.fields.awsConfig,
			}
			if got := d.getNameOfVM(&tt.args.vm); got != tt.want {
				t.Errorf("getNameOfVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeDiscovery_discoverFunctions(t *testing.T) {
	type fields struct {
		virtualMachineAPI EC2API
		functionAPI       LambdaAPI
		isDiscovering     bool
		awsConfig         *Client
	}
	mockClient := &Client{
		cfg: aws.Config{
			Region: "eu-central-1",
		},
	}
	// creationTime, _ := time.Parse(time.RFC3339, mockFunction1CreationTime)
	tests := []struct {
		name    string
		fields  fields
		want    []voc.Function
		wantErr bool
	}{
		// Test cases
		{
			"Test case 1 (no error)",
			fields{
				functionAPI: mockLambdaAPI{},
				awsConfig:   mockClient,
			},
			//args: args{client: mockClient},
			[]voc.Function{
				{Compute: &voc.Compute{
					CloudResource: &voc.CloudResource{
						ID:           mockFunction1ID,
						Name:         mockFunction1,
						CreationTime: int64(0),
						Type:         []string{"Function", "Compute", "Resource"},
						GeoLocation: voc.GeoLocation{
							Region: mockFunction1Region,
						},
					},
				}},
			},
			false,
		},
		{
			"Test case 3 (API error)",
			fields{
				functionAPI: mockLambdaAPIWithErrors{},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &computeDiscovery{
				virtualMachineAPI: tt.fields.virtualMachineAPI,
				functionAPI:       tt.fields.functionAPI,
				isDiscovering:     tt.fields.isDiscovering,
				awsConfig:         tt.fields.awsConfig,
			}
			got, err := d.discoverFunctions()
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverFunctions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discoverFunctions() got = %v, want %v", got, tt.want)
			}
		})
	}

	// Testing the case where two API Calls have to be made due to limit of returned functions
	d := computeDiscovery{
		functionAPI: mockLambdaAPI51LambdaFunctions{},
		awsConfig:   mockClient,
	}
	functions, err := d.discoverFunctions()
	assert.Nil(t, err)
	assert.Less(t, 50, len(functions))

}

func TestComputeDiscovery_NewComputeDiscovery(t *testing.T) {
	// Mock newFromConfigs and store the original functions back at the end of the test
	oldEC2 := newFromConfigEC2
	defer func() { newFromConfigEC2 = oldEC2 }()
	oldLambda := newFromConfigLambda
	defer func() { newFromConfigLambda = oldLambda }()

	newFromConfigEC2 = func(cfg aws.Config, optFns ...func(*ec2.Options)) *ec2.Client {
		return &ec2.Client{}
	}
	newFromConfigLambda = func(cfg aws.Config, optFns ...func(*lambda.Options)) *lambda.Client {
		return &lambda.Client{}
	}

	type args struct {
		client *Client
	}
	mockClient := &Client{
		cfg: aws.Config{
			Region: "eu-central-1",
		},
		accountID: aws.String("1234"),
	}
	tests := []struct {
		name string
		args args
		want discovery.Discoverer
	}{
		{
			args: args{client: mockClient},
			want: &computeDiscovery{
				virtualMachineAPI: &ec2.Client{},
				functionAPI:       &lambda.Client{},
				isDiscovering:     true,
				awsConfig:         mockClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewComputeDiscovery(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewComputeDiscovery() = %v, want %v", got, tt.want)
			}
		})
	}
}
