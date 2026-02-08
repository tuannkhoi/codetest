package sporttransform_test

import (
	"context"
	"testing"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms/sporttransform"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func TestTransformEvent_SkipsWhenEventTypeNotUpdated(t *testing.T) {
	client := sporttransform.NewSportTransformClient()

	partial := &model.Event{ID: "evt-1"}
	full := &model.Event{
		ID:        "evt-1",
		SportData: &model.SportEvent{},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output when EventTypeID not updated, got %#v", out)
	}
}

func TestTransformEvent_SkipsWhenSportNameAlreadySet(t *testing.T) {
	client := sporttransform.NewSportTransformClient()

	partial := &model.Event{
		ID:          "evt-2",
		EventTypeID: &model.OptionalString{Value: "soccer"},
	}
	full := &model.Event{
		ID:          "evt-2",
		EventTypeID: &model.OptionalString{Value: "soccer"},
		SportData:   &model.SportEvent{Name: &model.OptionalString{Value: "Soccer"}},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output when sport name already set, got %#v", out)
	}
}

func TestTransformEvent_SkipsUnknownSportType(t *testing.T) {
	client := sporttransform.NewSportTransformClient()

	partial := &model.Event{
		ID:          "evt-3",
		EventTypeID: &model.OptionalString{Value: "basketball"},
	}
	full := &model.Event{
		ID:          "evt-3",
		EventTypeID: &model.OptionalString{Value: "basketball"},
		SportData:   &model.SportEvent{},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output for unknown sport type, got %#v", out)
	}
}

func TestTransformEvent_SetsSportName(t *testing.T) {
	client := sporttransform.NewSportTransformClient()

	partial := &model.Event{
		ID:          "evt-4",
		EventTypeID: &model.OptionalString{Value: "soccer"},
	}
	full := &model.Event{
		ID:          "evt-4",
		EventTypeID: &model.OptionalString{Value: "soccer"},
		SportData:   &model.SportEvent{},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatalf("expected output event, got nil")
	}
	if out.ID != "evt-4" {
		t.Fatalf("expected ID %q, got %q", "evt-4", out.ID)
	}
	if out.SportData == nil || out.SportData.Name == nil {
		t.Fatalf("expected sport name to be set, got %#v", out.SportData)
	}
	if out.SportData.Name.Value != "Soccer" {
		t.Fatalf("expected sport name %q, got %q", "Soccer", out.SportData.Name.Value)
	}
}

func TestGetName(t *testing.T) {
	client := sporttransform.NewSportTransformClient()
	if got := client.GetName(); got != "SportsTransform" {
		t.Fatalf("expected name %q, got %q", "SportsTransform", got)
	}
}
