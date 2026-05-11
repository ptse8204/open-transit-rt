package main

import (
	"strings"
	"time"

	"open-transit-rt/internal/compliance"
)

type operationsGTFSQualityGuidanceView struct {
	GeneratedAt time.Time                      `json:"generated_at"`
	AgencyID    string                         `json:"agency_id"`
	Boundary    string                         `json:"boundary"`
	Workflow    []operationsGTFSQualityStep    `json:"workflow"`
	ClaimFlags  operationsGTFSQualityClaimFlag `json:"claim_flags"`
}

type operationsGTFSQualityStep struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Summary     string   `json:"summary"`
	AdminLinks  []string `json:"admin_links"`
	DocsLinks   []string `json:"docs_links"`
	DoesNotDo   string   `json:"does_not_do"`
	NextOutcome string   `json:"next_outcome"`
}

type operationsGTFSQualityClaimFlag struct {
	AutomaticGTFSEditEnabled        bool `json:"automatic_gtfs_edit_enabled"`
	DraftMutationEnabled            bool `json:"draft_mutation_enabled"`
	SchedulePublishEnabled          bool `json:"schedule_publish_enabled"`
	ValidatorSemanticsChanged       bool `json:"validator_semantics_changed"`
	ExternalEvidenceCreated         bool `json:"external_evidence_created"`
	ConsumerStatusesChanged         bool `json:"consumer_statuses_changed"`
	ComplianceClaimed               bool `json:"compliance_claimed"`
	AgencyApprovalClaimed           bool `json:"agency_approval_claimed"`
	ConsumerAcceptanceClaimed       bool `json:"consumer_acceptance_claimed"`
	PublicLaunchClaimed             bool `json:"public_launch_claimed"`
	ProductionReadinessClaimed      bool `json:"production_readiness_claimed"`
	ProductionGradeETAClaimed       bool `json:"production_grade_eta_claimed"`
	VendorCompatibilityClaimed      bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed    bool `json:"hardware_certification_claimed"`
	ProductionAVLReliabilityClaimed bool `json:"production_avl_reliability_claimed"`
}

func buildOperationsGTFSQualityGuidance(page operationsPage) operationsGTFSQualityGuidanceView {
	return operationsGTFSQualityGuidanceView{
		GeneratedAt: page.GeneratedAt,
		AgencyID:    page.AgencyID,
		Boundary:    "Private GTFS quality guidance only. It explains likely fixes for operator review, but it does not edit GTFS, mutate drafts, publish schedules, change validator semantics, create evidence, change consumer statuses, or claim compliance or approval.",
		Workflow: []operationsGTFSQualityStep{
			{
				ID:          "identify-owner",
				Label:       "Identify owner",
				Summary:     "Use the issue family, source, and sample to decide whether the fix belongs to the schedule planner, GTFS source export, GTFS Studio draft reviewer, GIS/shape maintainer, or technical maintainer.",
				AdminLinks:  []string{"/admin/operations/gtfs-quality", "/admin/operations/gtfs-import"},
				DocsLinks:   []string{"docs/tutorials/gtfs-validation-triage.md", "docs/tutorials/real-agency-gtfs-onboarding.md"},
				DoesNotDo:   "Does not infer agency approval or assign work to an external party.",
				NextOutcome: "A named operator or maintainer has enough context to fix the source data outside this page.",
			},
			{
				ID:          "fix-source",
				Label:       "Fix source data",
				Summary:     "Apply corrections in the source GTFS export, GTFS Studio draft, or source scheduling system, then re-export or republish through the existing import/publish flow.",
				AdminLinks:  []string{"/admin/operations/gtfs-import", "/admin/gtfs-studio"},
				DocsLinks:   []string{"docs/tutorials/real-agency-gtfs-onboarding.md", "docs/requirements-2a-2f.md"},
				DoesNotDo:   "Does not auto-edit GTFS rows, mutate drafts, publish schedules, or bypass review.",
				NextOutcome: "A reviewed ZIP or draft can go back through the normal import, publish, and validation paths.",
			},
			{
				ID:          "verify",
				Label:       "Verify",
				Summary:     "Re-import or publish the reviewed schedule, rerun the allowlisted static validator when appropriate, and check feed health/readiness before relying on the feed.",
				AdminLinks:  []string{"/admin/operations/validation-health", "/admin/operations/feed-health", "/admin/operations/readiness"},
				DocsLinks:   []string{"docs/tutorials/calitp-readiness-checklist.md", "docs/tutorials/operator-smoke-and-support-bundle.md"},
				DoesNotDo:   "Does not turn validator output into compliance, consumer acceptance, or public-launch proof.",
				NextOutcome: "Operators have a fresh supporting signal and a clear next action for any remaining blocker.",
			},
		},
		ClaimFlags: operationsGTFSQualityClaimFlag{},
	}
}

