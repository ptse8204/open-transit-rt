package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/tenant"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InstanceState string

const (
	StateExampleAvailable    InstanceState = "example_available"
	StateNotConfigured       InstanceState = "not_configured"
	StateConfiguredNotTested InstanceState = "configured_not_tested"
	StateDryRunPassed        InstanceState = "dry_run_passed"
	StateReadyForActivation  InstanceState = "ready_for_activation"
	StateActive              InstanceState = "active"
	StateBlocked             InstanceState = "blocked"
)

var SupportedInstanceStates = []InstanceState{
	StateExampleAvailable,
	StateNotConfigured,
	StateConfiguredNotTested,
	StateDryRunPassed,
	StateReadyForActivation,
	StateActive,
	StateBlocked,
}

type Instance struct {
	ID            int64           `json:"id"`
	AgencyID      string          `json:"agency_id"`
	ConnectorType string          `json:"connector_type"`
	ConnectorKind string          `json:"connector_kind"`
	DisplayName   string          `json:"display_name"`
	State         InstanceState   `json:"state"`
	Owner         string          `json:"owner"`
	ConfigJSON    json.RawMessage `json:"config_json,omitempty"`
	SecretRefs    []string        `json:"secret_refs,omitempty"`
	DryRunStatus  string          `json:"dry_run_status"`
	LastCheckedAt *time.Time      `json:"last_checked_at,omitempty"`
	ActivatedAt   *time.Time      `json:"activated_at,omitempty"`
	DisabledAt    *time.Time      `json:"disabled_at,omitempty"`
	CreatedBy     string          `json:"created_by"`
	UpdatedBy     string          `json:"updated_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type InstanceRepository interface {
	ListInstances(ctx context.Context, agencyID string) ([]Instance, error)
}

type PostgresInstanceStore struct {
	pool *pgxpool.Pool
}

func NewPostgresInstanceStore(pool *pgxpool.Pool) *PostgresInstanceStore {
	return &PostgresInstanceStore{pool: pool}
}

func (s *PostgresInstanceStore) ListInstances(ctx context.Context, agencyID string) ([]Instance, error) {
	if s == nil || s.pool == nil {
		return nil, fmt.Errorf("connector instance store is unavailable")
	}
	agencyID = strings.TrimSpace(agencyID)
	if err := tenant.ValidateAgencyID(agencyID); err != nil {
		return nil, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, agency_id, connector_type, connector_kind, display_name, state, owner,
		       config_json, secret_refs, dry_run_status, last_checked_at, activated_at,
		       disabled_at, created_by, updated_by, created_at, updated_at
		FROM connector_instance
		WHERE agency_id = $1
		ORDER BY connector_type, display_name, id
	`, agencyID)
	if err != nil {
		return nil, fmt.Errorf("list connector instances: %w", err)
	}
	defer rows.Close()
	var instances []Instance
	for rows.Next() {
		var inst Instance
		var state string
		var configRaw []byte
		var refsRaw []byte
		if err := rows.Scan(
			&inst.ID,
			&inst.AgencyID,
			&inst.ConnectorType,
			&inst.ConnectorKind,
			&inst.DisplayName,
			&state,
			&inst.Owner,
			&configRaw,
			&refsRaw,
			&inst.DryRunStatus,
			&inst.LastCheckedAt,
			&inst.ActivatedAt,
			&inst.DisabledAt,
			&inst.CreatedBy,
			&inst.UpdatedBy,
			&inst.CreatedAt,
			&inst.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan connector instance: %w", err)
		}
		normalized, err := ParseInstanceState(state)
		if err != nil {
			return nil, err
		}
		inst.State = normalized
		inst.ConfigJSON = safeJSONObject(configRaw)
		if len(refsRaw) > 0 {
			if err := json.Unmarshal(refsRaw, &inst.SecretRefs); err != nil {
				return nil, fmt.Errorf("decode connector secret refs: %w", err)
			}
		}
		inst.SecretRefs = cleanStringList(inst.SecretRefs)
		instances = append(instances, inst)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector instances: %w", err)
	}
	return instances, nil
}

func ParseInstanceState(value string) (InstanceState, error) {
	normalized := InstanceState(strings.TrimSpace(value))
	for _, allowed := range SupportedInstanceStates {
		if normalized == allowed {
			return normalized, nil
		}
	}
	return "", fmt.Errorf("unsupported connector instance state %q", value)
}

func (i Instance) DeploymentConfigExists() bool {
	return len(i.ConfigKeys()) > 0 || len(i.SecretRefs) > 0
}

func (i Instance) ConfigKeys() []string {
	if len(i.ConfigJSON) == 0 {
		return nil
	}
	var obj map[string]any
	if err := json.Unmarshal(i.ConfigJSON, &obj); err != nil {
		return nil
	}
	keys := make([]string, 0, len(obj))
	for key := range obj {
		key = strings.TrimSpace(key)
		if key != "" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}

func safeJSONObject(raw []byte) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage(`{}`)
	}
	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return json.RawMessage(`{}`)
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return json.RawMessage(out)
}

func cleanStringList(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}
