package compliance

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	GTFSQualitySourceCanonicalValidator = "canonical_validator"
	GTFSQualitySourceInternalImporter   = "internal_importer"

	GTFSQualityBlocking      = "blocking"
	GTFSQualityNeedsReview   = "needs_review"
	GTFSQualityInformational = "informational"
	GTFSQualityUnknown       = "unknown"

	GTFSQualityMaxGroups       = 100
	GTFSQualityMaxSamples      = 5
	GTFSQualityMaxSampleLength = 400
)

type GTFSQualityTriageInput struct {
	Canonical              *ValidationReportRecord
	InternalImporter       *ValidationReportRecord
	ActiveFeedVersionID    string
	ActiveFeedRevisionTime *time.Time
}

type GTFSQualityTriage struct {
	Canonical        GTFSQualitySection
	InternalImporter GTFSQualitySection
}

type GTFSQualitySection struct {
	Source              string
	SourceLabel         string
	Status              string
	ResultStatus        string
	FeedVersionID       string
	ValidationTimestamp *time.Time
	IsStale             bool
	StaleReason         string
	Groups              []GTFSQualityGroup
	OverflowCount       int
	OperatorSummary     string
	RecommendedAction   string
}

type GTFSQualityGroup struct {
	Source            string
	Family            string
	Codes             []string
	Severity          string
	RiskLevel         string
	Count             int
	OperatorSummary   string
	WhyItMatters      string
	RecommendedAction string
	Samples           []string
	OverflowCount     int
}

type gtfsQualityGroupBuilder struct {
	group GTFSQualityGroup
}

func BuildGTFSQualityTriage(input GTFSQualityTriageInput) GTFSQualityTriage {
	return GTFSQualityTriage{
		Canonical:        buildGTFSQualitySection(input.Canonical, GTFSQualitySourceCanonicalValidator, "Canonical MobilityData static validator", input.ActiveFeedVersionID, input.ActiveFeedRevisionTime),
		InternalImporter: buildGTFSQualitySection(input.InternalImporter, GTFSQualitySourceInternalImporter, "Open Transit RT internal import validation", input.ActiveFeedVersionID, input.ActiveFeedRevisionTime),
	}
}

func buildGTFSQualitySection(record *ValidationReportRecord, source string, label string, activeFeedVersionID string, activeFeedRevisionTime *time.Time) GTFSQualitySection {
	section := GTFSQualitySection{
		Source:            source,
		SourceLabel:       label,
		Status:            GTFSQualityUnknown,
		OperatorSummary:   "No validation result is recorded for this source.",
		RecommendedAction: "Run the appropriate validation or publish/import a schedule before using this as an operator signal.",
	}
	if record == nil {
		return section
	}
	result := record.Result
	createdAt := record.CreatedAt.UTC()
	section.ValidationTimestamp = &createdAt
	section.FeedVersionID = result.FeedVersionID
	section.ResultStatus = result.Status
	section.Status = statusFromValidationResult(result)
	section.OperatorSummary = summaryForSection(source, section.Status, result)
	section.RecommendedAction = actionForSection(source, section.Status)
	section.Groups, section.OverflowCount = groupsFromReport(source, result)
	if activeFeedVersionID != "" && result.FeedVersionID != "" && result.FeedVersionID != activeFeedVersionID {
		section.IsStale = true
		section.Status = maxSeverity(section.Status, GTFSQualityNeedsReview)
		section.StaleReason = "latest result is for a different feed version than the active schedule"
		section.RecommendedAction = "Rerun the static MobilityData validator against the active schedule feed version."
	}
	if activeFeedRevisionTime != nil && section.ValidationTimestamp != nil && section.ValidationTimestamp.Before(activeFeedRevisionTime.UTC()) {
		section.IsStale = true
		section.Status = maxSeverity(section.Status, GTFSQualityNeedsReview)
		section.StaleReason = "latest result predates the active schedule feed timestamp"
		section.RecommendedAction = "Rerun the static MobilityData validator against the active schedule feed version."
	}
	return section
}

