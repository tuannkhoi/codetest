// Package racetransform supplies a raceTransformClient
package racetransform

import (
	"context"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

type raceTransformClient struct{}

// NewRaceTransformClient creates a new raceTransformClient
func NewRaceTransformClient() transforms.TransformClient {
	return &raceTransformClient{}
}

var stateByRaceCourse = map[string]string{
	"flemington":    "VIC",
	"caulfield":     "VIC",
	"moonee_valley": "VIC",

	"royal_randwick":   "NSW",
	"rosehill_gardens": "NSW",
	"canterbury_park":  "NSW",

	"eagle_farm": "QLD",
	"doomben":    "QLD",
	"gold_coast": "QLD",

	"morphettville": "SA",
	"cheltenham":    "SA",
	"victoria_park": "SA",

	"ascot":        "WA",
	"belmont_park": "WA",
	"bunbury":      "WA",

	"hobart":     "TAS",
	"devonport":  "TAS",
	"launceston": "TAS",

	"canberra": "ACT",

	"darwin":        "NT",
	"alice_springs": "NT",
}

func (t *raceTransformClient) TransformEvent(_ context.Context, partialUpdate, fullModel *model.Event) (
	*model.Event, error,
) {
	var outDelta *model.Event

	if !partialUpdate.GetRaceData().HasRaceCourse() {
		return outDelta, nil // if the RaceCourse didn't update on this update skip processing the event
	}

	if fullModel.GetRaceData().HasState() {
		return outDelta, nil // no need to update the state if the state already set
	}

	stateName, ok := stateByRaceCourse[fullModel.GetRaceData().GetRaceCourse().GetValue()]
	if !ok {
		return outDelta, nil // unknown racecourse, can't determine the state so don't change anything
	}

	outDelta = &model.Event{
		ID: partialUpdate.ID,
		RaceData: &model.RaceEvent{
			State: &model.OptionalString{
				Value: stateName,
			},
		},
	}

	return outDelta, nil
}

func (t *raceTransformClient) GetName() string {
	return "RaceTransform"
}
