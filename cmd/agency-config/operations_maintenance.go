package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type operationsMaintenanceView struct {
	GeneratedAt     time.Time                        `json:"generated_at"`
	AgencyID        string                           `json:"agency_id"`
	Boundary        string                           `json:"boundary"`
	OverallStatus   string                           `json:"overall_status"`
	SummaryRows     []operationsMaintenanceRow       `json:"summary_rows"`
	Diagnostics     operationsMaintenanceDiagnostics `json:"diagnostics"`
	BackupRestore   operationsMaintenancePanel       `json:"backup_restore"`
	UpgradeRollback operationsMaintenancePanel       `json:"upgrade_rollback"`
	SupportReview   operationsMaintenancePanel       `json:"support_review"`
	CadencePlan     operationsMaintenancePanel       `json:"cadence_plan"`
	Infrastructure  operationsMaintenancePanel       `json:"infrastructure"`
	Tasks           []operationsMaintenanceTask      `json:"tasks"`
	SupportSummary  operationsMaintenanceSupport     `json:"support_summary"`
	ClaimFlags      operationsMaintenanceClaimFlags  `json:"claim_flags"`
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

type operationsMaintenanceDiagnostics struct {
	Status   string                            `json:"status"`
	Boundary string                            `json:"boundary"`
	Rows     []operationsMaintenanceDiagnostic `json:"rows"`
}

type operationsMaintenanceDiagnostic struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	SourceRef     string `json:"source_ref"`
	GeneratedAt   string `json:"generated_at"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsMaintenancePanel struct {
	Status     string                          `json:"status"`
	Boundary   string                          `json:"boundary"`
	NextAction string                          `json:"next_action"`
	Rows       []operationsMaintenancePanelRow `json:"rows"`
}

type operationsMaintenancePanelRow struct {
	ID                  string `json:"id"`
	Label               string `json:"label"`
	Status              string `json:"status"`
	CurrentSignal       string `json:"current_signal"`
	OperatorStep        string `json:"operator_step"`
	TechnicalHelperStep string `json:"technical_helper_step"`
	DoesNotProve        string `json:"does_not_prove"`
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
		maintenanceTask("weekly_feed_health", "weekly", "Check five configured feed paths and feed health next actions.", cockpitFeedHealthStatus(page), "agency operator", "Open Feed Health and review each row."),
		maintenanceTask("weekly_validators", "weekly", "Review validator health and stale reports.", cockpitValidationStatus(page), "agency operator or technical helper", "Open Validator Health; use off-host validation when the server is too small."),
		maintenanceTask("weekly_telemetry", "weekly", "Review telemetry freshness, stale rows, and assignment confidence.", cockpitTelemetryStatus(page), "agency operator", "Open Devices and Telemetry."),
		maintenanceTask("weekly_alerts", "weekly", "Check whether active alerts need to be created, updated, or archived.", alertStatus(page), "agency operator", "Open the Alerts Console and Alerts feed row."),
		maintenanceTask("monthly_gtfs_update", "monthly", "Import updated GTFS and review counts, quality, validators, and feed health.", cockpitActiveFeedStatus(page), "schedule owner with agency operator", "Open Browser GTFS Import."),
		maintenanceTask("monthly_backup_restore", "monthly", "Confirm backup and restore-drill configuration presence.", maintenanceBackupTaskStatus(), "technical helper", "Configure private backup/restore values and keep secret values out of docs."),
		maintenanceTask("as_needed_support_summary", "as needed", "Generate a local support bundle only when a technical helper needs diagnostics.", operationsStatusDiagnosticOnly, "technical helper", "Run the support-bundle helper from an operator shell and redact before sharing."),
	}
	overall := maintenanceOverall(rows, tasks)
	return operationsMaintenanceView{
		GeneratedAt:     page.GeneratedAt,
		AgencyID:        page.AgencyID,
		Boundary:        "Private maintenance diagnostics only. This page summarizes configured/nonconfigured signals and next tasks without creating evidence, changing consumer statuses, claiming compliance, claiming production readiness, claiming SLA or uptime coverage, or exposing secret values.",
		OverallStatus:   overall,
		SummaryRows:     rows,
		Diagnostics:     buildOperationsMaintenanceDiagnostics(),
		BackupRestore:   buildOperationsMaintenanceBackupRestore(),
		UpgradeRollback: buildOperationsMaintenanceUpgradeRollback(),
		SupportReview:   buildOperationsMaintenanceSupportReview(),
		CadencePlan:     buildOperationsMaintenanceCadencePlan(),
		Infrastructure:  buildOperationsMaintenanceInfrastructureChecks(),
		Tasks:           tasks,
		SupportSummary: operationsMaintenanceSupport{
			Status:     operationsStatusDiagnosticOnly,
			Command:    "make support-bundle",
			OutputPath: ".cache/support-bundles/<timestamp>",
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

func maintenancePanelRow(id, label, status, signal, operatorStep, helperStep, doesNotProve string) operationsMaintenancePanelRow {
	return operationsMaintenancePanelRow{ID: id, Label: label, Status: status, CurrentSignal: signal, OperatorStep: operatorStep, TechnicalHelperStep: helperStep, DoesNotProve: doesNotProve}
}

func buildOperationsMaintenanceBackupRestore() operationsMaintenancePanel {
	rows := []operationsMaintenancePanelRow{
		maintenancePanelRow(
			"backup_configuration_presence",
			"Backup configuration presence",
			maintenanceEnvStatus("BACKUP_DIR", "BACKUP_PATH", "OPEN_TRANSIT_BACKUP_DIR"),
			maintenanceEnvPresenceSignal("BACKUP_DIR", "BACKUP_PATH", "OPEN_TRANSIT_BACKUP_DIR"),
			"Keep backup configuration marked missing until a deployment owner supplies a private backup target.",
			"Configure backup output outside the browser and keep backup contents out of docs/evidence unless a separate evidence phase is authorized.",
			"Configuration presence does not prove a successful backup exists.",
		),
		maintenancePanelRow(
			"restore_drill_configuration_presence",
			"Restore-drill configuration presence",
			maintenanceEnvStatus("RESTORE_DRILL_DATABASE_URL", "RESTORE_DRILL_TARGET", "OPEN_TRANSIT_RESTORE_DRILL"),
			maintenanceEnvPresenceSignal("RESTORE_DRILL_DATABASE_URL", "RESTORE_DRILL_TARGET", "OPEN_TRANSIT_RESTORE_DRILL"),
			"Keep restore-drill readiness marked missing until a private restore target is configured.",
			"Run restore drills only from an operator shell against an explicit non-live target; never paste restore URLs into the browser.",
			"Configuration presence does not prove a restore succeeded or that disaster recovery coverage exists.",
		),
		maintenancePanelRow(
			"deployment_doctor_backup_restore",
			"Deployment doctor backup/restore summary",
			operationsStatusDiagnosticOnly,
			"deployment-doctor records backup_readiness and restore_readiness summary fields when a local helper has run",
			"Use the Local Diagnostic Summaries table to see the latest safe deployment-doctor status.",
			"Run `make deployment-doctor` from an operator shell when a fresh private summary is needed.",
			"Deployment-doctor status does not execute backup or restore actions.",
		),
		maintenancePanelRow(
			"browser_destructive_actions",
			"Browser destructive actions",
			operationsStatusReady,
			"disabled: no backup, restore, rollback, migration, or package action is exposed by this page",
			"Use the browser only for review and next-step guidance.",
			"Use documented shell workflows with explicit operator confirmation for destructive work.",
			"Disabled browser actions do not prove the shell workflows have been run.",
		),
	}
	return operationsMaintenancePanel{
		Status:     maintenancePanelOverall(rows),
		Boundary:   "Backup and restore review is private guidance only. This page shows configuration presence and safe summary pointers; it does not create backups, restore databases, run migrations, read backup dumps, or create retained evidence.",
		NextAction: "Resolve missing configuration first, then run private backup or restore-drill workflows from an operator shell when authorized.",
		Rows:       rows,
	}
}

func buildOperationsMaintenanceUpgradeRollback() operationsMaintenancePanel {
	rows := []operationsMaintenancePanelRow{
		maintenancePanelRow(
			"upgrade_precheck",
			"Upgrade precheck",
			operationsStatusNeedsReview,
			"review release notes, current commit, database migration status, active feed version, and backup configuration before upgrade",
			"Use this checklist before a local/source upgrade; do not treat the browser as an upgrade executor.",
			"Run `make check`, `make validate`, and deployment-specific prechecks from an operator shell before changing a deployment.",
			"Precheck guidance does not prove release readiness or production readiness.",
		),
		maintenancePanelRow(
			"rollback_precheck",
			"Rollback precheck",
			operationsStatusNeedsReview,
			"confirm the rollback target, active feed version, migration direction, and restore owner before rollback",
			"Keep rollback marked review-required until a deployment owner confirms the target and data implications.",
			"Use documented rollback procedures outside the browser and avoid destructive database changes without a tested restore path.",
			"A rollback checklist does not prove rollback success.",
		),
		maintenancePanelRow(
			"migration_safety",
			"Migration safety",
			operationsStatusNeedsReview,
			"migrations are not run from the browser and must remain backward-compatible or explicitly reviewed",
			"Review migration status and backup/restore readiness before any schema change.",
			"Run `go run ./cmd/migrate status` or the deployment doctor from an operator shell when needed.",
			"Migration status does not prove data-loss safety by itself.",
		),
		maintenancePanelRow(
			"release_artifact_boundary",
			"Release artifact boundary",
			operationsStatusReady,
			"no tag, package, image, push, or release publish action is exposed by this page",
			"Use this page only to understand maintenance readiness.",
			"Keep release-cut cleanup separate until a maintainer explicitly authorizes release actions.",
			"Absence of browser release actions does not prove a release candidate is ready.",
		),
	}
	return operationsMaintenancePanel{
		Status:     maintenancePanelOverall(rows),
		Boundary:   "Upgrade and rollback review is checklist-only. The browser does not tag, package, publish, run migrations, roll back services, or restore databases.",
		NextAction: "Use the checklist to decide whether a technical helper must run shell-based upgrade, rollback, migration-status, or restore-readiness checks.",
		Rows:       rows,
	}
}

func buildOperationsMaintenanceSupportReview() operationsMaintenancePanel {
	rows := []operationsMaintenancePanelRow{
		maintenancePanelRow(
			"support_bundle_output_scope",
			"Support bundle output scope",
			operationsStatusDiagnosticOnly,
			"support bundles are private local diagnostics under .cache/support-bundles/<timestamp>",
			"Use support bundles only when a technical helper needs a bounded private diagnostic snapshot.",
			"Run `make support-bundle` from an operator shell; the browser does not generate, upload, or inspect raw bundle files.",
			"Support bundle creation does not prove evidence, compliance, production readiness, or external acceptance.",
		),
		maintenancePanelRow(
			"redaction_review",
			"Redaction review",
			operationsStatusNeedsReview,
			"review manifest, excluded categories, generated file list, and private-value warnings before sharing",
			"Treat every generated support bundle as private until a human has reviewed it.",
			"Use the redaction policy and rerun the bundle if a private value appears in a shareable summary.",
			"A helper review does not prove the output is safe for public release.",
		),
		maintenancePanelRow(
			"evidence_boundary",
			"Evidence boundary",
			operationsStatusReady,
			"support bundles stay outside docs/evidence unless a separate retained-evidence phase is authorized",
			"Do not rename support bundles as evidence or attach them to consumer packet records by default.",
			"Follow the evidence redaction policy only after separate written maintainer approval for retained evidence.",
			"Support diagnostics do not prove final-root readiness, deployment proof, compliance, or consumer acceptance.",
		),
		maintenancePanelRow(
			"private_output_warning",
			"Private output warning",
			operationsStatusNeedsReview,
			"do not expose raw logs, database URLs, access tokens, cookies, private payloads, backup dumps, or unredacted files",
			"Share only bounded facts required for the support question.",
			"Remove private values and convert raw output into public-safe summaries before moving anything outside the operator environment.",
			"A warning row does not prove local files contain no private data.",
		),
	}
	return operationsMaintenancePanel{
		Status:     maintenancePanelOverall(rows),
		Boundary:   "Support-bundle guidance is private and review-only. The browser does not generate bundles, inspect raw bundle files, upload diagnostics, or create retained evidence.",
		NextAction: "Run support bundles from an operator shell only when needed, then perform a human redaction review before sharing any output.",
		Rows:       rows,
	}
}

func buildOperationsMaintenanceCadencePlan() operationsMaintenancePanel {
	rows := []operationsMaintenancePanelRow{
		maintenancePanelRow(
			"daily_operating_check",
			"Daily operating check",
			operationsStatusDiagnosticOnly,
			"review feed health, telemetry freshness, alerts, and local diagnostic summary status on service days",
			"Use the private Operations Console to decide whether an issue needs same-day attention.",
			"Escalate stale telemetry, blocked feed checks, or alert lifecycle gaps to the technical helper.",
			"Daily checks do not prove uptime, SLA coverage, or public consumer display.",
		),
		maintenancePanelRow(
			"weekly_maintenance_check",
			"Weekly maintenance check",
			operationsStatusDiagnosticOnly,
			"review validators, five-feed checks, backup and restore configuration, notification drafts, and reliability summaries",
			"Confirm that routine health rows have a named owner before service changes pile up.",
			"Run validator-health, deployment-doctor, operations-notify, and operations-reliability helpers from an operator shell when a fresh private summary is needed.",
			"Weekly review does not prove validator-clean feeds, compliance, or release readiness.",
		),
		maintenancePanelRow(
			"monthly_recovery_check",
			"Monthly recovery check",
			operationsStatusNeedsReview,
			"review GTFS update cadence, backup target, restore-drill target, rollback owner, and upgrade notes",
			"Keep recovery readiness marked review-required until a deployment owner confirms the non-live restore target.",
			"Run restore drills only against explicit non-live targets and keep secret connection values out of browser and docs.",
			"Recovery planning does not prove restore success or disaster recovery coverage.",
		),
		maintenancePanelRow(
			"as_needed_support_check",
			"As-needed support check",
			operationsStatusDiagnosticOnly,
			"generate support bundles only for a specific support question or blocker",
			"Capture the support question first so the bundle can be reviewed and summarized narrowly.",
			"Run `make support-bundle`, review redaction, then share only bounded facts needed for the question.",
			"As-needed support diagnostics do not prove compliance, external acceptance, or public launch readiness.",
		),
	}
	return operationsMaintenancePanel{
		Status:     maintenancePanelOverall(rows),
		Boundary:   "Maintenance cadence rows are planning guidance only. They do not schedule jobs, send notifications, run validators, execute backups, or create proof.",
		NextAction: "Use cadence rows to assign owners and decide when a technical helper should run private shell diagnostics.",
		Rows:       rows,
	}
}

func maintenancePanelOverall(rows []operationsMaintenancePanelRow) string {
	status := operationsStatusReady
	for _, row := range rows {
		status = worseMaintenanceStatus(status, row.Status)
	}
	return status
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
