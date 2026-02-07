// Package commontransform supplies a commonTransformClient
// which does transformation on fields that are shared among all event types
package commontransform

import (
	"context"
	"time"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

type commonTransformClient struct{}

// NewCommonTransformClient creates a new commonTransformClient
func NewCommonTransformClient() transforms.TransformClient {
	return &commonTransformClient{}
}

// TransformEvent performs common event transformation on the Event
func (t *commonTransformClient) TransformEvent(_ context.Context, partialUpdate, fullModel *model.Event) (
	*model.Event, error,
) {
	var outDelta *model.Event

	marketCloseAtTime := time.Now()

	marketClosed := make(map[string]bool)

	for _, market := range partialUpdate.Markets {
		if market.GetBettingStatus().GetValue() == model.BettingStatus_BettingClosed {
			marketClosed[market.GetID()] = true
		}
	}

	if len(marketClosed) == 0 {
		return outDelta, nil // no market closed, nothing to update
	}

	var hasChange bool

	for _, market := range fullModel.Markets {
		if marketClosed[market.GetID()] && !market.HasClosedAt() {
			hasChange = true
			market.SetClosedAt(&model.OptionalInt64{
				Value: marketCloseAtTime.UnixNano(),
			})
		}
	}

	if !hasChange {
		return outDelta, nil // no need to update if no new market is closed
	}

	outDelta = &model.Event{
		ID:      partialUpdate.ID,
		Markets: fullModel.Markets,
	}

	return outDelta, nil
}

func (t *commonTransformClient) GetName() string {
	return "CommonTransform"
}