func statusFromValidationResult(result ValidationResult) string {
	switch strings.ToLower(strings.TrimSpace(result.Status)) {
	case "failed", "error", "blocked":
		return GTFSQualityBlocking
	case "warning", "warnings":
		return GTFSQualityNeedsReview
	case "passed":
		if result.ErrorCount > 0 {
			return GTFSQualityBlocking
		}
		if result.WarningCount > 0 {
			return GTFSQualityNeedsReview
		}
		if result.InfoCount > 0 {
			return GTFSQualityInformational
		}
		return GTFSQualityInformational
	case "not_run":
		return GTFSQualityBlocking
	default:
		if result.ErrorCount > 0 {
			return GTFSQualityBlocking
		}
		if result.WarningCount > 0 {
			return GTFSQualityNeedsReview
		}
		if result.InfoCount > 0 {
			return GTFSQualityInformational
		}
		return GTFSQualityUnknown
	}
}

func summaryForSection(source string, status string, result ValidationResult) string {
	if source == GTFSQualitySourceCanonicalValidator {
		return fmt.Sprintf("Static validator result is %s with %d error, %d warning, and %d informational notice counts.", status, result.ErrorCount, result.WarningCount, result.InfoCount)
	}
	return fmt.Sprintf("Internal import validation result is %s with %d error, %d warning, and %d informational message counts.", status, result.ErrorCount, result.WarningCount, result.InfoCount)
}

func actionForSection(source string, status string) string {
	if source == GTFSQualitySourceCanonicalValidator {
		switch status {
		case GTFSQualityBlocking:
			return "Fix blocking GTFS issues in the source schedule, then rerun the static MobilityData validator."
		case GTFSQualityNeedsReview:
			return "Review warning groups, decide whether agency data needs correction, then rerun after operator-reviewed source changes."
		default:
			return "Keep this as a diagnostic signal and rerun after the active schedule changes."
		}
	}
	switch status {
	case GTFSQualityBlocking:
		return "Fix import-blocking GTFS structure or references, then import or publish the schedule again."
	case GTFSQualityNeedsReview:
		return "Review importer warnings before relying on the published schedule for downstream feed generation."
	default:
		return "Import or publish a schedule again when source data changes."
	}
}

func groupsFromReport(source string, result ValidationResult) ([]GTFSQualityGroup, int) {
	notices, malformed := extractNoticeMaps(source, result.Report)
	if len(notices) == 0 {
		total := result.ErrorCount + result.WarningCount + result.InfoCount
		if malformed || total > 0 || strings.TrimSpace(result.Status) == "not_run" {
			group := unknownGroup(source, max(1, total))
			group.Samples = appendSample(group.Samples, sampleFromMap(result.Report))
			group.OverflowCount = max(0, group.Count-len(group.Samples))
			return []GTFSQualityGroup{group}, 0
		}
		return nil, 0
	}
	builders := map[string]*gtfsQualityGroupBuilder{}
	for _, notice := range notices {
		code := noticeCode(notice)
		severity := noticeSeverity(notice, result)
		family := noticeFamily(code, notice)
		key := severity + "\x00" + family + "\x00" + code
		builder := builders[key]
		if builder == nil {
			guidance := guidanceForNotice(family, code, severity)
			builder = &gtfsQualityGroupBuilder{group: GTFSQualityGroup{
				Source:            source,
				Family:            family,
				Codes:             []string{code},
				Severity:          severity,
				RiskLevel:         riskLevelForNotice(family, severity),
				OperatorSummary:   guidance.summary,
				WhyItMatters:      guidance.why,
				RecommendedAction: guidance.action,
			}}
			builders[key] = builder
		}
		builder.group.Count += noticeCount(notice)
		if len(builder.group.Samples) < GTFSQualityMaxSamples {
			builder.group.Samples = appendSample(builder.group.Samples, sampleFromMap(notice))
		}
	}
	groups := make([]GTFSQualityGroup, 0, len(builders))
	for _, builder := range builders {
		builder.group.OverflowCount = max(0, builder.group.Count-len(builder.group.Samples))
		groups = append(groups, builder.group)
	}
	sortGTFSQualityGroups(groups)
	overflow := 0
	if len(groups) > GTFSQualityMaxGroups {
		for _, group := range groups[GTFSQualityMaxGroups:] {
			overflow += group.Count
		}
		groups = groups[:GTFSQualityMaxGroups]
	}
	return groups, overflow
}

