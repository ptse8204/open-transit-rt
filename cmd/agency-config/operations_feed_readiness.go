package main

import (
	"fmt"
	"strings"
)

type operationsFeedReadinessView struct {
	Boundary        string                              `json:"boundary"`
	Rows            []operationsFeedReadinessRow        `json:"rows"`
	Metadata        []operationsFeedReadinessMetadata   `json:"metadata"`
	FutureChecklist []operationsFeedReadinessFutureGate `json:"future_checklist"`
	ClaimFlags      operationsFeedReadinessClaimFlags   `json:"claim_flags"`
}

type operationsFeedReadinessRow struct {
	ID                 string   `json:"id"`
	Label              string   `json:"label"`
	Status             string   `json:"status"`
	PublicPath         string   `json:"public_path"`
	ConfiguredURL      string   `json:"configured_url"`
	CopyValue          string   `json:"copy_value"`
	MetadataSource     string   `json:"metadata_source"`
	MetadataStatus     string   `json:"metadata_status"`
	ValidationContext  string   `json:"validation_context"`
	PublicFetchContext string   `json:"public_fetch_context"`
	Meaning            string   `json:"meaning"`
	CopyGuidance       string   `json:"copy_guidance"`
	DoesNotProve       string   `json:"does_not_prove"`
	ReviewChecklist    []string `json:"review_checklist"`
	DocsLink           string   `json:"docs_link"`
}

type operationsFeedReadinessMetadata struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Status        string `json:"status"`
	CurrentSignal string `json:"current_signal"`
	NextAction    string `json:"next_action"`
	DoesNotProve  string `json:"does_not_prove"`
}

type operationsFeedReadinessFutureGate struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	CurrentStatus string `json:"current_status"`
	NextAction    string `json:"next_action"`
	Boundary      string `json:"boundary"`
}

type operationsFeedReadinessClaimFlags struct {
	ExternalEvidenceCreated       bool `json:"external_evidence_created"`
	FinalRootEvidenceCreated      bool `json:"final_root_evidence_created"`
	ConsumerStatusesChanged       bool `json:"consumer_statuses_changed"`
	ConsumerAcceptanceClaimed     bool `json:"consumer_acceptance_claimed"`
	ComplianceClaimed             bool `json:"compliance_claimed"`
	ProductionReadinessClaimed    bool `json:"production_readiness_claimed"`
	PublicLaunchClaimed           bool `json:"public_launch_claimed"`
	HostedSaaSClaimed             bool `json:"hosted_saas_claimed"`
	VendorCompatibilityClaimed    bool `json:"vendor_compatibility_claimed"`
	HardwareCertificationClaimed  bool `json:"hardware_certification_claimed"`
	ProductionGradeETAClaimed     bool `json:"production_grade_eta_claimed"`
	RealWorldETAAccuracyClaimed   bool `json:"real_world_eta_accuracy_claimed"`
	SLACoverageClaimed            bool `json:"sla_coverage_claimed"`
	UptimeGuaranteeClaimed        bool `json:"uptime_guarantee_claimed"`
	ConsumerSubmissionClaimed     bool `json:"consumer_submission_claimed"`
	FinalRootReadinessClaimed     bool `json:"final_root_readiness_claimed"`
	ExternalBrowserFetchPerformed bool `json:"external_browser_fetch_performed"`
}

func buildOperationsFeedReadiness(page operationsPage) operationsFeedReadinessView {
	return operationsFeedReadinessView{
		Boundary:        "Private public-feed readiness review only. This page reads configured metadata and local diagnostic summaries; it does not fetch final roots, contact consumers, write evidence, change consumer statuses, or prove final-root, compliance, public launch, hosted service, production, vendor, hardware, SLA, uptime, or ETA-quality outcomes.",
		Rows:            feedReadinessRows(page),
		Metadata:        feedReadinessMetadata(page),
		FutureChecklist: feedReadinessFutureGates(),
		ClaimFlags:      operationsFeedReadinessClaimFlags{},
	}
}

