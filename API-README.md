# API (gRPC)

The gRPC server listens on `localhost:50051` by default.

Service: `core.Service` (see `core/core.proto`)

## Update

Path: `core.Service/Update`

Request: `UpdateRequest`
- `Event` (model.Event)

Response: `UpdateResponse`
- `Message` (string)

## GetSportEvent

Path: `core.Service/GetSportEvent`

Request: `GetSportEventRequest`
- `EventID` (string)

Response: `GetSportEventResponse`
- `Event` (SportEvent)

## GetRaceEvent

Path: `core.Service/GetRaceEvent`

Request: `GetRaceEventRequest`
- `EventID` (string)

Response: `GetRaceEventResponse`
- `Event` (RaceEvent)

## SearchEvents

Path: `core.Service/SearchEvents`

Request: `SearchEventsRequest`
- `Filter` (SearchEventsFilter)
  - `BettingStatus` (optional model.BettingStatus)
  - `EventVisibility` (optional model.EventVisibility)
  - `StartDate` (optional google.protobuf.Timestamp)
  - `EndDate` (optional google.protobuf.Timestamp)
- `PageSize` (optional uint64)
- `PageToken` (optional string)

Response: `SearchEventsResponse`
- `SportEvents` (repeated SportEvent)
- `RaceEvents` (repeated RaceEvent)
- `NextPageToken` (string)

For authoritative field definitions, see `core/core.proto` and `model/event.proto`.
You can also use `test.http` to run example requests directly from GoLand.
