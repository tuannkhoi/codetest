package merger_test

import (
	"context"
	"testing"

	"git.neds.sh/technology/pricekinetics/tools/codetest/merger"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func TestMergeSelectionSlice_MergesAndSorts(t *testing.T) {
	ctx := context.Background()

	left := []*model.Selection{
		{ID: "b", Name: &model.OptionalString{Value: "left-b"}},
		{ID: "d", Name: &model.OptionalString{Value: "left-d"}},
	}
	right := []*model.Selection{
		{ID: "a", Name: &model.OptionalString{Value: "right-a"}},
		{ID: "b", Name: &model.OptionalString{Value: "right-b"}},
		{ID: "c", Name: &model.OptionalString{Value: "right-c"}},
	}

	out := merger.MergeSelectionSlice(ctx, left, right)

	wantIDs := []string{"a", "b", "c", "d"}
	if len(out) != len(wantIDs) {
		t.Fatalf("expected %d selections, got %d", len(wantIDs), len(out))
	}
	for i, id := range wantIDs {
		if got := out[i].GetID(); got != id {
			t.Fatalf("expected selection %d to have ID %q, got %q", i, id, got)
		}
	}
	if got := out[1].GetName().GetValue(); got != "right-b" {
		t.Fatalf("expected merged selection name to be %q, got %q", "right-b", got)
	}
	if got := out[3].GetName().GetValue(); got != "left-d" {
		t.Fatalf("expected left-only selection name to be %q, got %q", "left-d", got)
	}
}

func TestMergeMarketSlice_MergesSelections(t *testing.T) {
	ctx := context.Background()

	left := []*model.Market{
		{
			ID: "m2",
			Selections: []*model.Selection{
				{ID: "s9", Name: &model.OptionalString{Value: "left-s9"}},
			},
		},
		{
			ID: "m1",
			Selections: []*model.Selection{
				{ID: "s1", Name: &model.OptionalString{Value: "left-s1"}},
			},
		},
	}
	right := []*model.Market{
		{
			ID: "m1",
			Selections: []*model.Selection{
				{ID: "s2", Name: &model.OptionalString{Value: "right-s2"}},
				{ID: "s1", Name: &model.OptionalString{Value: "right-s1"}},
			},
		},
		{ID: "m0"},
	}

	out := merger.MergeMarketSlice(ctx, left, right)

	wantIDs := []string{"m0", "m1", "m2"}
	if len(out) != len(wantIDs) {
		t.Fatalf("expected %d markets, got %d", len(wantIDs), len(out))
	}
	for i, id := range wantIDs {
		if got := out[i].GetID(); got != id {
			t.Fatalf("expected market %d to have ID %q, got %q", i, id, got)
		}
	}

	m1 := out[1]
	if m1.GetID() != "m1" {
		t.Fatalf("expected m1 at index 1, got %q", m1.GetID())
	}
	if len(m1.GetSelections()) != 2 {
		t.Fatalf("expected 2 selections for m1, got %d", len(m1.GetSelections()))
	}
	if got := m1.GetSelections()[0].GetID(); got != "s1" {
		t.Fatalf("expected first selection ID %q, got %q", "s1", got)
	}
	if got := m1.GetSelections()[0].GetName().GetValue(); got != "right-s1" {
		t.Fatalf("expected merged selection name to be %q, got %q", "right-s1", got)
	}
	if got := m1.GetSelections()[1].GetID(); got != "s2" {
		t.Fatalf("expected second selection ID %q, got %q", "s2", got)
	}
}
