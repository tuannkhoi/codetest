package racetransform_test

import (
	"context"
	"testing"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms/racetransform"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func TestTransformEvent_SkipsWhenRaceCourseNotUpdated(t *testing.T) {
	client := racetransform.NewRaceTransformClient()

	partial := &model.Event{ID: "race-1"}
	full := &model.Event{
		ID:       "race-1",
		RaceData: &model.RaceEvent{},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output when RaceCourse not updated, got %#v", out)
	}
}

func TestTransformEvent_SkipsWhenStateAlreadySet(t *testing.T) {
	client := racetransform.NewRaceTransformClient()

	partial := &model.Event{
		ID: "race-2",
		RaceData: &model.RaceEvent{
			RaceCourse: &model.OptionalString{Value: "flemington"},
		},
	}
	full := &model.Event{
		ID: "race-2",
		RaceData: &model.RaceEvent{
			RaceCourse: &model.OptionalString{Value: "flemington"},
			State:      &model.OptionalString{Value: "VIC"},
		},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output when state already set, got %#v", out)
	}
}

func TestTransformEvent_SkipsUnknownRaceCourse(t *testing.T) {
	client := racetransform.NewRaceTransformClient()

	partial := &model.Event{
		ID: "race-3",
		RaceData: &model.RaceEvent{
			RaceCourse: &model.OptionalString{Value: "unknown_course"},
		},
	}
	full := &model.Event{
		ID: "race-3",
		RaceData: &model.RaceEvent{
			RaceCourse: &model.OptionalString{Value: "unknown_course"},
		},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output for unknown racecourse, got %#v", out)
	}
}

func TestTransformEvent_SetsState(t *testing.T) {
	client := racetransform.NewRaceTransformClient()

	partial := &model.Event{
		ID: "race-4",
		RaceData: &model.RaceEvent{
			RaceCourse: &model.OptionalString{Value: "flemington"},
		},
	}
	full := &model.Event{
		ID: "race-4",
		RaceData: &model.RaceEvent{
			RaceCourse: &model.OptionalString{Value: "flemington"},
		},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatalf("expected output event, got nil")
	}
	if out.ID != "race-4" {
		t.Fatalf("expected ID %q, got %q", "race-4", out.ID)
	}
	if out.RaceData == nil || out.RaceData.State == nil {
		t.Fatalf("expected state to be set, got %#v", out.RaceData)
	}
	if out.RaceData.State.Value != "VIC" {
		t.Fatalf("expected state %q, got %q", "VIC", out.RaceData.State.Value)
	}
}

func TestGetName(t *testing.T) {
	client := racetransform.NewRaceTransformClient()
	if got := client.GetName(); got != "RaceTransform" {
		t.Fatalf("expected name %q, got %q", "RaceTransform", got)
	}
}
