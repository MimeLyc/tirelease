import * as React from "react";

import { Stack } from "@mui/material";
import Button from "@mui/material/Button";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";

import { useMutation, useQueryClient } from "react-query";
import { url } from "../../utils";
import axios from "axios";
import { HotfixAddBaseInfo } from "./HotfixAddBaseInfo";
import { HotfixAddPrecheck } from "./HotfixAddPrecheck";
import { HotfixAddEnvInfo } from "./HotfixAddEnvInfo";
import { HotfixAddBuildInfo } from "./HotfixAddReleaseInfo";

import storage from '../common/LocalStorage';

export const HotfixAdd = ({ open, onClose, hotfixes }) => {
  let user = storage.getUser();

  const [hotfixPrecheck, setHotfixPrecheck] = React.useState(false)

  const [hotfixBase, setHotfixBase] = React.useState(
    {
      creator_email: user?.email,
      creator: {
        email: user?.email,
        name: user?.name
      },
      customer: "",
      operator_email: user?.email,
      oncall_prefix: "oncall",
      oncall_id: "",
      oncall_url: "",
      is_debug: false,
      platform: "",
    }
  )

  const updateBase = (base) => {
    setHotfixBase({ ...hotfixBase, ...base })
  }

  const [hotfixEnv, setHotfixEnv] = React.useState(
    {
      major: -1,
      minor: -1,
      patch: -1,
      artifact_archs: ["x86"],
      artifact_editions: ["enterprise"],
      artifact_types: ["image"],
    }
  )

  const updateEnv = (env) => {
    setHotfixEnv({ ...hotfixEnv, ...env })
  }

  const [hotfixRelease, setHotfixRelease] = React.useState({
    release_infos: [],
  })

  const updateRelease = (release) => {
    setHotfixRelease({ ...hotfixRelease, ...release })
  }

  const create = useMutation(
    (data) => {
      return axios.post(url("hotfix"), data);
    },
    {
      onSuccess: () => {
        // queryClient.invalidateQueries("versions");
        onClose();
      },
      onError: (e) => {
        console.log("error", e);
      },
    }
  );

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="lg"
      fullWidth>
      <DialogTitle>Create New Hotfix</DialogTitle>

      <Stack direction="column">
        <Stack>
          <DialogContent>
            <HotfixAddPrecheck onChange={(value) => { setHotfixPrecheck(value) }} />
          </DialogContent>
        </Stack>


        <Stack>
          {hotfixPrecheck ?
            <div>
              <DialogContent>
                <HotfixAddBaseInfo hotfixes={hotfixes} onUpdate={updateBase} hotfixBase={hotfixBase} />
              </DialogContent>

              <DialogContent>
                <HotfixAddEnvInfo hotfixes={hotfixes} onUpdate={updateEnv} hotfixEnv={hotfixEnv} />
              </DialogContent>

              <DialogContent>
                <HotfixAddBuildInfo hotfixes={hotfixes} onUpdate={updateRelease} hotfixRelease={hotfixRelease} />
              </DialogContent>

            </div>
            : <div />
          }
        </Stack>

        <Stack>
          <DialogActions>
            <Button onClick={onClose}>Close</Button>
            {hotfixPrecheck ?
              <Button
                onClick={() => {
                  if (hotfixBase.major === -1 || hotfixBase.minor === -1 || hotfixBase.patch === -1) {
                    alert(
                      "Hotfix is not complete, major, minor and patch of base Version are required"
                    );
                    return;
                  }
                  create.mutate({
                    pass_precheck: true,
                    ...hotfixBase,
                    ...hotfixEnv,
                    ...hotfixRelease,
                    status: "pending_approval",
                    operator_email: user?.email,
                    base_version: `${hotfixEnv.major}.${hotfixEnv.minor}.${hotfixEnv.patch}`,
                  });
                }}
                variant="contained"
              >
                Apply
              </Button>
              : <div />
            }
          </DialogActions>
        </Stack>
      </Stack>
    </Dialog>
  );
};
