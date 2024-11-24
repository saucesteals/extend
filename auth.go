package extend

import (
	"context"
	"time"
)

type Authenticator interface {
	GetAccessToken(ctx context.Context) (string, error)
	Expiry() time.Time
	Refresh(ctx context.Context) (string, error)
}
