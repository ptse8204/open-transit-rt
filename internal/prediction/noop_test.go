package prediction

import (
	"context"
	"testing"
)

func TestNoopAdapterReturnsExplicitDiagnostics(t *testing.T) {
	adapter := NewNoopAdapter()
	if adapter.Name() != "noop" {
		t.Fatalf("adapter name = %q, want noop", adapter.Name())
	}
	result, err := adapter.PredictTripUpdates(context.Background(), Request{AgencyID: "demo-agency"})
	if err != nil {
		t.Fatalf("predict trip updates: %v", err)
	}
	if len(result.TripUpdates) != 0 {
		t.Fatalf("trip updates = %d, want empty no-op output", len(result.TripUpdates))
	}
	if result.Diagnostics.Status != StatusNoop || result.Diagnostics.Reason != ReasonNoopAdapter {
		t.Fatalf("diagnostics = %+v, want noop/noop_adapter", result.Diagnostics)
	}
}
