// Package sporttransform supplies a sportTransformClient
// which does transformation on fields that are specific to sport events
package sporttransform

import (
	"context"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

type sportTransformClient struct{}

// NewSportTransformClient creates a new sportTransformClient
func NewSportTransformClient() transforms.TransformClient {
	return &sportTransformClient{}
}

var sportTypeMap = map[string]string{
	"soccer":       "Soccer",
	"rugby_league": "Rugby League",
}

// TransformEvent performs sport specific transformation on the Event
func (t *sportTransformClient) TransformEvent(_ context.Context, partialUpdate, fullModel *model.Event) (
	*model.Event, error,
) {
	var outDelta *model.Event
	if partialUpdate.EventTypeID == nil {
		return outDelta, nil // if the EventTypeID didn't update on this update skip processing the event
	}

	if fullModel.GetSportData().GetName() != nil {
		return outDelta, nil // no need to update the name if the name already set
	}

	sportName, ok := sportTypeMap[fullModel.GetEventTypeID().GetValue()]
	if !ok {
		return outDelta, nil // unknown sport type so don't change anything
	}

	outDelta = &model.Event{
		ID:        partialUpdate.ID,
		SportData: &model.SportEvent{Name: &model.OptionalString{Value: sportName}},
	}

	return outDelta, nil
}

func (t *sportTransformClient) GetName() string {
	return "SportsTransform"
}