func extractNoticeMaps(source string, report map[string]any) ([]map[string]any, bool) {
	if len(report) == 0 {
		return nil, true
	}
	if source == GTFSQualitySourceInternalImporter {
		var notices []map[string]any
		for _, entry := range []struct {
			key      string
			severity string
		}{
			{"errors", GTFSQualityBlocking},
			{"warnings", GTFSQualityNeedsReview},
			{"info", GTFSQualityInformational},
		} {
			items, ok := report[entry.key].([]any)
			if !ok {
				continue
			}
			for _, item := range items {
				if notice, ok := item.(map[string]any); ok {
					copyNotice := shallowNoticeCopy(notice)
					copyNotice["severity"] = entry.severity
					notices = append(notices, copyNotice)
				} else {
					notices = append(notices, map[string]any{"code": "unknown", "severity": GTFSQualityUnknown, "message": fmt.Sprintf("%T", item)})
				}
			}
		}
		return notices, len(notices) == 0
	}
	raw, ok := report["raw_report"]
	if !ok {
		return nil, true
	}
	rawMap, ok := raw.(map[string]any)
	if !ok {
		return nil, true
	}
	rawNotices, ok := rawMap["notices"]
	if !ok {
		return nil, true
	}
	items, ok := rawNotices.([]any)
	if !ok {
		return nil, true
	}
	notices := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if notice, ok := item.(map[string]any); ok {
			notices = append(notices, shallowNoticeCopy(notice))
		} else {
			notices = append(notices, map[string]any{"code": "unknown", "severity": GTFSQualityUnknown, "message": fmt.Sprintf("%T", item)})
		}
	}
	return notices, false
}

func shallowNoticeCopy(in map[string]any) map[string]any {
	out := make(map[string]any, min(len(in), 12))
	for key, value := range in {
		lower := strings.ToLower(key)
		if lower == "raw_report" || lower == "stdout" || lower == "stderr" || lower == "argv" || lower == "args" || strings.Contains(lower, "path") {
			continue
		}
		out[key] = value
	}
	return out
}

func noticeCode(notice map[string]any) string {
	for _, key := range []string{"code", "notice_code", "type", "name"} {
		if value := strings.TrimSpace(scalarString(notice[key])); value != "" {
			return strings.ToLower(value)
		}
	}
	return "unknown"
}

func noticeSeverity(notice map[string]any, result ValidationResult) string {
	raw := strings.ToLower(strings.TrimSpace(scalarString(notice["severity"])))
	switch raw {
	case "error", "errors", GTFSQualityBlocking:
		return GTFSQualityBlocking
	case "warning", "warn", GTFSQualityNeedsReview:
		return GTFSQualityNeedsReview
	case "info", "information", GTFSQualityInformational:
		return GTFSQualityInformational
	case GTFSQualityUnknown:
		return GTFSQualityUnknown
	default:
		return statusFromValidationResult(result)
	}
}

func noticeCount(notice map[string]any) int {
	for _, key := range []string{"totalNotices", "total_notices", "count"} {
		switch typed := notice[key].(type) {
		case int:
			if typed > 0 {
				return typed
			}
		case int64:
			if typed > 0 {
				return int(typed)
			}
		case float64:
			if typed > 0 {
				return int(typed)
			}
		case string:
			parsed, err := strconv.Atoi(strings.TrimSpace(typed))
			if err == nil && parsed > 0 {
				return parsed
			}
		}
	}
	return 1
}

