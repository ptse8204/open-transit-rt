package main

import (
	"fmt"
	"strings"
)

const expectedConsumerPreparedTargets = 7

type operationsConsumerPreparationView struct {
	Boundary           string
	Status             string
	Summary            string
	OperatorRule       string
	PreparedCount      int
	ExpectedCount      int
	RuntimeRecordCount int
	Targets            []operationsConsumerPreparationTarget
	FutureGates        []operationsConsumerPreparationGate
	Separations        []operationsConsumerPreparationSeparation
	ClaimFlags         operationsConsumerPreparationClaimFlags
}

type operationsConsumerPreparationTarget struct {
	ID              string
	Name            string
	Status          string
	CurrentPath     string
	PacketPath      string
	Meaning         string
	NextAction      string
	DoesNotProve    string
	RuntimeObserved string
}

type operationsConsumerPreparationGate struct {
	ID                    string
	Label                 string
	Status                string
	RequiredAuthorization string
	BlockedAction         string
	DoesNotProve          string
}

type operationsConsumerPreparationSeparation struct {
	ID               string
	Label            string
	Boundary         string
	OperatorHandling string
}

type operationsConsumerPreparationClaimFlags struct {
	ConsumerStatusesChanged    bool
	ConsumerSubmissionClaimed  bool
	ConsumerReviewClaimed      bool
	ConsumerAcceptanceClaimed  bool
	ConsumerIngestionClaimed   bool
	ConsumerListingClaimed     bool
	ConsumerDisplayClaimed     bool
	ExternalContactPerformed   bool
	ExternalEvidenceCreated    bool
	FinalRootEvidenceCreated   bool
	ComplianceClaimed          bool
	ProductionReadinessClaimed bool
	HostedSaaSClaimed          bool
	PublicLaunchClaimed        bool
}

func buildOperationsConsumerPreparation(page operationsPage) operationsConsumerPreparationView {
	targets, preparedCount := consumerPreparationTargets(page.Consumers)
	status := operationsStatusNeedsReview
	if preparedCount == expectedConsumerPreparedTargets && len(targets) == expectedConsumerPreparedTargets {
		status = operationsStatusReady
	}
	return operationsConsumerPreparationView{
		Boundary:           "Private prepared-only consumer packet review. This page explains existing docs-tracker records; it does not create evidence, contact consumers, submit packets, or move tracker status.",
		Status:             status,
		Summary:            fmt.Sprintf("%d of %d docs tracker targets are visible as prepared-only records.", preparedCount, expectedConsumerPreparedTargets),
		OperatorRule:       "Keep every target at prepared unless a separately authorized evidence or consumer-submission phase supplies explicit written authorization and target-originated records.",
		PreparedCount:      preparedCount,
		ExpectedCount:      expectedConsumerPreparedTargets,
		RuntimeRecordCount: len(page.RuntimeConsumers),
		Targets:            targets,
		FutureGates:        consumerPreparationFutureGates(),
		Separations:        consumerPreparationSeparations(),
	}
}

func consumerPreparationTargets(rows []consumerStatusView) ([]operationsConsumerPreparationTarget, int) {
	targets := make([]operationsConsumerPreparationTarget, 0, len(rows))
	preparedCount := 0
	for _, row := range rows {
		status := strings.TrimSpace(row.Status)
		if strings.EqualFold(status, "prepared") {
			preparedCount++
		}
		targets = append(targets, operationsConsumerPreparationTarget{
			ID:           consumerPreparationID(row.Name),
			Name:         row.Name,
			Status:       firstNonEmpty(status, "unknown"),
			CurrentPath:  row.CurrentPath,
			PacketPath:   row.PacketPath,
			Meaning:      "Prepared means the repository has a target-specific packet record path and current pointer for future review.",
			NextAction:   "Review feed URLs and source-of-truth metadata privately; do not submit or change status without separate written authorization.",
			DoesNotProve: "Does not prove submission, review, acceptance, ingestion, listing, display, target approval, compliance, or public launch.",
		})
	}
	return targets, preparedCount
}

func consumerPreparationFutureGates() []operationsConsumerPreparationGate {
	return []operationsConsumerPreparationGate{
		{
			ID:                    "final_root_proof",
			Label:                 "Final-root proof",
			Status:                operationsStatusBlocked,
			RequiredAuthorization: "Requires separate written authorization for retained final-root evidence.",
			BlockedAction:         "No final-root evidence collection or protected evidence writes run from this page.",
			DoesNotProve:          "No final-root readiness or agency-owned-domain proof exists here.",
		},
		{
			ID:                    "consumer_submission",
			Label:                 "Consumer packet submission",
			Status:                operationsStatusBlocked,
			RequiredAuthorization: "Requires separate written authorization before contacting any target or portal.",
			BlockedAction:         "No browser action sends email, opens a portal workflow, uploads files, or performs a network submission.",
			DoesNotProve:          "No target review, ingestion, listing, display, or acceptance exists here.",
		},
		{
			ID:                    "status_movement",
			Label:                 "Tracker status movement",
			Status:                operationsStatusBlocked,
			RequiredAuthorization: "Requires separate written authorization plus target-originated records before editing tracker status.",
			BlockedAction:         "This product track must leave every docs/evidence consumer target at prepared.",
			DoesNotProve:          "Prepared is not a higher consumer state.",
		},
		{
			ID:                    "agency_or_vendor_contact",
			Label:                 "Agency, vendor, or portal contact",
			Status:                operationsStatusBlocked,
			RequiredAuthorization: "Requires separate written authorization and credentials outside this product track.",
			BlockedAction:         "No external contact, vendor compatibility test, hardware certification, or portal automation is performed.",
			DoesNotProve:          "No agency approval, vendor compatibility, or hardware certification exists here.",
		},
	}
}

func consumerPreparationSeparations() []operationsConsumerPreparationSeparation {
	return []operationsConsumerPreparationSeparation{
		{
			ID:               "docs_tracker",
			Label:            "Docs tracker rows",
			Boundary:         "Server-rendered pointers to prepared records under docs/evidence; this page does not edit those files.",
			OperatorHandling: "Use these paths as a private map for future authorized review, not as proof of target action.",
		},
		{
			ID:               "runtime_records",
			Label:            "Runtime workflow records",
			Boundary:         "Database workflow notes are shown separately and do not override the docs tracker prepared state.",
			OperatorHandling: "Treat runtime rows as local deployment notes only.",
		},
		{
			ID:               "feed_readiness",
			Label:            "Feed URL readiness",
			Boundary:         "Feed URL copy/review pages help private operators inspect configured URLs before any external process.",
			OperatorHandling: "Use feed readiness to find missing metadata; do not interpret a URL as target review.",
		},
		{
			ID:               "evidence_gates",
			Label:            "Future evidence gates",
			Boundary:         "Optional final-root, submission, agency pilot, vendor/device, ETA quality, and compliance gates stay separate.",
			OperatorHandling: "Start those gates only after separate written authorization.",
		},
	}
}

func consumerPreparationID(name string) string {
	id := strings.ToLower(strings.TrimSpace(name))
	replacer := strings.NewReplacer(" ", "-", ".", "-", "/", "-", "_", "-", "(", "", ")", "")
	id = replacer.Replace(id)
	id = strings.Trim(id, "-")
	if id == "" {
		return "unknown-target"
	}
	return id
}
