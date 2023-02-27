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

import storage from '../common/LocalStorage';

export const HotfixAdd = ({ open, onClose, hotfixes }) => {
  let user = storage.getUser();

  const [hotfixPrecheck, setHotfixPrecheck] = React.useState(false)

  const [hotfixBase, setHotfixBase] = React.useState(
    {
      submitor: user?.email,
      oncall_prefix: "oncall",
      oncall_id: -1,
      oncall_url: "",
      is_debug: false,
      is_on_hotfix: false,
      has_control_switch: true,
      rollback_method: "",
      trigger_reason: "",
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
      artifact_arch: "x86",
      artifact_edition: "enterprise",
      artifact_type: "image",
    }
  )

  const updateEnv = (env) => {
    setHotfixEnv({ ...hotfixEnv, ...env })
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

            </div>
            : <div />
          }
        </Stack>

        <Stack>
          <DialogActions>
            <Button onClick={onClose}>Close</Button>
            <Button
              onClick={() => {
                if (hotfixBase.major === -1 || hotfixBase.minor === -1 || hotfixBase.patch === -1) {
                  alert(
                    "Hotfix is not complete, major, minor and patch of base Version are required"
                  );
                  return;
                }
                create.mutate({
                  name: `date-${hotfixBase.major}.${hotfixBase.minor}.${hotfixBase.patch}-customer`,
                  base_version: `${hotfixBase.major}.${hotfixBase.minor}.${hotfixBase.patch}`,
                  status: "pending_approval",
                  creator_email: user?.email,
                  operator_email: user?.email,
                });
              }}
              variant="contained"
            >
              Apply
            </Button>
          </DialogActions>
        </Stack>
      </Stack>
    </Dialog>
  );
};
