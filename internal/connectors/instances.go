package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"open-transit-rt/internal/tenant"

	"github.com/jackc/pgx/v5"
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

type DryRunJob struct {
	ID                  int64           `json:"id"`
	AgencyID            string          `json:"agency_id"`
	ConnectorInstanceID int64           `json:"connector_instance_id"`
	Status              string          `json:"status"`
	StartedAt           time.Time       `json:"started_at"`
	FinishedAt          time.Time       `json:"finished_at"`
	RedactedSummary     json.RawMessage `json:"redacted_summary,omitempty"`
	AcceptedCount       int             `json:"accepted_count"`
	RejectedCount       int             `json:"rejected_count"`
	DroppedCount        int             `json:"dropped_count"`
	RedactionScanStatus string          `json:"redaction_scan_status"`
	CreatedBy           string          `json:"created_by"`
	CreatedAt           time.Time       `json:"created_at"`
}

type UpsertInstanceInput struct {
	AgencyID      string
	ConnectorType string
	ConnectorKind string
	DisplayName   string
	State         InstanceState
	Owner         string
	ConfigJSON    json.RawMessage
	SecretRefs    []string
	DryRunStatus  string
	ActorID       string
	Now           time.Time
}

type CreateDryRunJobInput struct {
	AgencyID            string
	ConnectorInstanceID int64
	Status              string
	RedactedSummary     json.RawMessage
	AcceptedCount       int
	RejectedCount       int
	DroppedCount        int
	RedactionScanStatus string
	ActorID             string
	Now                 time.Time
}

type UpdateInstanceStateInput struct {
	AgencyID string
	ID       int64
	State    InstanceState
	ActorID  string
	Now      time.Time
}

type InstanceRepository interface {
	ListInstances(ctx context.Context, agencyID string) ([]Instance, error)
	ListDryRunJobs(ctx context.Context, agencyID string, limit int) ([]DryRunJob, error)
}

type InstanceWriter interface {
	UpsertInstance(ctx context.Context, input UpsertInstanceInput) (Instance, error)
	CreateDryRunJob(ctx context.Context, input CreateDryRunJobInput) (DryRunJob, error)
	UpdateInstanceState(ctx context.Context, input UpdateInstanceStateInput) (Instance, error)
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
		inst, err := scanConnectorInstance(rows)
		if err != nil {
			return nil, fmt.Errorf("scan connector instance: %w", err)
		}
		instances = append(instances, inst)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector instances: %w", err)
	}
	return instances, nil
}

