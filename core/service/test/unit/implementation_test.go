package service_test

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/repository/mock"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/service"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms/sporttransform"
	"git.neds.sh/technology/pricekinetics/tools/codetest/merger"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func TestService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRepository(ctrl)
	host := &service.Service{
		Upstreams: &service.Upstreams{
			MergerClient: merger.NewInlineMergerClient(),
			Repo:         repo,
			Transforms: []transforms.TransformClient{
				sporttransform.NewSportTransformClient(),
			},
		},
	}

	ctx := context.Background()

	newEvent := &model.Event{
		ID:          "unit-test-1",
		Name:        &model.OptionalString{Value: "Test event"},
		EventTypeID: &model.OptionalString{Value: "soccer"},
		StartTime:   &model.OptionalInt64{Value: 1758244443000000000},
		EventVisibility: &model.OptionalEventVisibility{
			Value: model.EventVisibility_VisibilityDisplayed,
		},
	}

	repo.EXPECT().GetEventByID(ctx, newEvent.ID).Return(nil, nil)
	repo.EXPECT().UpdateEvent(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, evt *model.Event) error {
		if evt.GetID() != newEvent.ID {
			t.Fatalf("expected event ID %q, got %q", newEvent.ID, evt.GetID())
		}
		if evt.GetName().GetValue() != "Test event" {
			t.Fatalf("expected name %q, got %q", "Test event", evt.GetName().GetValue())
		}
		if evt.GetSportData().GetName().GetValue() != "Soccer" {
			t.Fatalf("expected sport name %q, got %q", "Soccer", evt.GetSportData().GetName().GetValue())
		}
		if evt.GetEventVisibility().GetValue() != model.EventVisibility_VisibilityDisplayed {
			t.Fatalf("expected visibility %q, got %q", "VisibilityDisplayed", evt.GetEventVisibility().GetValue().String())
		}
		return nil
	})

	_, err := host.Update(ctx, &core.UpdateRequest{Event: newEvent})
	if err != nil {
		t.Fatalf("unexpected error on new event: %v", err)
	}

	existing := &model.Event{
		ID:          "unit-test-2",
		Name:        &model.OptionalString{Value: "Old name"},
		EventTypeID: &model.OptionalString{Value: "soccer"},
		StartTime:   &model.OptionalInt64{Value: 1758244443000000000},
		EventVisibility: &model.OptionalEventVisibility{
			Value: model.EventVisibility_VisibilityDisplayed,
		},
	}
	update := &model.Event{
		ID:   existing.ID,
		Name: &model.OptionalString{Value: "Updated name"},
	}

	repo.EXPECT().GetEventByID(ctx, existing.ID).Return(existing, nil)
	repo.EXPECT().UpdateEvent(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, evt *model.Event) error {
		if evt.GetID() != existing.ID {
			t.Fatalf("expected event ID %q, got %q", existing.ID, evt.GetID())
		}
		if evt.GetName().GetValue() != "Updated name" {
			t.Fatalf("expected name %q, got %q", "Updated name", evt.GetName().GetValue())
		}
		if evt.GetEventTypeID().GetValue() != "soccer" {
			t.Fatalf("expected event type %q, got %q", "soccer", evt.GetEventTypeID().GetValue())
		}
		if evt.GetEventVisibility().GetValue() != model.EventVisibility_VisibilityDisplayed {
			t.Fatalf("expected visibility %q, got %q", "VisibilityDisplayed", evt.GetEventVisibility().GetValue().String())
		}
		return nil
	})

	_, err = host.Update(ctx, &core.UpdateRequest{Event: update})
	if err != nil {
		t.Fatalf("unexpected error on existing event: %v", err)
	}
}