func feedReadinessRows(page operationsPage) []operationsFeedReadinessRow {
	firstRunRows := firstRunFeedURLs(page)
	rows := make([]operationsFeedReadinessRow, 0, len(firstRunRows))
	for _, source := range firstRunRows {
		rows = append(rows, operationsFeedReadinessRow{
			ID:                 source.ID,
			Label:              source.Label,
			Status:             feedReadinessRowStatus(page, source.ID, source.URL),
			PublicPath:         feedHealthPublicPath(source.ID),
			ConfiguredURL:      firstNonEmpty(source.URL, "missing"),
			CopyValue:          source.CopyValue,
			MetadataSource:     source.Source,
			MetadataStatus:     feedReadinessMetadataStatus(page, source.ID),
			ValidationContext:  feedReadinessValidationContext(page, source.ID),
			PublicFetchContext: feedReadinessFetchContext(page, source.ID),
			Meaning:            source.Meaning,
			CopyGuidance:       "Copy this configured URL only after the source-of-truth metadata checklist below has been reviewed. Copying a URL from the private console is not outside proof.",
			DoesNotProve:       source.DoesNotProve,
			ReviewChecklist:    feedReadinessChecklist(source.ID),
			DocsLink:           source.DocsLink,
		})
	}
	return rows
}

func feedReadinessRowStatus(page operationsPage, id string, url string) string {
	if strings.TrimSpace(url) == "" || page.DiscoveryError != "" {
		return operationsStatusMissing
	}
	switch id {
	case "feeds_json":
		if page.Discovery.Readiness.AllRequiredFeedsListed && page.Discovery.Readiness.LicenseComplete && page.Discovery.Readiness.ContactComplete && page.Discovery.Readiness.HTTPSURLs {
			return operationsStatusReady
		}
	default:
		if feed, ok := feedHealthMetadata(page, id); ok && strings.TrimSpace(feed.CanonicalPublicURL) != "" {
			return operationsStatusReady
		}
	}
	return operationsStatusNeedsReview
}

func feedReadinessMetadataStatus(page operationsPage, id string) string {
	if page.DiscoveryError != "" {
		return "publication metadata is not configured"
	}
	if id == "feeds_json" {
		return fmt.Sprintf("public_base_url=%t; license=%t; contact=%t; all_required_listed=%t; https=%t; discoverable=%t",
			strings.TrimSpace(page.Discovery.PublicBaseURL) != "",
			page.Discovery.Readiness.LicenseComplete,
			page.Discovery.Readiness.ContactComplete,
			page.Discovery.Readiness.AllRequiredFeedsListed,
			page.Discovery.Readiness.HTTPSURLs,
			page.Discovery.Readiness.Discoverable,
		)
	}
	feed, ok := feedHealthMetadata(page, id)
	if !ok {
		return "published feed metadata missing"
	}
	return fmt.Sprintf("activation=%s; license=%t; contact=%t; active_feed_version=%s",
		firstNonEmpty(feed.ActivationStatus, "unknown"),
		strings.TrimSpace(firstNonEmpty(feed.LicenseName, page.Discovery.License.Name)) != "",
		strings.TrimSpace(firstNonEmpty(feed.ContactEmail, page.Discovery.TechnicalContactEmail)) != "",
		firstNonEmpty(feed.ActiveFeedVersionID, "not available"),
	)
}

func feedReadinessValidationContext(page operationsPage, id string) string {
	if id == "feeds_json" {
		return "feeds.json is metadata, not a GTFS validator artifact"
	}
	row := feedHealthValidationRow(page, id)
	if row == nil {
		return "no private validator-health row recorded"
	}
	return fmt.Sprintf("latest=%s; stale=%s; validator=%s; latest_feed_version=%s",
		firstNonEmpty(row.LatestResultStatus, "not available"),
		firstNonEmpty(row.StaleStatus, "unknown"),
		firstNonEmpty(row.ValidatorID, "not available"),
		firstNonEmpty(row.LatestResultFeedVersionID, "not available"),
	)
}

func feedReadinessFetchContext(page operationsPage, id string) string {
	if id == "feeds_json" {
		return "no separate local public fetch snapshot is recorded for feeds.json"
	}
	row := feedHealthReliabilityRow(page, id)
	if row == nil {
		return "local public fetch status is not recorded"
	}
	return fmt.Sprintf("endpoint_available=%s; health=%s; snapshot=%s",
		boolPtrText(row.EndpointAvailable),
		firstNonEmpty(row.Status, "unknown"),
		formatTimeForText(row.SnapshotAt),
	)
}