func noticeFamily(code string, notice map[string]any) string {
	lower := strings.ToLower(code + " " + scalarString(notice["message"]) + " " + scalarString(notice["filename"]) + " " + scalarString(notice["fileName"]) + " " + scalarString(notice["file"]))
	switch {
	case strings.Contains(lower, "missing_required") || strings.Contains(lower, "missing required") || strings.Contains(lower, "required_file") || strings.Contains(lower, "required file") || strings.Contains(lower, "file_not_found"):
		return "missing_required_file"
	case strings.Contains(lower, "expired_calendar"):
		return "expired_calendar"
	case strings.Contains(lower, "route_short_name_too_long"):
		return "route_short_name_too_long"
	case strings.Contains(lower, "unused_shape"):
		return "unused_shape"
	case strings.Contains(lower, "foreign") || strings.Contains(lower, "missing") && strings.Contains(lower, "reference"):
		return "missing_or_foreign_key_reference"
	case strings.Contains(lower, "agency") && (strings.Contains(lower, "name") || strings.Contains(lower, "timezone") || strings.Contains(lower, "url") || strings.Contains(lower, "lang") || strings.Contains(lower, "phone")):
		return "agency_metadata"
	case strings.Contains(lower, "feed_info") || strings.Contains(lower, "license") || strings.Contains(lower, "contact") || strings.Contains(lower, "attribution"):
		return "license_contact_metadata"
	case strings.Contains(lower, "stop_time") || strings.Contains(lower, "stop_times") || strings.Contains(lower, "arrival_time") || strings.Contains(lower, "departure_time"):
		return "bad_stop_times"
	case strings.Contains(lower, "duplicate"):
		return "duplicate_ids"
	case strings.Contains(lower, "calendar") || strings.Contains(lower, "service_date") || strings.Contains(lower, "service_id"):
		return "calendar_service_dates"
	case strings.Contains(lower, "shape_dist") || strings.Contains(lower, "shape") && strings.Contains(lower, "order"):
		return "shape_ordering"
	case strings.Contains(lower, "frequenc"):
		return "frequency_issues"
	case strings.Contains(lower, "block_id") || strings.Contains(lower, "block transition") || strings.Contains(lower, "block"):
		return "block_transition_issues"
	case strings.Contains(lower, "route"):
		return "route_metadata"
	case strings.Contains(lower, "stop_lat") || strings.Contains(lower, "stop_lon") || strings.Contains(lower, "stop_location") || strings.Contains(lower, "location_type") || strings.Contains(lower, "parent_station"):
		return "stop_location"
	default:
		return "unknown"
	}
}

type noticeGuidance struct {
	summary string
	why     string
	action  string
}

