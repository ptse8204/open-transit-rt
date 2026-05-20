package main

import (
	"fmt"
	"strings"
	"time"

	"open-transit-rt/internal/compliance"
)

const gtfsQualityFixPlannerMaxRows = 50

type operationsGTFSQualityGuidanceView struct {
	GeneratedAt time.Time                       `json:"generated_at"`
	AgencyID    string                          `json:"agency_id"`
	Boundary    string                          `json:"boundary"`
	Workflow    []operationsGTFSQualityStep     `json:"workflow"`
	FixPlanner  operationsGTFSQualityFixPlanner `json:"fix_planner"`
	ClaimFlags  operationsGTFSQualityClaimFlag  `json:"claim_flags"`
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
	DraftSuggestionRecordsCreated   bool `json:"draft_suggestion_records_created"`
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

type operationsGTFSQualityFixPlanner struct {
	Status              string                        `json:"status"`
	Summary             string                        `json:"summary"`
	Boundary            string                        `json:"boundary"`
	DraftSuggestionMode string                        `json:"draft_suggestion_mode"`
	BeforeValidation    []string                      `json:"before_validation"`
	AfterValidation     []string                      `json:"after_validation"`
	Rows                []operationsGTFSQualityFixRow `json:"rows"`
	DisplayedRows       int                           `json:"displayed_rows"`
	TotalRows           int                           `json:"total_rows"`
	HiddenRows          int                           `json:"hidden_rows"`
	Checklist           string                        `json:"checklist"`
}

type operationsGTFSQualityFixRow struct {
	ID                    string   `json:"id"`
	Source                string   `json:"source"`
	SourceLabel           string   `json:"source_label"`
	Severity              string   `json:"severity"`
	Family                string   `json:"family"`
	Codes                 []string `json:"codes"`
	Count                 int      `json:"count"`
	LikelyOwner           string   `json:"likely_owner"`
	RiskLevel             string   `json:"risk_level"`
	AffectedFiles         string   `json:"affected_files"`
	IssueSummary          string   `json:"issue_summary"`
	WhyItMatters          string   `json:"why_it_matters"`
	SafeFixSuggestion     string   `json:"safe_fix_suggestion"`
	DraftSuggestion       string   `json:"draft_suggestion"`
	DraftSuggestionRecord string   `json:"draft_suggestion_record"`
	BeforeValidationPlan  string   `json:"before_validation_plan"`
	AfterValidationPlan   string   `json:"after_validation_plan"`
	VerifyWith            string   `json:"verify_with"`
	EscalateIf            string   `json:"escalate_if"`
	Samples               []string `json:"samples"`
	NoAutoApplyBoundary   string   `json:"no_auto_apply_boundary"`
}

func buildOperationsGTFSQualityGuidance(page operationsPage) operationsGTFSQualityGuidanceView {
	view := operationsGTFSQualityGuidanceView{
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
	view.FixPlanner = buildOperationsGTFSQualityFixPlanner(page)
	return view
}

func buildOperationsGTFSQualityFixPlanner(page operationsPage) operationsGTFSQualityFixPlanner {
	planner := operationsGTFSQualityFixPlanner{
		Status:              gtfsQualityFixPlannerStatus(page.GTFSQuality),
		Summary:             "Private fix planner generated from sanitized GTFS validator and importer groups.",
		Boundary:            "Advisory private checklist only. It does not automatically edit production GTFS, write GTFS Studio draft rows, publish schedules, create evidence, contact external parties, change consumer statuses, or prove compliance.",
		DraftSuggestionMode: "manual_review_only",
		BeforeValidation: []string{
			"Confirm the active published feed version and whether the issue came from the canonical static validator or internal importer.",
			"Assign an operator, schedule owner, GIS/shapes maintainer, source-system admin, or technical maintainer before changing source data.",
			"Choose either source GTFS export correction or GTFS Studio draft authoring; do not patch production database rows.",
		},
		AfterValidation: []string{
			"Re-import the corrected GTFS ZIP or publish the reviewed GTFS Studio draft through the existing private flow.",
			"Rerun the allowlisted static validator when available, then review GTFS Quality, Validation Health, Feed Health, and Schedule Review.",
			"Keep remaining warnings in needs_review until a data owner documents why no source change is required.",
		},
	}
	planner.Rows = append(planner.Rows, gtfsQualityFixRowsFromSection("canonical_static", page.GTFSQuality.Canonical)...)
	planner.Rows = append(planner.Rows, gtfsQualityFixRowsFromSection("internal_importer", page.GTFSQuality.InternalImporter)...)
	planner.TotalRows = len(planner.Rows)
	if planner.TotalRows > gtfsQualityFixPlannerMaxRows {
		planner.HiddenRows = planner.TotalRows - gtfsQualityFixPlannerMaxRows
		planner.Rows = planner.Rows[:gtfsQualityFixPlannerMaxRows]
	}
	planner.DisplayedRows = len(planner.Rows)
	planner.Checklist = gtfsQualityFixChecklist(page, planner)
	if planner.TotalRows == 0 {
		planner.Summary = "No grouped GTFS quality issue rows are available. Import or publish GTFS and run validators before using this planner as an operator checklist."
	}
	return planner
}

func gtfsQualityFixPlannerStatus(triage compliance.GTFSQualityTriage) string {
	return worstGTFSQualityStatus(triage.Canonical.Status, triage.InternalImporter.Status)
}

func worstGTFSQualityStatus(statuses ...string) string {
	worst := compliance.GTFSQualityUnknown
	for _, status := range statuses {
		switch status {
		case compliance.GTFSQualityBlocking:
			return compliance.GTFSQualityBlocking
		case compliance.GTFSQualityNeedsReview:
			if worst != compliance.GTFSQualityBlocking {
				worst = compliance.GTFSQualityNeedsReview
			}
		case compliance.GTFSQualityInformational:
			if worst == "" || worst == compliance.GTFSQualityUnknown {
				worst = compliance.GTFSQualityInformational
			}
		case compliance.GTFSQualityUnknown, "":
			if worst == "" {
				worst = compliance.GTFSQualityUnknown
			}
		default:
			if worst == "" {
				worst = compliance.GTFSQualityUnknown
			}
		}
	}
	if worst == "" {
		return compliance.GTFSQualityUnknown
	}
	return worst
}

func gtfsQualityFixRowsFromSection(prefix string, section compliance.GTFSQualitySection) []operationsGTFSQualityFixRow {
	rows := make([]operationsGTFSQualityFixRow, 0, len(section.Groups))
	for index, group := range section.Groups {
		rows = append(rows, gtfsQualityFixRowFromGroup(prefix, section, group, index))
	}
	return rows
}

func gtfsQualityFixRowFromGroup(prefix string, section compliance.GTFSQualitySection, group compliance.GTFSQualityGroup, index int) operationsGTFSQualityFixRow {
	source := firstNonEmpty(group.Source, section.Source)
	return operationsGTFSQualityFixRow{
		ID:                    fmt.Sprintf("%s_%03d_%s", prefix, index+1, strings.ReplaceAll(group.Family, "_", "-")),
		Source:                source,
		SourceLabel:           section.SourceLabel,
		Severity:              group.Severity,
		Family:                group.Family,
		Codes:                 append([]string(nil), group.Codes...),
		Count:                 group.Count,
		LikelyOwner:           gtfsQualityLikelyOwner(group),
		RiskLevel:             gtfsQualityRiskLevel(group),
		AffectedFiles:         gtfsQualityAffectedFiles(group),
		IssueSummary:          group.OperatorSummary,
		WhyItMatters:          group.WhyItMatters,
		SafeFixSuggestion:     gtfsQualitySafeFixPath(section.Source, group),
		DraftSuggestion:       gtfsQualityDraftSuggestion(source, group),
		DraftSuggestionRecord: "Advisory only; no persisted draft suggestion record is created by this Operations Console route.",
		BeforeValidationPlan:  gtfsQualityBeforeValidationPlan(section, group),
		AfterValidationPlan:   gtfsQualityAfterValidationPlan(section, group),
		VerifyWith:            gtfsQualityVerifyWith(section.Source, group),
		EscalateIf:            gtfsQualityEscalation(group),
		Samples:               append([]string(nil), group.Samples...),
		NoAutoApplyBoundary:   "No automatic production edit. Do not auto-apply this suggestion to published GTFS, draft tables, feed versions, or validator records.",
	}
}

func gtfsQualityDraftSuggestion(source string, group compliance.GTFSQualityGroup) string {
	if source == compliance.GTFSQualitySourceInternalImporter || group.Source == compliance.GTFSQualitySourceInternalImporter {
		return "Use GTFS Studio only for reviewed manual draft authoring when the source export cannot be corrected first; otherwise fix the source GTFS and re-import."
	}
	switch group.Family {
	case "route_short_name_too_long":
		return "Manual draft candidate: review routes.txt naming fields with the route naming owner before changing route_short_name or route_long_name."
	case "route_metadata":
		return "Manual draft candidate: review routes.txt names, colors, route type, and agency links with the route owner before editing route metadata."
	case "stop_location":
		return "Manual draft candidate: review stops.txt coordinates and station hierarchy with the GIS or stop inventory owner before editing stops."
	case "agency_metadata":
		return "Manual draft candidate: review agency.txt metadata with the administrator before changing feed identity fields."
	case "license_contact_metadata":
		return "Manual draft candidate: review feed_info.txt, license, attribution, and contact metadata with the deployment owner before changing sharing prep fields."
	case "missing_required_file":
		return "Manual draft candidate only if GTFS Studio can author the missing file safely; otherwise regenerate the source export."
	case "expired_calendar", "calendar_service_dates":
		return "Manual draft candidate: review calendar.txt and calendar_dates.txt service coverage with the schedule owner before changing service dates."
	case "bad_stop_times":
		return "Manual draft candidate: review stop_times.txt ordering and after-midnight times with the schedule owner before editing trip timing."
	case "frequency_issues":
		return "Manual draft candidate: review frequencies.txt windows and trip references before changing headway-based service."
	case "unused_shape", "shape_ordering":
		return "Manual draft candidate: review shapes.txt and related trips with the GIS/shapes maintainer before removing or reordering geometry."
	case "block_transition_issues":
		return "Manual draft candidate: review trips.txt block_id continuity with operations staff before changing block assignments."
	default:
		return "Manual draft candidate only after a data owner confirms the source-system correction path is not the right first step."
	}
}

func gtfsQualityBeforeValidationPlan(section compliance.GTFSQualitySection, group compliance.GTFSQualityGroup) string {
	return fmt.Sprintf("Before validation: confirm source=%s, owner=%s, affected files=%s, and whether the fix belongs in the source export or a reviewed GTFS Studio draft.", firstNonEmpty(group.Source, section.Source), gtfsQualityLikelyOwner(group), gtfsQualityAffectedFiles(group))
}

func gtfsQualityAfterValidationPlan(section compliance.GTFSQualitySection, group compliance.GTFSQualityGroup) string {
	return fmt.Sprintf("After validation: %s Keep this row open until the new validator/importer result is reviewed against the active feed version.", gtfsQualityVerifyWith(section.Source, group))
}

func gtfsQualityFixChecklist(page operationsPage, planner operationsGTFSQualityFixPlanner) string {
	var b strings.Builder
	b.WriteString("Private GTFS quality fix checklist\n")
	b.WriteString("Agency: " + page.AgencyID + "\n")
	if page.ActiveFeedVersion != "" {
		b.WriteString("Active schedule feed version: " + page.ActiveFeedVersion + "\n")
	} else {
		b.WriteString("Active schedule feed version: missing\n")
	}
	b.WriteString("Status: " + firstNonEmpty(planner.Status, compliance.GTFSQualityUnknown) + "\n")
	b.WriteString("Boundary: advisory only; no automatic production edit, no draft mutation, no schedule publish, no evidence write, no consumer status change, no compliance claim.\n\n")
	if len(planner.Rows) == 0 {
		b.WriteString("No grouped issue rows are available. Import or publish GTFS and run validators before using this as a fix checklist.\n")
		return b.String()
	}
	for index, row := range planner.Rows {
		b.WriteString(fmt.Sprintf("%d. [%s] %s / %s (%d notice(s))\n", index+1, row.Severity, row.SourceLabel, row.Family, row.Count))
		b.WriteString("   Owner: " + row.LikelyOwner + "\n")
		b.WriteString("   Risk: " + row.RiskLevel + "\n")
		b.WriteString("   Files: " + row.AffectedFiles + "\n")
		b.WriteString("   Safe fix: " + row.SafeFixSuggestion + "\n")
		b.WriteString("   Safe draft suggestion: " + row.DraftSuggestion + "\n")
		b.WriteString("   Before validation plan: " + row.BeforeValidationPlan + "\n")
		b.WriteString("   After validation plan: " + row.AfterValidationPlan + "\n")
		b.WriteString("   Boundary: " + row.NoAutoApplyBoundary + "\n")
	}
	if planner.HiddenRows > 0 {
		b.WriteString(fmt.Sprintf("\n%d additional grouped issue row(s) are hidden by the private planner display cap.\n", planner.HiddenRows))
	}
	return b.String()
}

func gtfsQualityLikelyOwner(group compliance.GTFSQualityGroup) string {
	switch group.Family {
	case "expired_calendar", "calendar_service_dates":
		return "Schedule planner or GTFS source owner"
	case "route_short_name_too_long":
		return "Route naming owner"
	case "route_metadata":
		return "Route naming owner"
	case "stop_location":
		return "GIS or stop inventory owner"
	case "agency_metadata", "license_contact_metadata":
		return "Administrator with GTFS source owner"
	case "unused_shape", "shape_ordering":
		return "GIS or shapes maintainer"
	case "bad_stop_times", "frequency_issues", "block_transition_issues":
		return "Schedule planner with operations review"
	case "missing_required_file", "missing_or_foreign_key_reference", "duplicate_ids":
		return "GTFS export owner or source-system admin"
	default:
		if group.Source == compliance.GTFSQualitySourceInternalImporter {
			return "Technical maintainer with GTFS source owner"
		}
		return "GTFS source owner with technical maintainer review"
	}
}

func gtfsQualityRiskLevel(group compliance.GTFSQualityGroup) string {
	if strings.TrimSpace(group.RiskLevel) != "" {
		return group.RiskLevel
	}
	switch group.Severity {
	case compliance.GTFSQualityBlocking:
		return "blocks import or reliable feed use"
	case compliance.GTFSQualityNeedsReview:
		switch group.Family {
		case "expired_calendar", "calendar_service_dates", "missing_required_file", "missing_or_foreign_key_reference", "bad_stop_times", "frequency_issues", "block_transition_issues":
			return "can break service availability or realtime usefulness"
		case "agency_metadata", "license_contact_metadata":
			return "can block sharing preparation or operator trust"
		case "route_metadata", "stop_location", "shape_ordering", "unused_shape":
			return "can degrade maps, matching, or downstream display"
		default:
			return "needs source-owner review before relying on the feed"
		}
	case compliance.GTFSQualityInformational:
		return "track during normal data-quality review"
	default:
		return "unclassified impact; maintainer review needed"
	}
}

func gtfsQualityAffectedFiles(group compliance.GTFSQualityGroup) string {
	switch group.Family {
	case "missing_required_file":
		return "required GTFS file named in the validator notice"
	case "expired_calendar", "calendar_service_dates":
		return "calendar.txt / calendar_dates.txt"
	case "agency_metadata":
		return "agency.txt"
	case "license_contact_metadata":
		return "feed_info.txt / agency contact metadata"
	case "route_metadata":
		return "routes.txt"
	case "stop_location":
		return "stops.txt"
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
	case "route_metadata":
		return "Correct route type, names, colors, URLs, and agency links in routes.txt source data, then rerun validation."
	case "stop_location":
		return "Correct stop coordinates, location_type, parent station, and naming fields with the GIS or stop inventory owner."
	case "agency_metadata":
		return "Correct agency.txt name, timezone, URL, language, or contact context in the source export, then rerun validation."
	case "license_contact_metadata":
		return "Review feed_info.txt, license, attribution, and contact values with the administrator before authorized sharing preparation."
	case "missing_required_file":
		return "Regenerate the GTFS ZIP with the missing required file, then rerun import and static validation."
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
