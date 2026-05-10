package lms

import (
	"context"
	"time"

	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/httpclient"
)

type RunnerClient struct {
	client *httpclient.Client
}

func NewRunnerClient(baseURL string) *RunnerClient {
	return &RunnerClient{
		client: httpclient.New(baseURL, 30*time.Second),
	}
}

func (rc *RunnerClient) Run(ctx context.Context, req models.RunRequest) (*models.RunResult, error) {
	var result models.RunResult
	if err := rc.client.Post(ctx, "/internal/run", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
