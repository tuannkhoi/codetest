//go:build integration

package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/repository"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/service"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms/sporttransform"
	"git.neds.sh/technology/pricekinetics/tools/codetest/merger"
	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

func setupIntegrationService(t *testing.T, eventID string) (*service.Service, repository.Repository, func()) {
	t.Helper()

	repo, err := repository.NewMongoRepository(context.Background(), repository.MongoConfig{
		Host: "localhost",
		Port: 27017,
	})
	require.NoError(t, err)

	host := &service.Service{
		Upstreams: &service.Upstreams{
			MergerClient: merger.NewInlineMergerClient(),
			Repo:         repo,
			Transforms: []transforms.TransformClient{
				sporttransform.NewSportTransformClient(),
			},
		},
	}

	cleanup := func() {
		_ = repo.DeleteEventByID(context.Background(), eventID)
	}

	return host, repo, cleanup
}

func TestService_Update(t *testing.T) {
	const eventID = "integration-test-1"
	host, _, cleanup := setupIntegrationService(t, eventID)
	defer cleanup()

	output, err := host.Update(context.Background(), &core.UpdateRequest{
		Event: &model.Event{
			ID:          eventID,
			Name:        &model.OptionalString{Value: "Test event"},
			EventTypeID: &model.OptionalString{Value: "soccer"},
			StartTime:   &model.OptionalInt64{Value: 1758244443000000000}, // Friday, September 19, 2025 11:14:03 AM GMT+10:00
			EventVisibility: &model.OptionalEventVisibility{
				Value: model.EventVisibility_VisibilityDisplayed,
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("New Event born %s", eventID), output.Message)

	output, err = host.Update(context.Background(), &core.UpdateRequest{
		Event: &model.Event{
			ID: eventID,
			Markets: []*model.Market{
				{
					ID:   "mkt01",
					Name: &model.OptionalString{Value: "New Market"},
				},
			},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "Success", output.Message)

	final, err := host.GetSportEvent(context.Background(), &core.GetSportEventRequest{EventID: eventID})
	assert.NoError(t, err)
	assert.Equal(t, "Test event", final.Event.Name)
	assert.Contains(t, final.Event.StartTime, "2025-09-19")
	assert.Equal(t, "soccer", final.Event.SportTypeID)
	assert.Equal(t, "Soccer", final.Event.SportName)
	assert.Equal(t, "New Market", final.Event.Markets[0].Name.Value)
	assert.Equal(t, "VisibilityDisplayed", final.Event.EventVisibility)
}

func TestService_GetSportEvent(t *testing.T) {
	const eventID = "integration-test-2"
	host, _, cleanup := setupIntegrationService(t, eventID)
	defer cleanup()

	_, err := host.Update(context.Background(), &core.UpdateRequest{
		Event: &model.Event{
			ID:          eventID,
			Name:        &model.OptionalString{Value: "GetSportEvent"},
			EventTypeID: &model.OptionalString{Value: "soccer"},
			StartTime:   &model.OptionalInt64{Value: 1758244443000000000}, // Friday, September 19, 2025 11:14:03 AM GMT+10:00
			EventVisibility: &model.OptionalEventVisibility{
				Value: model.EventVisibility_VisibilityDisplayed,
			},
		},
	})
	require.NoError(t, err)

	resp, err := host.GetSportEvent(context.Background(), &core.GetSportEventRequest{EventID: eventID})
	require.NoError(t, err)

	assert.Equal(t, eventID, resp.Event.ID)
	assert.Equal(t, "GetSportEvent", resp.Event.Name)
	assert.Contains(t, resp.Event.StartTime, "2025-09-19")
	assert.Equal(t, "soccer", resp.Event.SportTypeID)
	assert.Equal(t, "Soccer", resp.Event.SportName)
	assert.Equal(t, "VisibilityDisplayed", resp.Event.EventVisibility)
}

func TestService_GetRaceEvent(t *testing.T) {
	const eventID = "integration-test-3"
	host, _, cleanup := setupIntegrationService(t, eventID)
	defer cleanup()

	_, err := host.Update(context.Background(), &core.UpdateRequest{
		Event: &model.Event{
			ID:        eventID,
			Name:      &model.OptionalString{Value: "GetRaceEvent"},
			StartTime: &model.OptionalInt64{Value: 1758244443000000000}, // Friday, September 19, 2025 11:14:03 AM GMT+10:00
			EventVisibility: &model.OptionalEventVisibility{
				Value: model.EventVisibility_VisibilityDisplayed,
			},
			RaceData: &model.RaceEvent{
				Category:   &model.OptionalRaceCategory{Value: model.RaceCategory_RaceCategoryHorse},
				Distance:   &model.OptionalInt64{Value: 1200},
				RaceCourse: &model.OptionalString{Value: "flemington"},
				State:      &model.OptionalString{Value: "VIC"},
			},
		},
	})
	require.NoError(t, err)

	resp, err := host.GetRaceEvent(context.Background(), &core.GetRaceEventRequest{EventID: eventID})
	require.NoError(t, err)

	assert.Equal(t, eventID, resp.Event.ID)
	assert.Equal(t, "GetRaceEvent", resp.Event.Name)
	assert.Contains(t, resp.Event.StartTime, "2025-09-19")
	assert.Equal(t, "RaceCategoryHorse", resp.Event.Category)
	assert.Equal(t, int64(1200), resp.Event.Distance)
	assert.Equal(t, "flemington", resp.Event.RaceCourse)
	assert.Equal(t, "VIC", resp.Event.State)
	assert.Equal(t, "VisibilityDisplayed", resp.Event.EventVisibility)
}

func TestService_SearchEvents(t *testing.T) {
	const seedPrefix = "integration-search-"
	host, repo, _ := setupIntegrationService(t, "")

	seed := seedEvents(t, host, repo, seedPrefix, 15)
	defer func() {
		for _, id := range seed {
			_ = repo.DeleteEventByID(context.Background(), id)
		}
	}()

	baseTime := time.Date(2025, 9, 19, 1, 14, 3, 0, time.UTC)
	filter := &core.SearchEventsFilter{
		BettingStatus:   model.BettingStatus_BettingOpen.Enum(),
		EventVisibility: model.EventVisibility_VisibilityDisplayed.Enum(),
		StartDate:       timestamppb.New(baseTime.Add(2 * time.Minute)),
		EndDate:         timestamppb.New(baseTime.Add(10 * time.Minute)),
	}

	pageSize := uint64(5)
	resp, err := host.SearchEvents(context.Background(), &core.SearchEventsRequest{
		Filter:   filter,
		PageSize: &pageSize,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp.SportEvents)
	assert.NotEmpty(t, resp.RaceEvents)
	assert.NotEmpty(t, resp.NextPageToken)

	pageToken := resp.NextPageToken
	resp2, err := host.SearchEvents(context.Background(), &core.SearchEventsRequest{
		Filter:    filter,
		PageSize:  &pageSize,
		PageToken: &pageToken,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resp2.SportEvents)
	assert.NotEmpty(t, resp2.RaceEvents)
	if resp2.NextPageToken != "" {
		assert.NotEqual(t, resp.NextPageToken, resp2.NextPageToken)
	}
}

func seedEvents(
	t *testing.T,
	host *service.Service,
	repo repository.Repository,
	prefix string,
	count int,
) []string {
	t.Helper()

	var ids []string
	baseTime := time.Date(2025, 9, 19, 1, 14, 3, 0, time.UTC)
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("%s%02d", prefix, i)
		event := &model.Event{
			ID:        id,
			Name:      &model.OptionalString{Value: "Search Seed"},
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

		_, err := host.Update(context.Background(), &core.UpdateRequest{Event: event})
		if err != nil {
			t.Fatalf("failed to seed event: %v", err)
		}

		ids = append(ids, id)
	}

	return ids
}
