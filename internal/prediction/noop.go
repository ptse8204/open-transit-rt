package prediction

import "context"

type NoopAdapter struct{}

func NewNoopAdapter() NoopAdapter {
	return NoopAdapter{}
}

func (NoopAdapter) Name() string {
	return "noop"
}

func (NoopAdapter) PredictTripUpdates(context.Context, Request) (Result, error) {
	return Result{
		TripUpdates: nil,
		Diagnostics: Diagnostics{
			Status: StatusNoop,
			Reason: ReasonNoopAdapter,
		},
	}, nil
}
