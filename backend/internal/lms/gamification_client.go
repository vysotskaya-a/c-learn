package lms

import (
	"context"
	"time"

	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/httpclient"
)

type GamificationClient struct {
	client *httpclient.Client
}

func NewGamificationClient(baseURL string) *GamificationClient {
	return &GamificationClient{
		client: httpclient.New(baseURL, 5*time.Second),
	}
}

func (gc *GamificationClient) AwardXP(ctx context.Context, req models.XPAwardRequest) (*models.XPAwardResponse, error) {
	var result models.XPAwardResponse
	if err := gc.client.Post(ctx, "/internal/xp/award", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