func gtfsQualityLikelyOwner(group compliance.GTFSQualityGroup) string {
	switch group.Family {
	case "expired_calendar", "calendar_service_dates":
		return "Schedule planner or GTFS source owner"
	case "route_short_name_too_long":
		return "Route naming owner"
	case "unused_shape", "shape_ordering":
		return "GIS or shapes maintainer"
	case "bad_stop_times", "frequency_issues", "block_transition_issues":
		return "Schedule planner with operations review"
	case "missing_or_foreign_key_reference", "duplicate_ids":
		return "GTFS export owner or source-system admin"
	default:
		if group.Source == compliance.GTFSQualitySourceInternalImporter {
			return "Technical maintainer with GTFS source owner"
		}
		return "GTFS source owner with technical maintainer review"
	}
}

func gtfsQualityAffectedFiles(group compliance.GTFSQualityGroup) string {
	switch group.Family {
	case "expired_calendar", "calendar_service_dates":
		return "calendar.txt / calendar_dates.txt"
	case "route_short_name_too_long":
		return "routes.txt"
	case "unused_shape":
		return "shapes.txt / trips.txt"
	case "shape_ordering":
		return "shapes.txt"
	case "bad_stop_times":
		return "stop_times.txt / trips.txt"
	case "frequency_issues":
		return "frequencies.txt / trips.txt"
	case "block_transition_issues":
		return "trips.txt / stop_times.txt"
	case "missing_or_foreign_key_reference":
		return "reported referring file and referenced GTFS file"
	case "duplicate_ids":
		return "reported ID column and related references"
	default:
		return "reported file from the validator sample"
	}
}

func gtfsQualitySafeFixPath(source string, group compliance.GTFSQualityGroup) string {
	if source == compliance.GTFSQualitySourceInternalImporter || group.Source == compliance.GTFSQualitySourceInternalImporter {
		return "Fix the source GTFS or GTFS Studio draft, then re-import or publish through the existing flow. Do not patch database rows or bypass importer validation."
	}
	switch group.Family {
	case "unused_shape":
		return "Remove or reconnect shapes only after confirming no active trip should use them, then rerun the static validator."
	case "route_short_name_too_long":
		return "Move descriptive text to route_long_name or agency source notes when appropriate, then export and rerun validation."
	case "bad_stop_times":
		return "Correct times, stop sequence order, and after-midnight formatting in the source schedule, then rerun validation and matching smoke checks."
	case "frequency_issues":
		return "Review headway windows and trip references in the source schedule, then validate frequency-based service before relying on realtime matching."
	case "block_transition_issues":
		return "Confirm block_id continuity and trip order with operations staff, then validate block transition behavior."
	default:
		return "Correct the source GTFS or GTFS Studio draft, re-export or publish, then rerun the allowlisted static validator."
	}
}

func gtfsQualityVerifyWith(source string, group compliance.GTFSQualityGroup) string {
	if source == compliance.GTFSQualitySourceInternalImporter || group.Source == compliance.GTFSQualitySourceInternalImporter {
		return "Re-import through browser or CLI, then review this internal importer section again."
	}
	switch group.Family {
	case "bad_stop_times", "frequency_issues", "block_transition_issues":
		return "Rerun static validation, then review `/admin/operations/telemetry` and feed health after synthetic telemetry checks."
	case "expired_calendar", "calendar_service_dates":
		return "Rerun static validation and confirm the active service day appears in feed health/readiness."
	default:
		return "Rerun the allowlisted static validator, then review GTFS quality, validator health, and feed health."
	}
}

func gtfsQualityEscalation(group compliance.GTFSQualityGroup) string {
	switch group.Severity {
	case compliance.GTFSQualityBlocking:
		return "Escalate before relying on the schedule feed or realtime outputs."
	case compliance.GTFSQualityNeedsReview:
		return "Review with the data owner before treating the warning as acceptable."
	case compliance.GTFSQualityInformational:
		return "Track for the next data-quality review unless it masks a known service issue."
	default:
		return "Classify the recurring notice with a maintainer before closing it."
	}
}

func gtfsQualityGuidanceClass(group compliance.GTFSQualityGroup) string {
	return strings.ReplaceAll(strings.TrimSpace(group.Family), "_", "-")
}
