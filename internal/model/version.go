package model

import (
	"tirelease/internal/entity"
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
	_, err := version.VersionStateContext.Trans(toStateText)

	return err
}

func (version *ReleaseVersion) IsFrozen() bool {
	return version.Status == entity.ReleaseVersionStatusFrozen
}

func (version *ReleaseVersion) IsActive() bool {
	return version.Status == entity.ReleaseVersionStatusUpcoming || version.Status == entity.ReleaseVersionStatusFrozen
}

func initReleaseVersion(version ReleaseVersion) (ReleaseVersion, error) {
	version.Name = version.ComposeVersionName()
	version.Type = version.ComposeVersionType()
	version.Status = InitVersionStatus(version.Type)
	version.ReleaseBranch = version.ComposeVersionBranch()

	return version, nil

}
