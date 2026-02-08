package merger_test

import (
	"context"
	"testing"

	"git.neds.sh/technology/pricekinetics/tools/codetest/merger"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func TestMergeOptionalString(t *testing.T) {
	left := &model.OptionalString{Value: "left", Deleted: true}
	right := &model.OptionalString{Value: "right", Deleted: false}

	if got := merger.MergeOptionalString(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeOptionalString(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeOptionalString(context.Background(), left, right)
	if out.Value != "right" || out.Deleted != false {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeOptionalDouble(t *testing.T) {
	left := &model.OptionalDouble{Value: 1.25, Deleted: true}
	right := &model.OptionalDouble{Value: 2.5, Deleted: false}

	if got := merger.MergeOptionalDouble(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeOptionalDouble(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeOptionalDouble(context.Background(), left, right)
	if out.Value != 2.5 || out.Deleted != false {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeOptionalInt64(t *testing.T) {
	left := &model.OptionalInt64{Value: 1, Deleted: true}
	right := &model.OptionalInt64{Value: 2, Deleted: false}

	if got := merger.MergeOptionalInt64(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeOptionalInt64(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeOptionalInt64(context.Background(), left, right)
	if out.Value != 2 || out.Deleted != false {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeOptionalBettingStatus(t *testing.T) {
	left := &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen, Deleted: true}
	right := &model.OptionalBettingStatus{Value: model.BettingStatus_BettingClosed, Deleted: false}

	if got := merger.MergeOptionalBettingStatus(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeOptionalBettingStatus(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeOptionalBettingStatus(context.Background(), left, right)
	if out.Value != model.BettingStatus_BettingClosed || out.Deleted != false {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeOptionalEventVisibility(t *testing.T) {
	left := &model.OptionalEventVisibility{Value: model.EventVisibility_VisibilityDisplayed, Deleted: true}
	right := &model.OptionalEventVisibility{Value: model.EventVisibility_VisibilityHidden, Deleted: false}

	if got := merger.MergeOptionalEventVisibility(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeOptionalEventVisibility(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeOptionalEventVisibility(context.Background(), left, right)
	if out.Value != model.EventVisibility_VisibilityHidden || out.Deleted != false {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeOptionalRaceCategory(t *testing.T) {
	left := &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryGreyhound, Deleted: true}
	right := &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryHorse, Deleted: false}

	if got := merger.MergeOptionalRaceCategory(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeOptionalRaceCategory(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeOptionalRaceCategory(context.Background(), left, right)
	if out.Value != model.RaceCategory_RaceCategoryHorse || out.Deleted != false {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeSportEvent(t *testing.T) {
	left := &model.SportEvent{
		Name:   &model.OptionalString{Value: "Left"},
		Region: &model.OptionalString{Value: "LeftRegion"},
		League: &model.OptionalString{Value: "LeftLeague"},
		Round:  &model.OptionalString{Value: "LeftRound"},
	}
	right := &model.SportEvent{
		Name:   &model.OptionalString{Value: "Right"},
		Region: &model.OptionalString{Value: "RightRegion"},
		League: &model.OptionalString{Value: "RightLeague"},
		Round:  &model.OptionalString{Value: "RightRound"},
	}

	if got := merger.MergeSportEvent(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeSportEvent(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeSportEvent(context.Background(), left, right)
	if out.Name.Value != "Right" || out.Region.Value != "RightRegion" || out.League.Value != "RightLeague" ||
		out.Round.Value != "RightRound" {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeRaceEvent(t *testing.T) {
	left := &model.RaceEvent{
		Category:   &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryGreyhound},
		Distance:   &model.OptionalInt64{Value: 900},
		RaceCourse: &model.OptionalString{Value: "left-course"},
		State:      &model.OptionalString{Value: "QLD"},
	}
	right := &model.RaceEvent{
		Category:   &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryHorse},
		Distance:   &model.OptionalInt64{Value: 1200},
		RaceCourse: &model.OptionalString{Value: "right-course"},
		State:      &model.OptionalString{Value: "VIC"},
	}

	if got := merger.MergeRaceEvent(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeRaceEvent(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeRaceEvent(context.Background(), left, right)
	if out.Category.Value != model.RaceCategory_RaceCategoryHorse ||
		out.Distance.Value != 1200 ||
		out.RaceCourse.Value != "right-course" ||
		out.State.Value != "VIC" {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeSelection(t *testing.T) {
	left := &model.Selection{
		ID:            "sel-1",
		Name:          &model.OptionalString{Value: "Left"},
		BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen},
		Price:         &model.OptionalDouble{Value: 1.1},
	}
	right := &model.Selection{
		ID:            "sel-1",
		Name:          &model.OptionalString{Value: "Right"},
		BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingClosed},
		Price:         &model.OptionalDouble{Value: 2.2},
	}

	if got := merger.MergeSelection(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeSelection(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeSelection(context.Background(), left, right)
	if out.ID != "sel-1" {
		t.Fatalf("expected ID %q, got %q", "sel-1", out.ID)
	}
	if out.Name.Value != "Right" || out.BettingStatus.Value != model.BettingStatus_BettingClosed ||
		out.Price.Value != 2.2 {
		t.Fatalf("expected right values, got %+v", out)
	}
}

func TestMergeMarket(t *testing.T) {
	left := &model.Market{
		ID:            "mkt-1",
		Name:          &model.OptionalString{Value: "Left"},
		StartTime:     &model.OptionalInt64{Value: 1},
		ClosedAt:      &model.OptionalInt64{Value: 10},
		BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen},
		Selections: []*model.Selection{
			{ID: "s1", Name: &model.OptionalString{Value: "LeftS1"}},
		},
	}
	right := &model.Market{
		ID:            "mkt-1",
		Name:          &model.OptionalString{Value: "Right"},
		StartTime:     &model.OptionalInt64{Value: 2},
		ClosedAt:      &model.OptionalInt64{Value: 20},
		BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingClosed},
		Selections: []*model.Selection{
			{ID: "s1", Name: &model.OptionalString{Value: "RightS1"}},
			{ID: "s2", Name: &model.OptionalString{Value: "RightS2"}},
		},
	}

	if got := merger.MergeMarket(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeMarket(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeMarket(context.Background(), left, right)
	if out.ID != "mkt-1" {
		t.Fatalf("expected ID %q, got %q", "mkt-1", out.ID)
	}
	if out.Name.Value != "Right" || out.StartTime.Value != 2 ||
		out.ClosedAt.Value != 20 || out.BettingStatus.Value != model.BettingStatus_BettingClosed {
		t.Fatalf("expected right values, got %+v", out)
	}
	if len(out.Selections) != 2 {
		t.Fatalf("expected 2 selections, got %d", len(out.Selections))
	}
	if out.Selections[0].ID != "s1" || out.Selections[0].Name.Value != "RightS1" {
		t.Fatalf("expected merged selection s1, got %+v", out.Selections[0])
	}
	if out.Selections[1].ID != "s2" {
		t.Fatalf("expected selection s2, got %+v", out.Selections[1])
	}
}

func TestMergeEvent(t *testing.T) {
	left := &model.Event{
		ID:          "evt-1",
		Name:        &model.OptionalString{Value: "Left"},
		EventTypeID: &model.OptionalString{Value: "soccer"},
		StartTime:   &model.OptionalInt64{Value: 1},
		BettingStatus: &model.OptionalBettingStatus{
			Value: model.BettingStatus_BettingOpen,
		},
		EventVisibility: &model.OptionalEventVisibility{
			Value: model.EventVisibility_VisibilityDisplayed,
		},
		SportData: &model.SportEvent{
			Name:   &model.OptionalString{Value: "LeftSport"},
			League: &model.OptionalString{Value: "LeftLeague"},
		},
		RaceData: &model.RaceEvent{
			Category:   &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryGreyhound},
			Distance:   &model.OptionalInt64{Value: 900},
			RaceCourse: &model.OptionalString{Value: "left-course"},
			State:      &model.OptionalString{Value: "QLD"},
		},
		Markets: []*model.Market{
			{ID: "m1", Name: &model.OptionalString{Value: "LeftM1"}, ClosedAt: &model.OptionalInt64{Value: 10}},
		},
	}
	right := &model.Event{
		ID:          "evt-1",
		Name:        &model.OptionalString{Value: "Right"},
		EventTypeID: &model.OptionalString{Value: "soccer"},
		StartTime:   &model.OptionalInt64{Value: 2},
		BettingStatus: &model.OptionalBettingStatus{
			Value: model.BettingStatus_BettingClosed,
		},
		EventVisibility: &model.OptionalEventVisibility{
			Value: model.EventVisibility_VisibilityHidden,
		},
		SportData: &model.SportEvent{
			Name:   &model.OptionalString{Value: "RightSport"},
			League: &model.OptionalString{Value: "RightLeague"},
		},
		RaceData: &model.RaceEvent{
			Category:   &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryHorse},
			Distance:   &model.OptionalInt64{Value: 1200},
			RaceCourse: &model.OptionalString{Value: "right-course"},
			State:      &model.OptionalString{Value: "VIC"},
		},
		Markets: []*model.Market{
			{ID: "m1", Name: &model.OptionalString{Value: "RightM1"}, ClosedAt: &model.OptionalInt64{Value: 20}},
			{ID: "m2", Name: &model.OptionalString{Value: "RightM2"}},
		},
	}

	if got := merger.MergeEvent(context.Background(), nil, right); got != right {
		t.Fatalf("expected right when left nil")
	}
	if got := merger.MergeEvent(context.Background(), left, nil); got != left {
		t.Fatalf("expected left when right nil")
	}

	out := merger.MergeEvent(context.Background(), left, right)
	if out.ID != "evt-1" {
		t.Fatalf("expected ID %q, got %q", "evt-1", out.ID)
	}
	if out.Name.Value != "Right" || out.StartTime.Value != 2 ||
		out.BettingStatus.Value != model.BettingStatus_BettingClosed ||
		out.EventVisibility.Value != model.EventVisibility_VisibilityHidden {
		t.Fatalf("expected right values, got %+v", out)
	}
	if out.SportData.Name.Value != "RightSport" || out.SportData.League.Value != "RightLeague" {
		t.Fatalf("expected right sport data, got %+v", out.SportData)
	}
	if out.RaceData.Category.Value != model.RaceCategory_RaceCategoryHorse ||
		out.RaceData.Distance.Value != 1200 ||
		out.RaceData.RaceCourse.Value != "right-course" ||
		out.RaceData.State.Value != "VIC" {
		t.Fatalf("expected right race data, got %+v", out.RaceData)
	}
	if len(out.Markets) != 2 {
		t.Fatalf("expected 2 markets, got %d", len(out.Markets))
	}
	if out.Markets[0].ID != "m1" || out.Markets[0].Name.Value != "RightM1" ||
		out.Markets[0].ClosedAt.Value != 20 {
		t.Fatalf("expected merged market m1, got %+v", out.Markets[0])
	}
	if out.Markets[1].ID != "m2" {
		t.Fatalf("expected market m2, got %+v", out.Markets[1])
	}
}
