package sdk

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

//go:generate go run github.com/golang/mock/mockgen -destination=../../../mocks/mocksdk/mock_assume_clients.go -package=mocksdk -source ./assume_clients.go AssumeClients

type Client = any

// AssumeClients provides the clients by a role.
type AssumeClients[C Client] interface {
	// Get returns the Lambda for the provided role.
	// If the role is nil, will return the default Lambda without assuming any identity.
	Get(assumeRole *arn.ARN) C
}

// ClientProps for creating a new Client.
type ClientProps[C Client] struct {
	// Config to be used.
	Config aws.Config
	// NewClient factory func.
	NewClient func(aws.Config, aws.CredentialsProvider) C
}

// NewAssumeClients returns a AssumeClients for a default aws.Config.
func NewAssumeClients[C Client](props ClientProps[C]) AssumeClients[C] {
	return &assumeClients[C]{props: props, backing: make(map[string]C)}
}

type assumeClients[C Client] struct {
	props     ClientProps[C]
	stsClient *sts.Client
	backing   map[string]C
	mutex     sync.Mutex
}

func (l *assumeClients[C]) Get(assumeRole *arn.ARN) C {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	key := l.keyOf(assumeRole)
	if existing, ok := l.backing[key]; ok {
		return existing
	}
	created := l.create(assumeRole)
	l.backing[key] = created
	return created
}

func (l *assumeClients[C]) keyOf(role *arn.ARN) string {
	if role == nil {
		return ""
	}
	return role.String()
}

func (l *assumeClients[C]) create(role *arn.ARN) C {
	// In case no role was provided, use the default config
	if role == nil {
		return l.props.NewClient(l.props.Config, l.props.Config.Credentials)
	}
	// Create a sts sdk on demand
	if l.stsClient == nil {
		l.stsClient = sts.NewFromConfig(l.props.Config)
	}
	// Create a new client assuming the role and caching the credentials
	return l.props.NewClient(l.props.Config, l.props.Config.Credentials)
}
