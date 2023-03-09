import * as React from "react";
import { useMutation, useQueryClient } from "react-query";
import { url } from "../../utils";
import axios from "axios";

import { HotfixTriagePrecheck } from "./HotfixTriagePrecheck";
import { Stack, Button } from "@mui/material";
import { HotfixTriageBaseInfo } from "./HotfixTriageBaseInfo";
import { HotfixTriageEnvInfo } from "./HotfixTriageEnvInfo";
import { HotfixTriageBuildInfo } from "./HotfixTriageReleaseInfo";
import { HotfixTriageStepper } from "./HotfixTriageStepper";

import { useQuery } from "react-query";

import storage from '../common/LocalStorage';

export const HotfixTriage = ({ hotfixName }) => {
  let user = storage.getUser();

  const [hotfixBase, setHotfixBase] = React.useState()

  const updateBase = (base) => {
    setHotfixBase({ ...hotfixBase, ...base })
  }

  const [hotfixEnv, setHotfixEnv] = React.useState()

  const updateEnv = (env) => {
    setHotfixEnv({ ...hotfixEnv, ...env })
  }

  const [hotfixRelease, setHotfixRelease] = React.useState()

  const updateRelease = (release) => {
    setHotfixRelease({ ...hotfixRelease, ...release })
  }

  const update = useMutation(
    (data) => {
      return axios.post(url("hotfix"), data);
    },
    {
      onSuccess: () => {
        // TODO update hotfix info
        window.location.reload();
        // onClose();
      },
      onError: (e) => {
        console.log("error", e);
      },
    }
  );

  const { isLoading, error, data } = useQuery([`hotfix-${hotfixName}`], () => {
    console.log(url(`hotfix/${hotfixName}`))
    return fetch(url(`hotfix/${hotfixName}`)).then(async (res) => {
      return await res.json();
    });
  });

  if (isLoading) {
    return <p>Loading...</p>;
  }

  if (error) {
    return <p>Error: {error.message}</p>;
  }

  if (!hotfixBase && data != undefined) {
    var body = data.data
    setHotfixBase({
      customer: body.customer,
      oncall_url: body.oncall_url,
      creator: body.creator,
      oncall_prefix: body.oncall_prefix,
      oncall_id: body.oncall_id,
      oncall_url: body.oncall_url,
      is_debug: body.is_debug,
      platform: body.platform,
      status: body.status,
    })

    setHotfixEnv({
      base_version: body.base_version,
      artifact_archs: body.artifact_archs,
      artifact_editions: body.artifact_editions,
      artifact_types: body.artifact_types,
    })

    setHotfixRelease(
      {
        release_infos: body.release_infos.map(
          release => {
            return {
              all_prs_pushed: release.all_prs_pushed,
              repo: release.repo_full_name.split("/")[1],
              branch: release.branch,
              based_release_version: release.based_release_version,
              based_commit_sha: release.based_commit_sha,
              issues: release.issues,
              master_prs: release.master_prs,
              branch_prs: release.branch_prs,
            }
          }
        )
      }
    )
  }

  return (
    <Stack direction="column" spacing={2}>
      <Stack style={{ width: "100%" }}>
        <HotfixTriagePrecheck />
      </Stack>
      <Stack style={{ width: "100%" }}>
        <HotfixTriageStepper hotfixBase={hotfixBase} />
      </Stack>
      <Stack style={{ width: "100%" }}>
        <HotfixTriageBaseInfo onUpdate={updateBase} hotfixBase={hotfixBase} />
      </Stack>
      <Stack style={{ width: "100%" }}>
        <HotfixTriageEnvInfo onUpdate={updateEnv} hotfixEnv={hotfixEnv} />
      </Stack>
      <Stack style={{ width: "100%" }}>
        <HotfixTriageBuildInfo onUpdate={updateRelease} hotfixRelease={hotfixRelease} />
      </Stack>

      <Stack direction="row" alignItems="flex-end" spacing={2} justifyContent="flex-end">
        <Stack>
          <Button
            onClick={() => {
              update.mutate({
                ...hotfixBase,
                ...hotfixEnv,
                ...hotfixRelease,
                operator_email: user?.email,
              });
            }}
            sx={{ width: 100 }}
            variant="contained"
          >
            Update
          </Button>
        </Stack>
      </Stack >

    </Stack >
  )
}