func feedReadinessChecklist(id string) []string {
	if id == "feeds_json" {
		return []string{
			"Confirm public base URL is the intended operator-reviewed root.",
			"Confirm feeds.json lists schedule, Vehicle Positions, Trip Updates, and Alerts.",
			"Confirm contact and open-license metadata are complete.",
			"Confirm future source-of-truth website and final-root evidence remain separate authorization gates.",
		}
	}
	return []string{
		"Confirm the configured URL matches the expected public path.",
		"Review validator-health and local feed-health context before sharing.",
		"Confirm license/contact metadata remains visible through feeds.json.",
		"Keep consumer packet status prepared-only unless separate authorized retained proof exists.",
	}
}

func feedReadinessMetadata(page operationsPage) []operationsFeedReadinessMetadata {
	return []operationsFeedReadinessMetadata{
		feedReadinessMetadataRow("public_base_url", "Public base URL", strings.TrimSpace(page.Discovery.PublicBaseURL) != "" && page.DiscoveryError == "", firstNonEmpty(page.Discovery.PublicBaseURL, page.DiscoveryError, "missing"), "Set publication metadata before copying URLs outside the private console.", "Does not prove final-root ownership or public website source-of-truth listing."),
		feedReadinessMetadataRow("license", "License metadata", page.Discovery.Readiness.LicenseComplete && page.DiscoveryError == "", firstNonEmpty(page.Discovery.License.Name, "missing"), "Add license name and URL before external feed review.", "Does not prove legal approval or consumer acceptance."),
		feedReadinessMetadataRow("contact", "Technical contact metadata", page.Discovery.Readiness.ContactComplete && page.DiscoveryError == "", firstNonEmpty(page.Discovery.TechnicalContactEmail, "missing"), "Add a monitored technical contact before external feed review.", "Does not prove agency approval or managed support."),
		feedReadinessMetadataRow("https", "HTTPS configured URLs", page.Discovery.Readiness.HTTPSURLs && page.DiscoveryError == "", fmt.Sprintf("all_https=%t", page.Discovery.Readiness.HTTPSURLs), "Review any HTTP/local URL before using it outside local/reference contexts.", "Does not prove uptime, SLA, or hosted service availability."),
		feedReadinessMetadataRow("all_required_feeds", "Expected feed set", page.Discovery.Readiness.AllRequiredFeedsListed && page.DiscoveryError == "", fmt.Sprintf("%d feed records; all_required_listed=%t", len(page.Discovery.Feeds), page.Discovery.Readiness.AllRequiredFeedsListed), "List feeds.json, schedule, Vehicle Positions, Trip Updates, and Alerts before public feed review.", "Does not prove consumer ingestion, listing, display, or acceptance."),
	}
}

func feedReadinessMetadataRow(id string, label string, ok bool, signal string, next string, doesNotProve string) operationsFeedReadinessMetadata {
	status := operationsStatusNeedsReview
	if ok {
		status = operationsStatusReady
	}
	if strings.Contains(strings.ToLower(signal), "missing") || strings.TrimSpace(signal) == "" {
		status = operationsStatusMissing
	}
	return operationsFeedReadinessMetadata{
		ID:            id,
		Label:         label,
		Status:        status,
		CurrentSignal: signal,
		NextAction:    next,
		DoesNotProve:  doesNotProve,
	}
}

func feedReadinessFutureGates() []operationsFeedReadinessFutureGate {
	boundary := "Requires separate written authorization before collection, retention, portal contact, status movement, or stronger public wording."
	return []operationsFeedReadinessFutureGate{
		{ID: "source_of_truth_listing", Label: "Source-of-truth website listing", CurrentStatus: operationsStatusBlocked, NextAction: "Prepare a checklist only; do not fetch or retain final-root proof in Phase 87.", Boundary: boundary},
		{ID: "final_root_proof", Label: "Final-root evidence", CurrentStatus: operationsStatusBlocked, NextAction: "Keep final-root proof as a future optional evidence gate.", Boundary: boundary},
		{ID: "consumer_packet_use", Label: "Prepared consumer packet use", CurrentStatus: operationsStatusBlocked, NextAction: "Keep packet records prepared-only and do not modify protected packet/status files.", Boundary: boundary},
		{ID: "off_host_validation", Label: "Off-host validation record", CurrentStatus: operationsStatusNeedsReview, NextAction: "Document how an operator may run validation from their environment without turning output into retained evidence.", Boundary: boundary},
	}
}

func boolPtrText(value *bool) string {
	if value == nil {
		return "not observed"
	}
	if *value {
		return "true"
	}
	return "false"
}
