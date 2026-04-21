package alerts

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type ReconcileResult struct {
	CreatedOrUpdated int `json:"created_or_updated"`
	LinkedReviews    int `json:"linked_reviews"`
}

func (r *PostgresRepository) ReconcileCanceledTripAlerts(ctx context.Context, agencyID string, actorID string, at time.Time) (ReconcileResult, error) {
	if at.IsZero() {
		at = time.Now().UTC()
	}
	if actorID == "" {
		actorID = "system"
	}
	feedVersionID, err := r.activeFeedVersionID(ctx, agencyID)
	if err != nil {
		return ReconcileResult{}, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, vehicle_id, route_id, trip_id, start_date, start_time, expires_at
		FROM manual_override
		WHERE agency_id = $1
		  AND override_type = 'canceled_trip'
		  AND state = 'canceled'
		  AND cleared_at IS NULL
		  AND (expires_at IS NULL OR expires_at > $2)
		  AND trip_id IS NOT NULL
		  AND start_date IS NOT NULL
		  AND start_time IS NOT NULL
		ORDER BY trip_id, start_date, start_time, id
	`, agencyID, at)
	if err != nil {
		return ReconcileResult{}, fmt.Errorf("query canceled trip overrides: %w", err)
	}
	defer rows.Close()

	result := ReconcileResult{}
	for rows.Next() {
		var overrideID int64
		var vehicleID, tripID, startDate, startTime string
		var routeID sql.NullString
		var expiresAt sql.NullTime
		if err := rows.Scan(&overrideID, &vehicleID, &routeID, &tripID, &startDate, &startTime, &expiresAt); err != nil {
			return ReconcileResult{}, fmt.Errorf("scan canceled trip override: %w", err)
		}
		alertKey := fmt.Sprintf("canceled:%s:%s:%s", tripID, startDate, startTime)
		header := fmt.Sprintf("Trip %s canceled", tripID)
		description := "This trip has been canceled by operations."
		alert, err := r.UpsertAlert(ctx, UpsertInput{
			AgencyID:        agencyID,
			AlertKey:        alertKey,
			Cause:           "other_cause",
			Effect:          "no_service",
			HeaderText:      header,
			DescriptionText: description,
			ActiveStart:     &at,
			ActiveEnd:       timePtrFromNull(expiresAt),
			FeedVersionID:   feedVersionID,
			SourceType:      SourceCancellationReconciler,
			SourceID:        strconv.FormatInt(overrideID, 10),
			Metadata: map[string]any{
				"prediction_override_id": overrideID,
				"vehicle_id":             vehicleID,
			},
			ActorID: actorID,
			Publish: true,
			Now:     at,
			Entities: []InformedEntity{{
				AgencyID:  agencyID,
				RouteID:   routeID.String,
				TripID:    tripID,
				StartDate: startDate,
				StartTime: startTime,
			}},
		})
		if err != nil {
			return ReconcileResult{}, err
		}
		result.CreatedOrUpdated++
		linked, err := r.linkCancellationReviews(ctx, agencyID, vehicleID, tripID, startDate, startTime, alert.ID, actorID, at)
		if err != nil {
			return ReconcileResult{}, err
		}
		result.LinkedReviews += linked
	}
	if err := rows.Err(); err != nil {
		return ReconcileResult{}, fmt.Errorf("iterate canceled trip overrides: %w", err)
	}
	return result, nil
}

func (r *PostgresRepository) activeFeedVersionID(ctx context.Context, agencyID string) (string, error) {
	var feedVersionID string
	err := r.pool.QueryRow(ctx, `
		SELECT id
		FROM feed_version
		WHERE agency_id = $1 AND is_active
		ORDER BY activated_at DESC NULLS LAST, created_at DESC
		LIMIT 1
	`, agencyID).Scan(&feedVersionID)
	if err != nil {
		return "", fmt.Errorf("query active feed version for alerts: %w", err)
	}
	return feedVersionID, nil
}

func (r *PostgresRepository) linkCancellationReviews(ctx context.Context, agencyID string, vehicleID string, tripID string, startDate string, startTime string, alertID int64, actorID string, at time.Time) (int, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, details_json
		FROM incident
		WHERE agency_id = $1
		  AND incident_type = 'prediction_review'
		  AND status = 'open'
		  AND vehicle_id = $2
		  AND trip_id = $3
		  AND details_json->>'expected_alert_missing' = 'true'
		  AND COALESCE(details_json->>'start_date', '') = $4
		  AND COALESCE(details_json->>'start_time', '') = $5
	`, agencyID, vehicleID, tripID, startDate, startTime)
	if err != nil {
		return 0, fmt.Errorf("query cancellation review incidents: %w", err)
	}
	defer rows.Close()
	type review struct {
		id      int64
		details map[string]any
	}
	var reviews []review
	for rows.Next() {
		var item review
		var detailsBytes []byte
		if err := rows.Scan(&item.id, &detailsBytes); err != nil {
			return 0, fmt.Errorf("scan cancellation review incident: %w", err)
		}
		item.details = map[string]any{}
		if len(detailsBytes) > 0 {
			_ = json.Unmarshal(detailsBytes, &item.details)
		}
		reviews = append(reviews, item)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterate cancellation review incidents: %w", err)
	}
	for _, item := range reviews {
		item.details["expected_alert_missing"] = false
		item.details["service_alert_id"] = alertID
		item.details["cancellation_alert_linkage_status"] = "linked"
		payload, err := json.Marshal(item.details)
		if err != nil {
			return 0, fmt.Errorf("marshal linked review details: %w", err)
		}
		if _, err := r.pool.Exec(ctx, `
			UPDATE incident
			SET status = 'resolved',
			    resolved_at = $3,
			    details_json = $4::jsonb
			WHERE agency_id = $1 AND id = $2 AND incident_type = 'prediction_review'
		`, agencyID, item.id, at, string(payload)); err != nil {
			return 0, fmt.Errorf("update linked cancellation review: %w", err)
		}
	}
	return len(reviews), nil
}

func timePtrFromNull(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}