func guidanceForNotice(family string, code string, severity string) noticeGuidance {
	switch family {
	case "missing_required_file":
		return noticeGuidance{"A required GTFS file appears missing from the schedule package.", "Missing required files can block import, validation, feed discovery, and realtime matching setup.", "Regenerate the GTFS ZIP with all required files and rerun import plus static validation."}
	case "expired_calendar":
		return noticeGuidance{"Service calendar dates appear expired or out of useful range.", "Expired service can make a schedule look unavailable even when routes and trips exist.", "Confirm service dates with the agency schedule owner and update calendar or calendar_dates before rerunning validation."}
	case "agency_metadata":
		return noticeGuidance{"Agency metadata needs review.", "Agency name, timezone, language, URL, and contact context affect feed identity and downstream interpretation.", "Correct agency.txt metadata in the source system, then re-export and rerun validation."}
	case "license_contact_metadata":
		return noticeGuidance{"Feed license or contact metadata needs review.", "License and contact metadata help deployment owners prepare safe external sharing without implying acceptance.", "Review feed_info.txt, license, attribution, and contact values with the administrator before sharing URLs."}
	case "route_short_name_too_long":
		return noticeGuidance{"A route short name is longer than expected for rider-facing display.", "Long short names can render poorly in maps, signs, and downstream consumer displays.", "Review route naming with the operator and move descriptive text to route_long_name when appropriate."}
	case "route_metadata":
		return noticeGuidance{"Route metadata needs review.", "Route type, names, colors, URLs, or agency links affect schedule usability and downstream display.", "Review routes.txt with the route naming owner and correct source metadata before rerunning validation."}
	case "stop_location":
		return noticeGuidance{"Stop location or station metadata needs review.", "Bad stop coordinates or parent/location fields can break maps, rider guidance, and matching context.", "Review stops.txt coordinates and station hierarchy with the GIS or stop inventory owner."}
	case "unused_shape":
		return noticeGuidance{"A shape record is not referenced by active trips.", "Unused shapes add noise and can hide shape maintenance mistakes.", "Remove unused shapes only after confirming no planned trip should reference them."}
	case "missing_or_foreign_key_reference":
		return noticeGuidance{"A GTFS row references an ID that is missing from the related file.", "Broken references can block import, matching, or consumer parsing.", "Fix the referenced ID in the source GTFS file or add the missing related record, then rerun validation."}
	case "bad_stop_times":
		return noticeGuidance{"stop_times.txt contains invalid, inconsistent, or incomplete timing data.", "Bad stop times can break trip matching, realtime alignment, and trip presentation.", "Correct arrival/departure times, stop sequence ordering, and after-midnight time formatting in the source schedule."}
	case "duplicate_ids":
		return noticeGuidance{"A GTFS file contains duplicate identifiers where unique IDs are expected.", "Duplicate IDs make references ambiguous and can block deterministic imports.", "Deduplicate IDs in the source GTFS and update references consistently."}
	case "calendar_service_dates":
		return noticeGuidance{"Calendar or service-date records need review.", "Service availability controls whether trips are active on an agency-local service day.", "Check calendar.txt and calendar_dates.txt for coverage, exceptions, and agency-confirmed service periods."}
	case "shape_ordering":
		return noticeGuidance{"Shape ordering or shape_dist_traveled values need review.", "Shape ordering affects map drawing and vehicle progress along a trip.", "Verify shape point order and distance progression before approving source GTFS changes."}
	case "frequency_issues":
		return noticeGuidance{"Frequency-based service records need review.", "Frequency issues can affect repeated trip instances and realtime assignment windows.", "Review frequencies.txt start/end times, headways, and trip references."}
	case "block_transition_issues":
		return noticeGuidance{"Block or block transition data needs review.", "Block continuity affects vehicle trip chaining and conservative matching.", "Confirm block_id values and trip ordering with operations staff before changing source GTFS."}
	default:
		return noticeGuidance{fmt.Sprintf("Unrecognized validator notice %q has %s severity.", code, severity), "Unknown notices still need operator review because their downstream impact is not classified yet.", "Inspect the source validator report outside this page, classify the notice if it recurs, and rerun after operator-reviewed changes."}
	}
}

func riskLevelForNotice(family string, severity string) string {
	switch severity {
	case GTFSQualityBlocking:
		return "blocks import or reliable feed use"
	case GTFSQualityNeedsReview:
		switch family {
		case "expired_calendar", "calendar_service_dates", "missing_required_file", "missing_or_foreign_key_reference", "bad_stop_times", "frequency_issues", "block_transition_issues":
			return "can break service availability or realtime usefulness"
		case "agency_metadata", "license_contact_metadata":
			return "can block sharing preparation or operator trust"
		case "route_metadata", "stop_location", "shape_ordering", "unused_shape":
			return "can degrade maps, matching, or downstream display"
		default:
			return "needs source-owner review before relying on the feed"
		}
	case GTFSQualityInformational:
		return "track during normal data-quality review"
	default:
		return "unclassified impact; maintainer review needed"
	}
}

func unknownGroup(source string, count int) GTFSQualityGroup {
	guidance := guidanceForNotice("unknown", "unknown", GTFSQualityUnknown)
	return GTFSQualityGroup{
		Source:            source,
		Family:            "unknown",
		Codes:             []string{"unknown"},
		Severity:          GTFSQualityUnknown,
		RiskLevel:         riskLevelForNotice("unknown", GTFSQualityUnknown),
		Count:             count,
		OperatorSummary:   guidance.summary,
		WhyItMatters:      guidance.why,
		RecommendedAction: guidance.action,
	}
}

