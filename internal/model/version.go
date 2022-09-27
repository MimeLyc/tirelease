package model

import (
	"tirelease/internal/entity"
	"tirelease/internal/repository"
)

type ReleaseVersion struct {
	*entity.ReleaseVersion
	*VersionStateContext
}

func Parse2ReleaseVersion(versionEntity entity.ReleaseVersion) ReleaseVersion {
	releaseVersion := ReleaseVersion{
		ReleaseVersion: &versionEntity,
	}
	versionContext, _ := NewVersionStateContext(&releaseVersion)
	releaseVersion.VersionStateContext = versionContext
	return releaseVersion
}

func (version *ReleaseVersion) ChangeStatus(toStatus entity.ReleaseVersionStatus) error {
	toStateText := StateText(toStatus)
	version.ReleaseVersion.Status = toStatus
	// Update version status first because the next version may rely on the real status.
	repository.UpdateReleaseVersion(version.ReleaseVersion)

	_, err := version.VersionStateContext.Trans(toStateText)

	return err
}

func (version *ReleaseVersion) IsFrozen() bool {
	return version.Status == entity.ReleaseVersionStatusFrozen
}

func initReleaseVersion(version ReleaseVersion) (ReleaseVersion, error) {
	version.Name = version.ComposeVersionName()
	version.Type = version.ComposeVersionType()
	version.Status = InitVersionStatus(version.Type)
	version.ReleaseBranch = version.ComposeVersionBranch()

	return version, nil

}
