package ebs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
)

type EBSProtectedEntityTypeManager struct {
	session                                          session.Session
}

func NewEBSProtectedEntityTypeManagerFromConfig(params map[string]interface{}, s3URLBase string,
	logger logrus.FieldLogger) (*EBSProtectedEntityTypeManager, error) {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(params["aws-region"].(string)),
	}))


	return &EBSProtectedEntityTypeManager{
		session: *session,
	}, nil
}

func (this EBSProtectedEntityTypeManager) GetTypeName() string {
	panic("implement me")
}

func (this EBSProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id astrolabe.ProtectedEntityID) (astrolabe.ProtectedEntity, error) {
	panic("implement me")
}

func (this EBSProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]astrolabe.ProtectedEntityID, error) {
	svc := ec2.New(&this.session)
	input := &ec2.DescribeVolumesInput{}

	result, err := svc.DescribeVolumes(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	var returnIDs []astrolabe.ProtectedEntityID
	for _, curEBS := range result.Volumes {
		returnIDs = append(returnIDs, astrolabe.NewProtectedEntityID("ebs", *curEBS.VolumeId))
	}
	return returnIDs, nil
}

func (this EBSProtectedEntityTypeManager) Copy(ctx context.Context, pe astrolabe.ProtectedEntity, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	panic("implement me")
}

func (this EBSProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, info astrolabe.ProtectedEntityInfo, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	panic("implement me")
}
