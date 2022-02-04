package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Client struct {
}

func (c *Client) Create(ctx context.Context, username, email, password string) (id string, err error) {
	userUUID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("new UUID generation failed: %w", err)
	}

	zap.S().Infow("create user", "user_uuid", userUUID, "username", username, "email", email)

	return userUUID.String(), nil
}
