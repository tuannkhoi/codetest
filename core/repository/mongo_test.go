//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/repository"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func TestMongoRepo_UpdateGetDeleteEvent(t *testing.T) {
	cfg := repository.MongoConfig{
		Host: "localhost",
		Port: 27017,
	}

	repo, err := repository.NewMongoRepository(context.Background(), cfg)
	require.NoError(t, err)

	input := &model.Event{
		ID:            "mongo-e001",
		Name:          &model.OptionalString{Value: "Test Event"},
		StartTime:     &model.OptionalInt64{Value: 1758244443000000000}, // Friday, September 19, 2025 11:14:03 AM GMT+10:00
		BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen},
		EventTypeID:   &model.OptionalString{Value: "rugby_league"},
		EventVisibility: &model.OptionalEventVisibility{
			Value: model.EventVisibility_VisibilityDisplayed,
		},
		Markets: []*model.Market{
			{
				ID:            "mkt-1",
				Name:          &model.OptionalString{Value: "Head to Head"},
				StartTime:     &model.OptionalInt64{Value: 1758244443000000000}, // Friday, September 19, 2025 11:14:03 AM GMT+10:00
				BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen},
				Selections: []*model.Selection{
					{
						ID:            "sel-1",
						Name:          &model.OptionalString{Value: "Home Team"},
						BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingOpen},
						Price:         &model.OptionalDouble{Value: 1.80},
					},
					{
						ID:            "sel-2",
						Name:          &model.OptionalString{Value: "Away Team"},
						BettingStatus: &model.OptionalBettingStatus{Value: model.BettingStatus_BettingSuspended},
						Price:         &model.OptionalDouble{Value: 1.80},
					},
				},
			},
		},
	}

	defer func() {
		_ = repo.DeleteEventByID(context.Background(), input.ID)
	}()

	updErr := repo.UpdateEvent(context.Background(), input)
	assert.NoError(t, updErr)

	output, getErr := repo.GetEventByID(context.Background(), input.ID)
	assert.NoError(t, getErr)
	assert.Equal(t, input.Name.Value, output.Name.Value)
	assert.Equal(t, input.BettingStatus.Value, output.BettingStatus.Value)
	assert.Equal(t, input.EventVisibility.Value, output.EventVisibility.Value)
	assert.Equal(t, input.Markets[0].Name.Value, output.Markets[0].Name.Value)
	assert.Equal(t, input.Markets[0].Selections[0].Name.Value, output.Markets[0].Selections[0].Name.Value)

	delErr := repo.DeleteEventByID(context.Background(), input.ID)
	assert.NoError(t, delErr)
}

func TestMongoRepo_SearchEvents(t *testing.T) {
	cfg := repository.MongoConfig{
		Host: "localhost",
		Port: 27017,
	}
	repo, err := repository.NewMongoRepository(context.Background(), cfg)
	require.NoError(t, err)

	seed := seedEvents(t, repo, 200)
	defer func() {
		for _, id := range seed {
			_ = repo.DeleteEventByID(context.Background(), id)
		}
	}()

	baseTime := time.Date(2025, 9, 19, 1, 14, 3, 0, time.UTC)
	start := baseTime.Add(2 * time.Minute)
	end := baseTime.Add(10 * time.Minute)

	filter := &core.SearchEventsFilter{
		BettingStatus:   model.BettingStatus_BettingOpen.Enum(),
		EventVisibility: model.EventVisibility_VisibilityDisplayed.Enum(),
		StartDate:       timestamppb.New(start),
		EndDate:         timestamppb.New(end),
	}

	events, nextToken, err := repo.SearchEvents(context.Background(), filter, 5, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, events)
	assert.NotEmpty(t, nextToken)

	more, nextToken2, err := repo.SearchEvents(context.Background(), filter, 5, nextToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, more)
	if nextToken2 != "" {
		assert.NotEqual(t, nextToken, nextToken2)
	}
}

func seedEvents(t *testing.T, repo repository.Repository, count int) []string {
	t.Helper()

	var ids []string
	baseTime := time.Date(2025, 9, 19, 1, 14, 3, 0, time.UTC)
	for i := 0; i < count; i++ {
		id := "mongo-seed-" + time.Now().Format("150405") + "-" + time.Duration(i).String()
		event := &model.Event{
			ID:        id,
			Name:      &model.OptionalString{Value: "Seed Event"},
			StartTime: &model.OptionalInt64{Value: baseTime.Add(time.Duration(i) * time.Minute).UnixNano()},
			BettingStatus: &model.OptionalBettingStatus{
				Value: model.BettingStatus_BettingOpen,
			},
			EventVisibility: &model.OptionalEventVisibility{
				Value: model.EventVisibility_VisibilityDisplayed,
			},
		}

		if i%2 == 0 {
			event.EventTypeID = &model.OptionalString{Value: "soccer"}
			event.SportData = &model.SportEvent{
				Name: &model.OptionalString{Value: "Soccer"},
			}
		} else {
			event.RaceData = &model.RaceEvent{
				Category:   &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryHorse},
				Distance:   &model.OptionalInt64{Value: 1200},
				RaceCourse: &model.OptionalString{Value: "flemington"},
				State:      &model.OptionalString{Value: "VIC"},
			}
		}

		if i%5 == 0 {
			event.BettingStatus.Value = model.BettingStatus_BettingClosed
		}
		if i%5 == 0 {
			event.EventVisibility.Value = model.EventVisibility_VisibilityHidden
		}

		if err := repo.UpdateEvent(context.Background(), event); err != nil {
			t.Fatalf("failed to seed event: %v", err)
		}
		ids = append(ids, id)
	}

	return ids
}
