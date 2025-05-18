package gcpclient

import (
	"context"
	"encoding/json"
	"os"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type PubSubClient struct {
	Client *pubsub.Client
	Ctx    context.Context
}

func NewPubSubClient(projectId string) (*PubSubClient, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return nil, err
	}

	return &PubSubClient{
		Client: client,
		Ctx:    ctx,
	}, nil

}

func NewPubSubClientWithAuth(credentialsPath string) (*PubSubClient, error) {
	configData, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}

	var config PubSubConfig
	err = json.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.ProjectId, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, err
	}

	return &PubSubClient{
		Client: client,
		Ctx:    ctx,
	}, nil
}
