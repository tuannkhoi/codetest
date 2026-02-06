// Package merger merges instance of the codetest models
package merger

import (
	"context"

	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

// ServiceClient is an interface representing a service that can merge models together.
type ServiceClient interface {
	MergeEvent(ctx context.Context, left, right *model.Event) (*model.Event, error)
}

type inlineMergerClient struct {
}

// NewInlineMergerClient  creates a new instance of inlineMergerClient.
func NewInlineMergerClient() ServiceClient {
	return &inlineMergerClient{}
}

func (c *inlineMergerClient) MergeEvent(ctx context.Context, left, right *model.Event) (*model.Event, error) {
	update := MergeEvent(ctx, left, right)
	return update, nil
}
