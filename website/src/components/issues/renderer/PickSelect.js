import * as React from "react";
import FormControl from "@mui/material/FormControl";
import MenuItem from "@mui/material/MenuItem";
import Select from "@mui/material/Select";
import { useMutation } from "react-query";
import axios from "axios";
import { url } from "../../../utils";
import { mapPickStatusToBackend } from "./mapper"
import { fetchVersionByOption } from "../fetcher/fetchVersion";

export default function PickSelect({
  id,
  version = "master",
  patch = "master",
  pick = "unknown",
  onChange = () => { },
}) {
  const mutation = useMutation((newAffect) => {
    return axios.patch(url(`issue/${id}/cherrypick/${version}`), newAffect);
  });
  const [affects, setAffects] = React.useState(pick);

  const [isVersionFrozen, setIsVersionFrozen] = React.useState(false);
  // Query single version for the approved selector status.
  var versionOption = composeVersionOption(version)
  React.useEffect(() => {
    if (affects != "approved") {
      fetchVersionByOption({ page: 1, perPage: 1, option: versionOption })
        .then(
          (data) => {
            var versionResp = data.data;
            var hasResult = versionResp.length > 0;
            if (hasResult && versionResp[0].status == "frozen") {
              setIsVersionFrozen(true)
            } else {
              setIsVersionFrozen(false)
            }
          }
        )
    }
  }, []);

  const handleChange = (event) => {
    mutation.mutate({
      issue_id: id,
      version_name: version,
      triage_result: mapPickStatusToBackend(event.target.value),
      updated_vars: ["triage_result"]
    });
    onChange(event.target.value);
    setAffects(event.target.value);
  };

  return (
    <>
      {mutation.isLoading ? (
        <p>Updating...</p>
      ) : (
        <>
          {mutation.isError ? (
            <div>An error occurred: {mutation.error.message}</div>
          ) : null}
          <FormControl variant="standard" sx={{ m: 1, minWidth: 120 }}>
            <Select
              id="demo-simple-select-standard"
              value={affects}
              onChange={handleChange}
              label="Affection"
              disabled={pick.startsWith("released")}
            >
              <MenuItem value={"N/A"} disabled={true}>-</MenuItem>
              <MenuItem value={"unknown"}>unknown</MenuItem>
              {
                isVersionFrozen && affects != "approved" ? (
                  <MenuItem value={"approved(frozen)"}>
                    <div style={{ color: "CornflowerBlue", fontWeight: "bold" }}>
                      approved(frozen)
                    </div>
                  </MenuItem>
                ) : (
                  <MenuItem value={"approved"}>
                    <div style={{ color: "green", fontWeight: "bold" }}>
                      approved
                    </div>
                  </MenuItem>
                )
              }
              <MenuItem value={"later"}>later</MenuItem>
              <MenuItem value={"won't fix"}>won't fix</MenuItem>
              <MenuItem value={"released"} disabled={true}>
                released in {patch}
              </MenuItem>
            </Select>
          </FormControl>
        </>
      )}
    </>
  );
}

const PATCH_PATTERN = /\d+\.\d+\.\d+/
const MINOR_PATTERN = /\d+\.\d+/

function composeVersionOption(version) {
  var option = {}

  const versionItems = version.split(".")
  option["major"] = versionItems[0]
  option["minor"] = versionItems[1]

  if (PATCH_PATTERN.exec(version)) {
    option["short_type"] = "%d.%d.%d"
    option["patch"] = versionItems[2]
  } else if (MINOR_PATTERN.exec(version)) {
    option["short_type"] = "%d.%d"
    option["status_list"] = ["upcoming", "frozen"]
  }

  option["order_by"] = ["update_time"]

  return option
}