func TestService_GetSportEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRepository(ctrl)
	host := &service.Service{
		Upstreams: &service.Upstreams{
			MergerClient: merger.NewInlineMergerClient(),
			Repo:         repo,
			Transforms: []transforms.TransformClient{
				sporttransform.NewSportTransformClient(),
			},
		},
	}

	ctx := context.Background()
	event := &model.Event{
		ID:          "unit-get-1",
		Name:        &model.OptionalString{Value: "GetSportEvent"},
		StartTime:   &model.OptionalInt64{Value: 1758244443000000000},
		EventTypeID: &model.OptionalString{Value: "soccer"},
		BettingStatus: &model.OptionalBettingStatus{
			Value: model.BettingStatus_BettingOpen,
		},
		SportData: &model.SportEvent{
			Name:   &model.OptionalString{Value: "Soccer"},
			League: &model.OptionalString{Value: "Premier League"},
			Round:  &model.OptionalString{Value: "Round 1"},
			Region: &model.OptionalString{Value: "EU"},
		},
		EventVisibility: &model.OptionalEventVisibility{
			Value: model.EventVisibility_VisibilityDisplayed,
		},
		Markets: []*model.Market{
			{ID: "m1"},
		},
	}

	repo.EXPECT().GetEventByID(ctx, event.ID).Return(event, nil)

	resp, err := host.GetSportEvent(ctx, &core.GetSportEventRequest{EventID: event.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Event.ID != event.ID {
		t.Fatalf("expected event ID %q, got %q", event.ID, resp.Event.ID)
	}
	if resp.Event.Name != "GetSportEvent" {
		t.Fatalf("expected name %q, got %q", "GetSportEvent", resp.Event.Name)
	}
	if resp.Event.SportTypeID != "soccer" {
		t.Fatalf("expected sport type %q, got %q", "soccer", resp.Event.SportTypeID)
	}
	if resp.Event.SportName != "Soccer" {
		t.Fatalf("expected sport name %q, got %q", "Soccer", resp.Event.SportName)
	}
	if resp.Event.EventVisibility != "VisibilityDisplayed" {
		t.Fatalf("expected visibility %q, got %q", "VisibilityDisplayed", resp.Event.EventVisibility)
	}
	if resp.Event.League != "Premier League" {
		t.Fatalf("expected league %q, got %q", "Premier League", resp.Event.League)
	}
	if resp.Event.Round != "Round 1" {
		t.Fatalf("expected round %q, got %q", "Round 1", resp.Event.Round)
	}
	if resp.Event.Region != "EU" {
		t.Fatalf("expected region %q, got %q", "EU", resp.Event.Region)
	}
	if len(resp.Event.Markets) != 1 || resp.Event.Markets[0].ID != "m1" {
		t.Fatalf("expected market ID %q, got %#v", "m1", resp.Event.Markets)
	}
}

func TestService_GetRaceEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRepository(ctrl)
	host := &service.Service{
		Upstreams: &service.Upstreams{
			MergerClient: merger.NewInlineMergerClient(),
			Repo:         repo,
			Transforms: []transforms.TransformClient{
				sporttransform.NewSportTransformClient(),
			},
		},
	}

	ctx := context.Background()
	event := &model.Event{
		ID:        "unit-race-1",
		Name:      &model.OptionalString{Value: "GetRaceEvent"},
		StartTime: &model.OptionalInt64{Value: 1758244443000000000},
		BettingStatus: &model.OptionalBettingStatus{
			Value: model.BettingStatus_BettingOpen,
		},
		EventVisibility: &model.OptionalEventVisibility{
			Value: model.EventVisibility_VisibilityDisplayed,
		},
		RaceData: &model.RaceEvent{
			Category:   &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryHorse},
			Distance:   &model.OptionalInt64{Value: 1200},
			RaceCourse: &model.OptionalString{Value: "flemington"},
			State:      &model.OptionalString{Value: "VIC"},
		},
		Markets: []*model.Market{
			{ID: "m1"},
		},
	}

	repo.EXPECT().GetEventByID(ctx, event.ID).Return(event, nil)

	resp, err := host.GetRaceEvent(ctx, &core.GetRaceEventRequest{EventID: event.ID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Event.ID != event.ID {
		t.Fatalf("expected event ID %q, got %q", event.ID, resp.Event.ID)
	}
	if resp.Event.Name != "GetRaceEvent" {
		t.Fatalf("expected name %q, got %q", "GetRaceEvent", resp.Event.Name)
	}
	if resp.Event.Category != "RaceCategoryHorse" {
		t.Fatalf("expected category %q, got %q", "RaceCategoryHorse", resp.Event.Category)
	}
	if resp.Event.Distance != 1200 {
		t.Fatalf("expected distance %d, got %d", 1200, resp.Event.Distance)
	}
	if resp.Event.RaceCourse != "flemington" {
		t.Fatalf("expected race course %q, got %q", "flemington", resp.Event.RaceCourse)
	}
	if resp.Event.State != "VIC" {
		t.Fatalf("expected state %q, got %q", "VIC", resp.Event.State)
	}
	if resp.Event.EventVisibility != "VisibilityDisplayed" {
		t.Fatalf("expected visibility %q, got %q", "VisibilityDisplayed", resp.Event.EventVisibility)
	}
	if len(resp.Event.Markets) != 1 || resp.Event.Markets[0].ID != "m1" {
		t.Fatalf("expected market ID %q, got %#v", "m1", resp.Event.Markets)
	}
}
