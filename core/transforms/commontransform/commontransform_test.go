package commontransform_test

import (
	"context"
	"testing"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms/commontransform"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func TestTransformEvent_SkipsWhenNoMarketClosed(t *testing.T) {
	client := commontransform.NewCommonTransformClient()

	partial := &model.Event{
		ID: "evt-1",
		Markets: []*model.Market{
			{ID: "m1", BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen}},
		},
	}
	full := &model.Event{
		ID: "evt-1",
		Markets: []*model.Market{
			{ID: "m1", BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen}},
		},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output when no market closed, got %#v", out)
	}
}

func TestTransformEvent_SkipsWhenClosedAtAlreadySet(t *testing.T) {
	client := commontransform.NewCommonTransformClient()

	partial := &model.Event{
		ID: "evt-2",
		Markets: []*model.Market{
			{ID: "m1", BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingClosed}},
		},
	}
	full := &model.Event{
		ID: "evt-2",
		Markets: []*model.Market{
			{
				ID:            "m1",
				BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingClosed},
				ClosedAt:      &model.OptionalInt64{Value: 123},
			},
		},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != nil {
		t.Fatalf("expected nil output when ClosedAt already set, got %#v", out)
	}
}

func TestTransformEvent_SetsClosedAtForClosedMarkets(t *testing.T) {
	client := commontransform.NewCommonTransformClient()

	partial := &model.Event{
		ID: "evt-3",
		Markets: []*model.Market{
			{ID: "m1", BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingClosed}},
		},
	}
	full := &model.Event{
		ID: "evt-3",
		Markets: []*model.Market{
			{ID: "m1", BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingClosed}},
		},
	}

	out, err := client.TransformEvent(context.Background(), partial, full)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatalf("expected output event, got nil")
	}
	if out.ID != "evt-3" {
		t.Fatalf("expected ID %q, got %q", "evt-3", out.ID)
	}
	if len(out.Markets) != 1 || out.Markets[0].GetID() != "m1" {
		t.Fatalf("expected market ID %q, got %#v", "m1", out.Markets)
	}
	if out.Markets[0].GetClosedAt() == nil {
		t.Fatalf("expected ClosedAt to be set, got nil")
	}
	if out.Markets[0].GetClosedAt().GetValue() == 0 {
		t.Fatalf("expected ClosedAt value to be set")
	}
}

func TestGetName(t *testing.T) {
	client := commontransform.NewCommonTransformClient()
	if got := client.GetName(); got != "CommonTransform" {
		t.Fatalf("expected name %q, got %q", "CommonTransform", got)
	}
}
