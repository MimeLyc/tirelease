import PickSelect from "./PickSelect";
import * as React from "react";
import { getAffection } from "./Affection";
import { mapPickStatusToBackend, mapPickStatusToFrontend } from "./mapper"

export function getVersionTriageValue(versionTraige) {
  if (versionTraige === undefined) {
    return "N/A"
  }
  return mapPickStatusToFrontend(versionTraige.triage_result);
}

export function getPickTriageValue(version) {
  return (params) => {
    const affection = getAffection(version)(params);
    if (affection === "N/A" || affection === "no") {
      return <>not affect</>;
    }
    const version_triage = params.row.version_triages?.filter((t) =>
      t.version_name.startsWith(version)
    )[0];
    return getVersionTriageValue(version_triage)
  };
}

export function renderPickTriage(version) {
  return (params) => {

    const affection = getAffection(version)(params);
    if (affection === "N/A" || affection === "no") {
      return <>not affect</>;
    }
    let version_triage = params.row.version_triages?.filter((t) =>
      t.version_name.startsWith(version)
    )[0];
    const pick = version_triage === undefined ? "N/A" : mapPickStatusToFrontend(version_triage.triage_result);
    const patch = version_triage === undefined ? "N/A" : version_triage.version_name;

    const onChange = (value) => {
      if (params.row.version_triages) {
        params.row.version_triages.filter((t) =>
            t.version_name.startsWith(version)
        )[0].triage_result = mapPickStatusToBackend(value); 
      }
    }

    return (
      <>
        <PickSelect
          id={params.row.issue.issue_id}
          version={version}
          patch={patch}
          pick={pick}
          onChange={onChange}
        ></PickSelect>
      </>
    );
  };
}
