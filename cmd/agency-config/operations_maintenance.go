package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type operationsMaintenanceView struct {
	GeneratedAt    time.Time                       `json:"generated_at"`
	AgencyID       string                          `json:"agency_id"`
	Boundary       string                          `json:"boundary"`
	OverallStatus  string                          `json:"overall_status"`
	SummaryRows    []operationsMaintenanceRow      `json:"summary_rows"`
	Tasks          []operationsMaintenanceTask     `json:"tasks"`
	SupportSummary operationsMaintenanceSupport    `json:"support_summary"`
	ClaimFlags     operationsMaintenanceClaimFlags `json:"claim_flags"`
}

type operationsMaintenanceRow struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsMaintenanceTask struct {
	ID       string `json:"id"`
	Cadence  string `json:"cadence"`
	Task     string `json:"task"`
	Status   string `json:"status"`
	Owner    string `json:"owner"`
	NextStep string `json:"next_step"`
}

type operationsMaintenanceSupport struct {
	Status       string   `json:"status"`
	Command      string   `json:"command"`
	OutputPath   string   `json:"output_path"`
	Instructions []string `json:"instructions"`
}

type operationsMaintenanceClaimFlags struct {
	ExternalEvidenceCreated    bool `json:"external_evidence_created"`
	ConsumerStatusesChanged    bool `json:"consumer_statuses_changed"`
	ComplianceClaimed          bool `json:"compliance_claimed"`
	ProductionReadinessClaimed bool `json:"production_readiness_claimed"`
	SLAClaimed                 bool `json:"sla_claimed"`
	UptimeGuaranteeClaimed     bool `json:"uptime_guarantee_claimed"`
	HostedSaaSClaimed          bool `json:"hosted_saas_claimed"`
	AgencyAdoptionClaimed      bool `json:"agency_adoption_claimed"`
	ConsumerAcceptanceClaimed  bool `json:"consumer_acceptance_claimed"`
	VendorCompatibilityClaimed bool `json:"vendor_compatibility_claimed"`
	ProductionGradeETAClaimed  bool `json:"production_grade_eta_claimed"`
}

