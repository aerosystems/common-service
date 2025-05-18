package gcpclient

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func NewFirebaseClient(cfg *FirebaseConfig) (*auth.Client, error) {
	var opts []option.ClientOption
	if file := cfg.CredentialsPath; file != "" {
		opts = append(opts, option.WithCredentialsFile(file))
	}

	config := &firebase.Config{ProjectID: cfg.ProjectId}
	firebaseApp, err := firebase.NewApp(context.Background(), config, opts...)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %w\n", err)
	}

	authClient, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return nil, errors.New("unable to create firebase Auth client")
	}
	return authClient, nil
}
