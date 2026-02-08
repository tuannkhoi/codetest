//go:build integration

package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	repo, err := repository.NewRedisRepository(context.Background(), "localhost:6379", "")
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
