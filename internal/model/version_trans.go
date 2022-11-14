package model

import (
	"fmt"
	"tirelease/internal/entity"
)

var VersionStateTransMap = make(TransitionMap[*versionStateContext])

// Orders matters, the trans with from and to should be added firstly.
func init() {
	if len(VersionStateTransMap) > 0 {
		return
	}

	VersionStateTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.ReleaseVersionStatusFrozen),
	}] = VersionState2Frozen{}

	VersionStateTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.ReleaseVersionStatusUpcoming),
	}] = VersionState2Upcoming{}

	VersionStateTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.ReleaseVersionStatusReleased),
	}] = VersionState2Released{}

	VersionStateTransMap[StateTransitionMeta{
		FromState: EmptyStateText(),
		ToState:   StateText(entity.ReleaseVersionStatusCancelled),
	}] = VersionState2Cancelled{}
}

// Frozen unfinished issue triages
type VersionState2Frozen struct{}

func (trans VersionState2Frozen) FitConstraints(context *versionStateContext) (bool, error) {
	return true, nil
}

// Refrozen approved issues while the version changed to frozen.
func (trans VersionState2Frozen) Effect(context *versionStateContext) (bool, error) {
	issueTriages, err := context.Version.SelectHistoryIssueTriages()
	if err != nil {
		return false, nil
	}

	for _, triage := range issueTriages {
		stateResult := entity.VersionTriageResult(triage.PickTriage.State.getStateText())
		if stateResult == entity.VersionTriageResultAccept {

			if triage.GetMergeStatus() != entity.VersionTriageMergeStatusMerged {
				triage.TriagePickStatus(entity.VersionTriageResultAcceptFrozen)
			}
			err = CreateOrUpdateVersionTriageInfo(&triage)
			if err != nil {
				fmt.Printf("Change status of triage %d to accept frozen error %v", triage.ID, err)
			}
		}
	}

	return true, nil
}

// Create next version and inherit triages
type VersionState2Released struct{}

// TODO: If last patch version has not been released , return false
func (trans VersionState2Released) FitConstraints(context *versionStateContext) (bool, error) {
	return true, nil
}

func (trans VersionState2Released) Effect(context *versionStateContext) (bool, error) {
	if context.Version.Type == entity.ReleaseVersionTypeHotfix {
		return true, nil
	}
	nextVersion, err := CreateNextVersionIfNotExist(context.Version)
	if err != nil || nextVersion.ReleaseVersion == nil {
		return false, err
	}

	// Inherit unfinished triages to next version
	issueTriages, err := context.Version.SelectHistoryIssueTriages()
	if err != nil {
		return false, nil
	}
	for _, triage := range issueTriages {
		stateResult := entity.VersionTriageResult(triage.PickTriage.State.getStateText())

		// Deal with released issues
		if stateResult == entity.VersionTriageResultAccept ||
			stateResult == entity.VersionTriageResultAcceptFrozen {

			if triage.GetMergeStatus() == entity.VersionTriageMergeStatusMerged {
				triage.TriagePickStatus(entity.VersionTriageResultReleased)
			} else {
				triage.Version = nextVersion
			}

		} else if stateResult == entity.VersionTriageResultUnKnown ||
			stateResult == entity.VersionTriageResultLater {

			triage.Version = nextVersion

		}

		// TODO add logs while update error dumped
		err = CreateOrUpdateVersionTriageInfo(&triage)
		if err != nil {
			fmt.Printf("Change status of triage %d to accept error %v", triage.ID, err)
		}
	}

	// Change the status of next patch version to be upcoming
	nextVersion.ChangeStatus(entity.ReleaseVersionStatusUpcoming)

	return true, nil
}

type VersionState2Cancelled struct{}

func (trans VersionState2Cancelled) FitConstraints(context *versionStateContext) (bool, error) {
	return true, nil
}

func (trans VersionState2Cancelled) Effect(context *versionStateContext) (bool, error) {
	version := context.Version
	if version.Version.Type == entity.ReleaseVersionTypeHotfix ||
		version.Version.Type == entity.ReleaseVersionTypeMinor {
		return true, nil
	}

	lastPatch, err := SelectPrePatchVersion(*context.Version)
	if err != nil {
		return false, err
	}
	if lastPatch == nil {
		return false, fmt.Errorf("Last sprint of patch %s not founded.", version.Version.Name)
	}

	issueTriages, err := context.Version.SelectHistoryIssueTriages()
	if err != nil {
		return false, nil
	}

	for _, triage := range issueTriages {
		triage.Version = lastPatch
		err = CreateOrUpdateVersionTriageInfo(&triage)
		if err != nil {
			fmt.Printf("Change version of triage %d to last patch error %v", triage.ID, err)
		}
	}

	return true, nil

}

type VersionState2Upcoming struct{}

func (trans VersionState2Upcoming) FitConstraints(context *versionStateContext) (bool, error) {
	version := context.Version
	if version.Version.Type == entity.ReleaseVersionTypeHotfix ||
		version.Version.Type == entity.ReleaseVersionTypeMinor {
		return true, nil
	}

	lastPatch, err := SelectPrePatchVersion(*context.Version)
	if err != nil {
		return false, err
	}
	if lastPatch == nil || lastPatch.Status != entity.ReleaseVersionStatusReleased {
		return false, fmt.Errorf("Last sprint of patch %s has not been released.", version.Version.Name)
	}
	return true, nil
}

func (trans VersionState2Upcoming) Effect(context *versionStateContext) (bool, error) {
	issueTriages, err := context.Version.SelectHistoryIssueTriages()
	if err != nil {
		return false, nil
	}

	for _, triage := range issueTriages {
		stateResult := entity.VersionTriageResult(triage.PickTriage.State.getStateText())
		if stateResult == entity.VersionTriageResultAcceptFrozen {
			triage.TriagePickStatus(entity.VersionTriageResultAccept)
			err = CreateOrUpdateVersionTriageInfo(&triage)
			if err != nil {
				fmt.Printf("Change status of triage %d to accept error %v", triage.ID, err)
			}
		}
	}

	return true, nil
}