func (s *PostgresInstanceStore) ListDryRunJobs(ctx context.Context, agencyID string, limit int) ([]DryRunJob, error) {
	if s == nil || s.pool == nil {
		return nil, fmt.Errorf("connector instance store is unavailable")
	}
	agencyID = strings.TrimSpace(agencyID)
	if err := tenant.ValidateAgencyID(agencyID); err != nil {
		return nil, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx, `
		SELECT id, agency_id, connector_instance_id, status, started_at, finished_at,
		       redacted_summary, accepted_count, rejected_count, dropped_count,
		       redaction_scan_status, created_by, created_at
		FROM connector_dry_run_job
		WHERE agency_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2
	`, agencyID, limit)
	if err != nil {
		return nil, fmt.Errorf("list connector dry-run jobs: %w", err)
	}
	defer rows.Close()
	var jobs []DryRunJob
	for rows.Next() {
		job, err := scanDryRunJob(rows)
		if err != nil {
			return nil, fmt.Errorf("scan connector dry-run job: %w", err)
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connector dry-run jobs: %w", err)
	}
	return jobs, nil
}

func (s *PostgresInstanceStore) UpsertInstance(ctx context.Context, input UpsertInstanceInput) (Instance, error) {
	if s == nil || s.pool == nil {
		return Instance{}, fmt.Errorf("connector instance store is unavailable")
	}
	input, err := normalizeUpsertInstanceInput(input)
	if err != nil {
		return Instance{}, err
	}
	refsJSON, err := json.Marshal(input.SecretRefs)
	if err != nil {
		return Instance{}, fmt.Errorf("encode connector secret refs: %w", err)
	}
	row := s.pool.QueryRow(ctx, `
		INSERT INTO connector_instance (
		  agency_id, connector_type, connector_kind, display_name, state, owner,
		  config_json, secret_refs, dry_run_status, created_by, updated_by,
		  created_at, updated_at, activated_at, disabled_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8::jsonb, $9, $10, $10, $11, $11, NULL, NULL)
		ON CONFLICT (agency_id, connector_type, connector_kind, display_name) DO UPDATE
		SET state = EXCLUDED.state,
		    owner = EXCLUDED.owner,
		    config_json = EXCLUDED.config_json,
		    secret_refs = EXCLUDED.secret_refs,
		    dry_run_status = EXCLUDED.dry_run_status,
		    updated_by = EXCLUDED.updated_by,
		    updated_at = EXCLUDED.updated_at,
		    activated_at = NULL,
		    disabled_at = NULL
		RETURNING id, agency_id, connector_type, connector_kind, display_name, state, owner,
		          config_json, secret_refs, dry_run_status, last_checked_at, activated_at,
		          disabled_at, created_by, updated_by, created_at, updated_at
	`, input.AgencyID, input.ConnectorType, input.ConnectorKind, input.DisplayName, string(input.State), input.Owner, []byte(input.ConfigJSON), refsJSON, input.DryRunStatus, input.ActorID, input.Now)
	inst, err := scanConnectorInstance(row)
	if err != nil {
		return Instance{}, fmt.Errorf("upsert connector instance: %w", err)
	}
	return inst, nil
}

func (s *PostgresInstanceStore) CreateDryRunJob(ctx context.Context, input CreateDryRunJobInput) (DryRunJob, error) {
	if s == nil || s.pool == nil {
		return DryRunJob{}, fmt.Errorf("connector instance store is unavailable")
	}
	input, err := normalizeCreateDryRunJobInput(input)
	if err != nil {
		return DryRunJob{}, err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return DryRunJob{}, fmt.Errorf("begin connector dry-run transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	var agencyID string
	if err := tx.QueryRow(ctx, `
		SELECT agency_id
		FROM connector_instance
		WHERE id = $1 AND agency_id = $2
	`, input.ConnectorInstanceID, input.AgencyID).Scan(&agencyID); err != nil {
		return DryRunJob{}, fmt.Errorf("connector instance not found: %w", err)
	}
	nextState := StateConfiguredNotTested
	if input.Status == "passed" && input.RedactionScanStatus == "passed" {
		nextState = StateDryRunPassed
	}
	if input.Status == "failed" || input.Status == "blocked" || input.RedactionScanStatus != "passed" {
		nextState = StateBlocked
	}
	row := tx.QueryRow(ctx, `
		INSERT INTO connector_dry_run_job (
		  agency_id, connector_instance_id, status, started_at, finished_at,
		  redacted_summary, accepted_count, rejected_count, dropped_count,
		  redaction_scan_status, created_by, created_at
		)
		VALUES ($1, $2, $3, $4, $4, $5::jsonb, $6, $7, $8, $9, $10, $4)
		RETURNING id, agency_id, connector_instance_id, status, started_at, finished_at,
		          redacted_summary, accepted_count, rejected_count, dropped_count,
		          redaction_scan_status, created_by, created_at
	`, input.AgencyID, input.ConnectorInstanceID, input.Status, input.Now, []byte(input.RedactedSummary), input.AcceptedCount, input.RejectedCount, input.DroppedCount, input.RedactionScanStatus, input.ActorID)
	job, err := scanDryRunJob(row)
	if err != nil {
		return DryRunJob{}, fmt.Errorf("insert connector dry-run job: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		UPDATE connector_instance
		SET dry_run_status = $1,
		    state = $2,
		    last_checked_at = $3,
		    updated_by = $4,
		    updated_at = $3
		WHERE agency_id = $5 AND id = $6
	`, input.Status, string(nextState), input.Now, input.ActorID, input.AgencyID, input.ConnectorInstanceID); err != nil {
		return DryRunJob{}, fmt.Errorf("update connector dry-run status: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return DryRunJob{}, fmt.Errorf("commit connector dry-run transaction: %w", err)
	}
	return job, nil
}

func (s *PostgresInstanceStore) UpdateInstanceState(ctx context.Context, input UpdateInstanceStateInput) (Instance, error) {
	if s == nil || s.pool == nil {
		return Instance{}, fmt.Errorf("connector instance store is unavailable")
	}
	input, err := normalizeUpdateInstanceStateInput(input)
	if err != nil {
		return Instance{}, err
	}
	row := s.pool.QueryRow(ctx, `
		UPDATE connector_instance
		SET state = $1,
		    activated_at = CASE WHEN $1 = 'active' THEN $2 ELSE activated_at END,
		    disabled_at = CASE WHEN $1 = 'active' THEN NULL ELSE disabled_at END,
		    updated_by = $3,
		    updated_at = $2
		WHERE agency_id = $4 AND id = $5
		RETURNING id, agency_id, connector_type, connector_kind, display_name, state, owner,
		          config_json, secret_refs, dry_run_status, last_checked_at, activated_at,
		          disabled_at, created_by, updated_by, created_at, updated_at
	`, string(input.State), input.Now, input.ActorID, input.AgencyID, input.ID)
	inst, err := scanConnectorInstance(row)
	if err != nil {
		return Instance{}, fmt.Errorf("update connector instance state: %w", err)
	}
	return inst, nil
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

func normalizeUpdateInstanceStateInput(input UpdateInstanceStateInput) (UpdateInstanceStateInput, error) {
	input.AgencyID = strings.TrimSpace(input.AgencyID)
	if err := tenant.ValidateAgencyID(input.AgencyID); err != nil {
		return input, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	if input.ID <= 0 {
		return input, fmt.Errorf("connector instance id is required")
	}
	if _, err := ParseInstanceState(string(input.State)); err != nil {
		return input, err
	}
	input.ActorID = strings.TrimSpace(input.ActorID)
	if input.ActorID == "" {
		input.ActorID = "operator_console"
	}
	input.Now = input.Now.UTC()
	if input.Now.IsZero() {
		input.Now = time.Now().UTC()
	}
	return input, nil
}

func normalizeCreateDryRunJobInput(input CreateDryRunJobInput) (CreateDryRunJobInput, error) {
	input.AgencyID = strings.TrimSpace(input.AgencyID)
	if err := tenant.ValidateAgencyID(input.AgencyID); err != nil {
		return input, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	if input.ConnectorInstanceID <= 0 {
		return input, fmt.Errorf("connector instance id is required")
	}
	switch strings.TrimSpace(input.Status) {
	case "passed", "failed", "blocked":
		input.Status = strings.TrimSpace(input.Status)
	default:
		return input, fmt.Errorf("unsupported dry-run status")
	}
	switch strings.TrimSpace(input.RedactionScanStatus) {
	case "passed", "failed", "blocked":
		input.RedactionScanStatus = strings.TrimSpace(input.RedactionScanStatus)
	default:
		return input, fmt.Errorf("unsupported redaction scan status")
	}
	if input.AcceptedCount < 0 || input.RejectedCount < 0 || input.DroppedCount < 0 {
		return input, fmt.Errorf("dry-run counts must be non-negative")
	}
	if input.RedactedSummary == nil {
		input.RedactedSummary = json.RawMessage(`{}`)
	}
	if err := validateInstanceConfig(input.RedactedSummary); err != nil {
		return input, err
	}
	input.RedactedSummary = safeJSONObject(input.RedactedSummary)
	input.ActorID = strings.TrimSpace(input.ActorID)
	if input.ActorID == "" {
		input.ActorID = "operator_console"
	}
	input.Now = input.Now.UTC()
	if input.Now.IsZero() {
		input.Now = time.Now().UTC()
	}
	return input, nil
}

func normalizeUpsertInstanceInput(input UpsertInstanceInput) (UpsertInstanceInput, error) {
	input.AgencyID = strings.TrimSpace(input.AgencyID)
	if err := tenant.ValidateAgencyID(input.AgencyID); err != nil {
		return input, fmt.Errorf("agency_id must be path-safe: %w", err)
	}
	if !supportedConnectorType(input.ConnectorType) {
		return input, fmt.Errorf("unsupported connector type %q", input.ConnectorType)
	}
	input.ConnectorKind = strings.TrimSpace(input.ConnectorKind)
	if err := validateSafeConnectorLabel("connector kind", input.ConnectorKind); err != nil {
		return input, err
	}
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	if err := validateSafeConnectorLabel("display name", input.DisplayName); err != nil {
		return input, err
	}
	input.Owner = strings.TrimSpace(input.Owner)
	if len(input.Owner) > 160 {
		return input, fmt.Errorf("owner is too long")
	}
	if input.State == "" {
		input.State = StateConfiguredNotTested
	}
	if _, err := ParseInstanceState(string(input.State)); err != nil {
		return input, err
	}
	if input.ConfigJSON == nil {
		input.ConfigJSON = json.RawMessage(`{}`)
	}
	if err := validateInstanceConfig(input.ConfigJSON); err != nil {
		return input, err
	}
	input.ConfigJSON = safeJSONObject(input.ConfigJSON)
	input.SecretRefs = cleanStringList(input.SecretRefs)
	for _, ref := range input.SecretRefs {
		if err := ValidateSecretRef(ref); err != nil {
			return input, err
		}
	}
	input.DryRunStatus = strings.TrimSpace(input.DryRunStatus)
	if input.DryRunStatus == "" {
		input.DryRunStatus = "not_run"
	}
	input.ActorID = strings.TrimSpace(input.ActorID)
	if input.ActorID == "" {
		input.ActorID = "operator_console"
	}
	input.Now = input.Now.UTC()
	if input.Now.IsZero() {
		input.Now = time.Now().UTC()
	}
	return input, nil
}

func supportedConnectorType(value string) bool {
	for _, supported := range SupportedTypes {
		if value == supported {
			return true
		}
	}
	return false
}

func validateSafeConnectorLabel(field string, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	if len(value) > 160 {
		return fmt.Errorf("%s is too long", field)
	}
	if strings.ContainsAny(value, "\x00\r\n\t") {
		return fmt.Errorf("%s contains unsupported control characters", field)
	}
	if strings.Contains(value, "://") {
		return fmt.Errorf("%s must not contain endpoint URLs", field)
	}
	return nil
}

var secretRefPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,79}$`)

func ValidateSecretRef(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if !secretRefPattern.MatchString(value) {
		return fmt.Errorf("secret refs must be uppercase deployment reference labels")
	}
	return nil
}

func validateInstanceConfig(raw json.RawMessage) error {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return fmt.Errorf("config metadata must be valid JSON")
	}
	obj, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("config metadata must be a JSON object")
	}
	return rejectInstanceConfigSecrets("config_json", obj)
}

func rejectInstanceConfigSecrets(path string, value any) error {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			lower := strings.ToLower(key)
			if strings.Contains(lower, "password") || strings.Contains(lower, "token") || strings.Contains(lower, "secret") || strings.Contains(lower, "api_key") {
				return fmt.Errorf("%s.%s must use secret_refs instead of secret-like config keys", path, key)
			}
			if err := rejectInstanceConfigSecrets(path+"."+key, child); err != nil {
				return err
			}
		}
	case []any:
		for _, child := range typed {
			if err := rejectInstanceConfigSecrets(path+"[]", child); err != nil {
				return err
			}
		}
	case string:
		lower := strings.ToLower(typed)
		if strings.Contains(lower, "://") {
			return fmt.Errorf("%s must not contain endpoint URLs", path)
		}
		if strings.Contains(lower, "password=") || strings.Contains(lower, "token=") || strings.Contains(lower, "secret=") {
			return fmt.Errorf("%s must not contain inline secrets", path)
		}
	}
	return nil
}

type instanceScanner interface {
	Scan(dest ...any) error
}

func scanConnectorInstance(row instanceScanner) (Instance, error) {
	var inst Instance
	var state string
	var configRaw []byte
	var refsRaw []byte
	if err := row.Scan(
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
		if err == pgx.ErrNoRows {
			return Instance{}, err
		}
		return Instance{}, err
	}
	normalized, err := ParseInstanceState(state)
	if err != nil {
		return Instance{}, err
	}
	inst.State = normalized
	inst.ConfigJSON = safeJSONObject(configRaw)
	if len(refsRaw) > 0 {
		if err := json.Unmarshal(refsRaw, &inst.SecretRefs); err != nil {
			return Instance{}, fmt.Errorf("decode connector secret refs: %w", err)
		}
	}
	inst.SecretRefs = cleanStringList(inst.SecretRefs)
	return inst, nil
}

func scanDryRunJob(row instanceScanner) (DryRunJob, error) {
	var job DryRunJob
	var summaryRaw []byte
	if err := row.Scan(
		&job.ID,
		&job.AgencyID,
		&job.ConnectorInstanceID,
		&job.Status,
		&job.StartedAt,
		&job.FinishedAt,
		&summaryRaw,
		&job.AcceptedCount,
		&job.RejectedCount,
		&job.DroppedCount,
		&job.RedactionScanStatus,
		&job.CreatedBy,
		&job.CreatedAt,
	); err != nil {
		return DryRunJob{}, err
	}
	job.RedactedSummary = safeJSONObject(summaryRaw)
	return job, nil
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
