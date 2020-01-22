package ebs

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
	panic("implement me")
}

func (this EBSProtectedEntityTypeManager) Copy(ctx context.Context, pe astrolabe.ProtectedEntity, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	panic("implement me")
}

func (this EBSProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, info astrolabe.ProtectedEntityInfo, options astrolabe.CopyCreateOptions) (astrolabe.ProtectedEntity, error) {
	panic("implement me")
}