func buildOperationsMaintenance(page operationsPage) operationsMaintenanceView {
	rows := []operationsMaintenanceRow{
		maintenanceRow("deployed_version", "Deployed commit/version", maintenanceEnvStatus("OPEN_TRANSIT_RT_COMMIT", "GIT_COMMIT", "APP_VERSION"), maintenanceVersionSignal(), "Set a non-secret commit or version environment value during deployment so support can identify the running build.", "Does not prove a release was tagged or published."),
		maintenanceRow("active_feed_version", "Active feed version", maintenancePresentStatus(page.ActiveFeedVersion), firstNonEmpty(page.ActiveFeedVersion, "not available"), "Import or activate a schedule before reviewing realtime or validator state.", "Does not prove rollback execution is available in the browser."),
		maintenanceRow("last_gtfs_import", "Last GTFS import / feed metadata", maintenanceTimeStatus(page.FeedsUpdatedAt), formatTimeForText(page.FeedsUpdatedAt), "Use GTFS import and quality pages to review the latest schedule state.", "Does not prove the imported schedule is validator clean."),
		maintenanceRow("last_five_feed_check", "Last five-feed check", maintenanceLatestReliabilityStatus(page), maintenanceLatestReliabilitySignal(page), "Run local/reference feed checks or review reliability snapshots; missing data remains not configured.", "Does not prove public final-root readiness or consumer availability."),
		maintenanceRow("validator_state", "Validator state", maintenanceValidatorStatus(page), fmt.Sprintf("overall=%s; tooling=%s", page.ValidationHealth.OverallStatus, page.ValidationHealth.ToolingStatus), "Install or review pinned/off-host validators and rerun private validator health.", "Does not prove compliance or consumer acceptance."),
		maintenanceRow("backup_configuration", "Backup configuration", maintenanceEnvStatus("BACKUP_DIR", "BACKUP_PATH", "OPEN_TRANSIT_BACKUP_DIR"), maintenanceEnvPresenceSignal("BACKUP_DIR", "BACKUP_PATH", "OPEN_TRANSIT_BACKUP_DIR"), "Configure private backup output and document the restore owner before relying on backups.", "Does not prove a backup exists or a restore succeeded."),
		maintenanceRow("restore_drill_configuration", "Restore-drill configuration", maintenanceEnvStatus("RESTORE_DRILL_DATABASE_URL", "RESTORE_DRILL_TARGET", "OPEN_TRANSIT_RESTORE_DRILL"), maintenanceEnvPresenceSignal("RESTORE_DRILL_DATABASE_URL", "RESTORE_DRILL_TARGET", "OPEN_TRANSIT_RESTORE_DRILL"), "Configure restore-drill settings without printing secret values, then run the private drill workflow.", "Does not prove restore readiness or disaster recovery coverage."),
		maintenanceRow("telemetry_freshness", "Telemetry freshness", cockpitTelemetryStatus(page), telemetryEvidence(page), "Review devices and simulator guidance if Vehicle Positions are empty or stale.", "Does not prove real device reliability."),
		maintenanceRow("service_health", "Service health", operationsStatusUnknown, "not available from this single request unless an operator runs deployment diagnostics", "Run local/reference diagnostics or loopback health checks through the deployment helper when needed.", "Does not prove uptime, SLA, or hosted service availability."),
	}
	tasks := []operationsMaintenanceTask{
		maintenanceTask("weekly_feed_health", "weekly", "Check five public feed paths and feed health next actions.", cockpitFeedHealthStatus(page), "agency operator", "Open Feed Health and review each row."),
		maintenanceTask("weekly_validators", "weekly", "Review validator health and stale reports.", cockpitValidationStatus(page), "agency operator or technical helper", "Open Validator Health; use off-host validation when the server is too small."),
		maintenanceTask("weekly_telemetry", "weekly", "Review telemetry freshness, stale rows, and assignment confidence.", cockpitTelemetryStatus(page), "agency operator", "Open Devices and Telemetry."),
		maintenanceTask("weekly_alerts", "weekly", "Check whether active alerts need to be created, updated, or archived.", alertStatus(page), "agency operator", "Open the Alerts Console and Alerts feed row."),
		maintenanceTask("monthly_gtfs_update", "monthly", "Import updated GTFS and review counts, quality, validators, and feed health.", cockpitActiveFeedStatus(page), "schedule owner with agency operator", "Open Browser GTFS Import."),
		maintenanceTask("monthly_backup_restore", "monthly", "Confirm backup and restore-drill configuration presence.", maintenanceBackupTaskStatus(), "technical helper", "Configure private backup/restore values and keep secret values out of docs."),
		maintenanceTask("as_needed_support_summary", "as needed", "Generate a local support bundle only when a technical helper needs diagnostics.", operationsStatusDiagnosticOnly, "technical helper", "Run the support-bundle helper from an operator shell and redact before sharing."),
	}
	overall := maintenanceOverall(rows, tasks)
	return operationsMaintenanceView{
		GeneratedAt:   page.GeneratedAt,
		AgencyID:      page.AgencyID,
		Boundary:      "Private maintenance diagnostics only. This page summarizes configured/nonconfigured signals and next tasks without creating evidence, changing consumer statuses, claiming compliance, claiming production readiness, claiming SLA or uptime coverage, or exposing secret values.",
		OverallStatus: overall,
		SummaryRows:   rows,
		Tasks:         tasks,
		SupportSummary: operationsMaintenanceSupport{
			Status:     operationsStatusDiagnosticOnly,
			Command:    "make support-bundle",
			OutputPath: ".cache/support-bundle/<timestamp>",
			Instructions: []string{
				"Run only from an operator-controlled shell when a technical helper needs diagnostics.",
				"Review and redact output before sharing outside the operator environment.",
				"Do not place support bundles under docs/evidence unless a separate retained-evidence approval exists.",
			},
		},
		ClaimFlags: operationsMaintenanceClaimFlags{},
	}
}

