package git

// ============================================================= Constants
const CrossReferencedEvent = "CrossReferencedEvent"
const IssueComment = "IssueComment"
const CherryPickLabel = "cherry-pick-approved"
const NotCheryyPickLabel = "do-not-merge/cherry-pick-not-approved"
const LGT2Label = "status/LGT2"
const BugTypeLabel = "type/bug"
const AffectsLabel = "affects-%s"
const AffectsPrefixLabel = "affects-"
const MayAffectsLabel = "may-affects-%s"
const MayAffectsPrefixLabel = "may-affects-"
const SeverityLabel = "severity/"
const SeverityCriticalLabel = "severity/critical"
const SeverityMajorLabel = "severity/major"
const NoneReleaseNoteLabel = "release-note-none"
const TypeLabel = "type/"
const MergeableStateMergeable = "mergeable"
const MergeableStateUnstable = "unstable"
const MergeableStateUnknown = "unknown"
const MergeRetryComment = "/merge"
const OpenStatus = "open"
const ReleaseBranchPrefix = "release-"
const HeadRefPrefix = "refs/heads/"

type RefType string

const (
	RefTypeTag    = RefType("tag")
	RefTypeBranch = RefType("branch")
)
