// Package core contains the proto definitions of the Core Service
package core

import (
	"time"

	"git.neds.sh/technology/pricekinetics/tools/codetest/model"
)

//go:generate ./gen-proto.sh

// baseEventFieldsSetters are implemented by all type of events that contain below fields
type baseEventFieldsSetter interface {
	SetID(string)
	SetName(string)
	SetStartTime(string)
	SetBettingStatus(string)
	SetEventVisibility(string)
	SetMarkets([]*model.Market)
}

// fillBaseEventFields is a helper functions for ConvertFromModel functions below
// so that the logic for setting similar fields will not be repeated
func fillBaseEventFields(to baseEventFieldsSetter, m *model.Event) {
	to.SetID(m.GetID())
	to.SetName(m.GetName().GetValue())
	to.SetStartTime(time.Unix(0, m.GetStartTime().GetValue()).Format(time.RFC3339))
	to.SetBettingStatus(m.GetBettingStatus().GetValue().String())
	to.SetEventVisibility(m.GetEventVisibility().GetValue().String())
	to.SetMarkets(m.GetMarkets())
}

// ConvertFromModel converts a model.Event to a core.SportEvent
func (to *SportEvent) ConvertFromModel(model *model.Event) {
	fillBaseEventFields(to, model)

	to.SportTypeID = model.GetEventTypeID().GetValue()
	to.League = model.GetSportData().GetLeague().GetValue()
	to.SportName = model.GetSportData().GetName().GetValue()
	to.Round = model.GetSportData().GetRound().GetValue()
	to.Region = model.GetSportData().GetRegion().GetValue()
}

// ConvertFromModel converts a model.Event to a core.RaceEvent.
func (to *RaceEvent) ConvertFromModel(model *model.Event) {
	fillBaseEventFields(to, model)
	to.Category = model.GetRaceData().GetCategory().GetValue().String()
	to.Distance = model.GetRaceData().GetDistance().GetValue()
	to.RaceCourse = model.GetRaceData().GetRaceCourse().GetValue()
	to.State = model.GetRaceData().GetState().GetValue()
}