func maintenanceRow(id, label, status, signal, next, doesNotProve string) operationsMaintenanceRow {
	return operationsMaintenanceRow{ID: id, Label: label, Status: status, CurrentSignal: signal, NextAction: next, DoesNotProve: doesNotProve}
}

func maintenanceTask(id, cadence, task, status, owner, next string) operationsMaintenanceTask {
	return operationsMaintenanceTask{ID: id, Cadence: cadence, Task: task, Status: status, Owner: owner, NextStep: next}
}

func maintenanceOverall(rows []operationsMaintenanceRow, tasks []operationsMaintenanceTask) string {
	status := operationsStatusReady
	for _, row := range rows {
		status = worseMaintenanceStatus(status, row.Status)
	}
	for _, task := range tasks {
		status = worseMaintenanceStatus(status, task.Status)
	}
	return status
}

func worseMaintenanceStatus(current, next string) string {
	rank := map[string]int{
		operationsStatusReady:          0,
		operationsStatusDiagnosticOnly: 1,
		operationsStatusNeedsReview:    2,
		operationsStatusUnknown:        3,
		operationsStatusMissing:        4,
		operationsStatusBlocked:        5,
	}
	if rank[next] > rank[current] {
		return next
	}
	return current
}

func maintenancePresentStatus(value string) string {
	if strings.TrimSpace(value) == "" {
		return operationsStatusMissing
	}
	return operationsStatusReady
}

func maintenanceTimeStatus(t *time.Time) string {
	if t == nil || t.IsZero() {
		return operationsStatusMissing
	}
	return operationsStatusReady
}

func maintenanceEnvStatus(names ...string) string {
	for _, name := range names {
		if strings.TrimSpace(os.Getenv(name)) != "" {
			return operationsStatusReady
		}
	}
	return operationsStatusMissing
}

func maintenanceEnvPresenceSignal(names ...string) string {
	var present []string
	var missing []string
	for _, name := range names {
		if strings.TrimSpace(os.Getenv(name)) != "" {
			present = append(present, name+"=configured")
		} else {
			missing = append(missing, name+"=not configured")
		}
	}
	if len(present) > 0 {
		return strings.Join(present, "; ") + "; values withheld"
	}
	return strings.Join(missing, "; ")
}

func maintenanceVersionSignal() string {
	for _, name := range []string{"OPEN_TRANSIT_RT_COMMIT", "GIT_COMMIT", "APP_VERSION"} {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return name + " is configured; value withheld from this summary unless deployment docs explicitly allow it"
		}
	}
	return "not available"
}

func maintenanceLatestReliabilityStatus(page operationsPage) string {
	if len(page.Reliability.Feeds) == 0 {
		return operationsStatusMissing
	}
	for _, row := range page.Reliability.Feeds {
		switch row.Status {
		case "unhealthy":
			return operationsStatusBlocked
		case "missing":
			return operationsStatusMissing
		case "needs_review", "unknown":
			return operationsStatusNeedsReview
		}
	}
	return operationsStatusReady
}

func maintenanceLatestReliabilitySignal(page operationsPage) string {
	var latest *time.Time
	for _, row := range page.Reliability.Feeds {
		if row.SnapshotAt != nil && (latest == nil || row.SnapshotAt.After(*latest)) {
			t := row.SnapshotAt.UTC()
			latest = &t
		}
	}
	if latest == nil {
		return "not available"
	}
	return "latest private feed-health snapshot " + formatTimeForText(latest)
}

func maintenanceValidatorStatus(page operationsPage) string {
	return cockpitValidationStatus(page)
}

func maintenanceBackupTaskStatus() string {
	if maintenanceEnvStatus("BACKUP_DIR", "BACKUP_PATH", "OPEN_TRANSIT_BACKUP_DIR") == operationsStatusReady &&
		maintenanceEnvStatus("RESTORE_DRILL_DATABASE_URL", "RESTORE_DRILL_TARGET", "OPEN_TRANSIT_RESTORE_DRILL") == operationsStatusReady {
		return operationsStatusReady
	}
	return operationsStatusMissing
}
