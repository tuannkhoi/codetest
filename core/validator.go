package core

import "errors"

// Validate validates SearchEventsRequest
func (req *SearchEventsRequest) Validate() error {
	if !req.HasFilter() {
		return errors.New("no filter specified")
	}

	filter := req.GetFilter()
	startDate, endDate := filter.GetStartDate(), filter.GetEndDate()

	if filter.HasStartDate() && filter.HasEndDate() && startDate.AsTime().After(endDate.AsTime()) {
		return errors.New("start date cannot be after end date")
	}

	return nil
}