func sortGTFSQualityGroups(groups []GTFSQualityGroup) {
	sort.SliceStable(groups, func(i, j int) bool {
		a := groups[i]
		b := groups[j]
		if severityRank(a.Severity) != severityRank(b.Severity) {
			return severityRank(a.Severity) < severityRank(b.Severity)
		}
		if a.Family != b.Family {
			return a.Family < b.Family
		}
		ac := ""
		if len(a.Codes) > 0 {
			ac = a.Codes[0]
		}
		bc := ""
		if len(b.Codes) > 0 {
			bc = b.Codes[0]
		}
		if ac != bc {
			return ac < bc
		}
		if a.Count != b.Count {
			return a.Count > b.Count
		}
		return a.Source < b.Source
	})
}

func severityRank(severity string) int {
	switch severity {
	case GTFSQualityBlocking:
		return 0
	case GTFSQualityNeedsReview:
		return 1
	case GTFSQualityInformational:
		return 2
	default:
		return 3
	}
}

func maxSeverity(a string, b string) string {
	if severityRank(a) <= severityRank(b) {
		return a
	}
	return b
}

func sampleFromMap(values map[string]any) string {
	if len(values) == 0 {
		return ""
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		lower := strings.ToLower(key)
		if lower == "raw_report" || lower == "stdout" || lower == "stderr" || lower == "argv" || lower == "args" || strings.Contains(lower, "path") || strings.Contains(lower, "command") {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, min(len(keys), 6))
	for _, key := range keys {
		value := scalarString(values[key])
		if value == "" {
			continue
		}
		parts = append(parts, key+"="+value)
		if len(parts) >= 6 {
			break
		}
	}
	return truncateSample(strings.Join(parts, ", "))
}

func appendSample(samples []string, sample string) []string {
	sample = truncateSample(sample)
	if sample == "" {
		return samples
	}
	return append(samples, sample)
}

func scalarString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return scrubPrivateText(typed)
	case fmt.Stringer:
		return typed.String()
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return fmt.Sprintf("%T", typed)
	}
}

func truncateSample(value string) string {
	value = scrubPrivateText(value)
	value = strings.TrimSpace(value)
	if len(value) <= GTFSQualityMaxSampleLength {
		return value
	}
	return value[:GTFSQualityMaxSampleLength] + "..."
}

func scrubPrivateText(value string) string {
	for _, marker := range []string{"/tmp/private", "/var/folders", "/private/var", "\\tmp\\private", "/users/private", "/users/", "/var/lib", "/etc/"} {
		value = strings.ReplaceAll(value, marker, "{private_path}")
		value = strings.ReplaceAll(value, strings.ToUpper(marker), "{private_path}")
	}
	redacted := make([]string, 0, 8)
	for _, field := range strings.Fields(value) {
		lower := strings.ToLower(field)
		switch {
		case strings.Contains(lower, "/tmp/private"), strings.Contains(lower, "/users/"), strings.Contains(lower, "/var/lib"), strings.Contains(lower, "/etc/"), strings.Contains(lower, "\\tmp\\private"):
			redacted = append(redacted, "{private_path}")
		case strings.Contains(lower, "authorization:"), strings.Contains(lower, "bearer"), strings.Contains(lower, "token="), strings.Contains(lower, "secret"), strings.Contains(lower, "password"), strings.Contains(lower, "cookie"), strings.Contains(lower, "admin_session"):
			redacted = append(redacted, "{redacted_secret}")
		case strings.Contains(lower, "postgres://"), strings.Contains(lower, "postgresql://"), strings.Contains(lower, "database_url"), strings.Contains(lower, "restore_database_url"):
			redacted = append(redacted, "{redacted_database}")
		case strings.HasPrefix(lower, "http://"), strings.HasPrefix(lower, "https://"), strings.Contains(lower, "webhook"):
			redacted = append(redacted, "{redacted_url}")
		default:
			redacted = append(redacted, field)
		}
		if len(redacted) >= 64 {
			redacted = append(redacted, "{truncated}")
			break
		}
	}
	return strings.Join(redacted, " ")
}
