# Hosted Validator Records

- Environment: `hosted-pending`
- Capture date (UTC): 2026-04-22
- Operator: pending
- Status: missing

## Required Evidence

Clean validator records are required for:

- Static GTFS `schedule.zip`.
- GTFS-RT Vehicle Positions.
- GTFS-RT Trip Updates.
- GTFS-RT Alerts.

## Collection Commands

Run from the deployment host or a trusted operator workstation with admin access:

```sh
mkdir -p "$ENVIRONMENT_NAME/validation"

for feed_type in schedule vehicle_positions trip_updates alerts
do
  if [ "$feed_type" = "schedule" ]; then
    validator_id="static-mobilitydata"
  else
    validator_id="realtime-mobilitydata"
  fi

  curl -sS -X POST "$ADMIN_BASE_URL/admin/validation/run" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    --data "{\"validator_id\":\"$validator_id\",\"feed_type\":\"$feed_type\"}" \
    | tee "$ENVIRONMENT_NAME/validation/validate-$feed_type.json"
done
```

## Required Summary To Fill

| Feed type | Validator ID/version | Run timestamp UTC | Status | Errors | Warnings | Full output |
| --- | --- | --- | --- | ---: | ---: | --- |
| schedule | pending | pending | pending | pending | pending | pending |
| vehicle_positions | pending | pending | pending | pending | pending | pending |
| trip_updates | pending | pending | pending | pending | pending | pending |
| alerts | pending | pending | pending | pending | pending | pending |

## Blocker

No hosted validator outputs were available when this intake packet was created. Local Step 2 validation outputs are retained separately and all failed; they do not satisfy this requirement.
