// Package service contains the grpc/http server implementation
package service

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"git.neds.sh/technology/pricekinetics/tools/codetest/core"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/repository"
	"git.neds.sh/technology/pricekinetics/tools/codetest/core/transforms"
	"git.neds.sh/technology/pricekinetics/tools/codetest/merger"
)

// Upstreams defines dependencies the service has on other services
type Upstreams struct {
	MergerClient merger.ServiceClient
	Repo         repository.Repository
	Transforms   []transforms.TransformClient
}

// NewService creates a new instance of Service
func NewService(grpcPort, httpPort int, upstreams *Upstreams) *Service {
	return &Service{GRPCPort: grpcPort, HTTPPort: httpPort, Upstreams: upstreams}
}

// RegisterGRPCServerImplementations registers the grpc service contract implemented by this server
func (host *Service) RegisterGRPCServerImplementations(grpcServer *grpc.Server) {
	core.RegisterServiceServer(grpcServer, host)
}

// Update updates an Event and runs the pipeline of transformations
func (host *Service) Update(ctx context.Context, req *core.UpdateRequest) (*core.UpdateResponse, error) {
	existing, err := host.Upstreams.Repo.GetEventByID(ctx, req.GetEvent().GetID())
	if err != nil {
		logrus.WithError(err).Error("Update: failed to retrieve event")
		return nil, err
	}

	resp := &core.UpdateResponse{Message: "Success"}

	update := req.GetEvent()
	if existing == nil {
		resp.Message = fmt.Sprintf("New Event born %v", req.GetEvent().GetID())
	} else {
		update, err = host.Upstreams.MergerClient.MergeEvent(context.Background(), existing, req.GetEvent())
		if err != nil {
			logrus.WithError(err).Error("Update: failed to merge event")
			return nil, err
		}
	}

	for _, t := range host.Upstreams.Transforms {
		upd, tErr := t.TransformEvent(ctx, req.Event, update)
		if tErr != nil {
			logrus.WithError(tErr).Errorf("Update: failed to run transform %v", t.GetName())
		}
		if upd != nil {
			update, err = host.Upstreams.MergerClient.MergeEvent(context.Background(), update, upd)
			if err != nil {
				logrus.WithError(err).Errorf("Update: failed to merge event in transform %v", t.GetName())
				return nil, err
			}
		}
	}

	if err := host.Upstreams.Repo.UpdateEvent(ctx, update); err != nil {
		logrus.WithError(err).Error("Update: failed to update event")
		return nil, err
	}

	return resp, nil
}

// GetSportEvent retrieves a model.Event from the database and returns a core.SportEvent,
// this is a more UserConsumable representation of the model that is specific to sport events
func (host *Service) GetSportEvent(ctx context.Context, req *core.GetSportEventRequest) (
	*core.GetSportEventResponse, error,
) {
	existing, err := host.Upstreams.Repo.GetEventByID(ctx, req.GetEventID())
	if err != nil {
		logrus.WithError(err).Error("GetSportEvent: failed to retrieve event")
		return nil, err
	}

	resp := &core.GetSportEventResponse{}

	if existing == nil {
		return resp, nil
	}

	sportEvent := &core.SportEvent{}
	sportEvent.ConvertFromModel(existing)
	resp.Event = sportEvent

	return resp, nil
}

// GetRaceEvent retrieves a model.Event from the database and returns a core.RaceEvent,
// this is a more UserConsumable representation of the model that is specific to race events
func (host *Service) GetRaceEvent(ctx context.Context, req *core.GetRaceEventRequest) (
	*core.GetRaceEventResponse, error,
) {
	existing, err := host.Upstreams.Repo.GetEventByID(ctx, req.GetEventID())
	if err != nil {
		logrus.WithError(err).Error("GetRaceEvent: failed to retrieve event")
		return nil, err
	}

	resp := &core.GetRaceEventResponse{}

	if existing == nil {
		return resp, nil
	}

	raceEvent := &core.RaceEvent{}
	raceEvent.ConvertFromModel(existing)
	resp.Event = raceEvent

	return resp, nil
}

// SearchEvents allows users to search events by date and/or betting status and/or event visibility.
func (host *Service) SearchEvents(ctx context.Context, req *core.SearchEventsRequest) (
	*core.SearchEventsResponse, error,
) {
	if err := req.Validate(); err != nil {
		logrus.WithError(err).Error("SearchEvents: failed to validate request")
		return nil, err
	}

	events, nextToken, err := host.Upstreams.Repo.SearchEvents(ctx, req.GetFilter(), req.GetPageSize(), req.GetPageToken())
	if err != nil {
		logrus.WithError(err).Error("SearchEvents: failed to retrieve events")
		return nil, err
	}

	var raceEvents []*core.RaceEvent
	var sportEvents []*core.SportEvent

	for _, event := range events {
		// at the moment, an event being inserted is not required to have either RaceEvent or SportEvent
		// but in the real world (and also for the purpose of this task), it's reasonable to assume that
		// we can tell which event is race or sport by checking whether race or sport data is present
		if event.HasRaceData() {
			raceEvent := &core.RaceEvent{}
			raceEvent.ConvertFromModel(event)
			raceEvents = append(raceEvents, raceEvent)
		} else {
			sportEvent := &core.SportEvent{}
			sportEvent.ConvertFromModel(event)
			sportEvents = append(sportEvents, sportEvent)
		}
	}

	return &core.SearchEventsResponse{
		RaceEvents:    raceEvents,
		SportEvents:   sportEvents,
		NextPageToken: nextToken,
	}, nil
}
