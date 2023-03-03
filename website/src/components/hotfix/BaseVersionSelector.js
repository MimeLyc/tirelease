import * as React from "react";
import { useQuery } from "react-query";
import Typography from "@mui/material/Typography";
import {
  Stack
} from '@mui/material';

import { url } from "../../utils";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import FormControl from "@mui/material/FormControl";
import InputLabel from "@mui/material/InputLabel";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";


function aggregate(values) {
  const dedup = [...new Set(values)];
  const sorted = dedup.sort();
  return { sorted };
}

function getMajors(versions) {
  const majors = [];
  for (const version of versions) {
    const [major] = version.split(".");
    majors.push(parseInt(major));
  }
  return aggregate(majors);
}

function getMinors(versions, targetMajor) {
  const minors = [];
  for (const version of versions) {
    const [major, minor] = version.split(".");
    if (parseInt(major) === targetMajor) {
      minors.push(parseInt(minor));
    }
  }
  return aggregate(minors);
}

function getPatches(versions, targetMajor, targetMinor) {
  const patches = [];
  for (const version of versions) {
    const [major, minor, patch] = version.split(".");
    if (parseInt(major) === targetMajor && parseInt(minor) === targetMinor) {
      patches.push(parseInt(patch));
    }
  }
  return aggregate(patches);
}

export const BaseVersionSelector = ({ onMajorChange, onMinorChange, onPatchChange }) => {
  const { isLoading, error, data } = useQuery("versions", () => {
    return fetch(url("version")).then(async (res) => {
      return await res.json();
    });
  });

  var versions = []
  if (data) {
    console.log("version", data);
    versions = data.data.map((version) => version.name)
  }


  // const queryClient = useQueryClient();
  const [major, setMajor] = React.useState(-1);
  const [minor, setMinor] = React.useState(-1);
  const [patch, setPatch] = React.useState(-1);

  const majorData = getMajors(versions);
  const [minorData, setMinorData] = React.useState({ sorted: [] });
  const [patchData, setPatchData] = React.useState({ sorted: [] });

  return (
    < Stack direction="column" spacing={2} alignItems="top" >
      <DialogContentText>
        Base Version {major === -1 ? "[major]" : major}.
        {minor === -1 ? "[minor]" : minor}.
        {patch === -1 ? "[patch]" : patch}
      </DialogContentText>
      {isLoading && <p>Loading...</p>}
      <Stack direction="row" spacing={2} alignItems="flex-end">
        <FormControl fullWidth>
          <InputLabel id="create-version">Major *</InputLabel>
          <Select
            labelId="create-version"
            id="create-version-select"
            value={major === -1 ? "" : major}
            label="Version"
            onChange={(e) => {
              setMajor(e.target.value);
              setMinorData(getMinors(versions, e.target.value));
              onMajorChange(e.target.value);
            }}
            autoWidth
          >
            {majorData.sorted.map((v) => (
              <MenuItem value={v}>{v}</MenuItem>
            ))}
          </Select>
        </FormControl>
        <Typography fontSize={"2em"}>.</Typography>
        <FormControl fullWidth>
          <InputLabel id="create-version">Minor *</InputLabel>
          <Select
            labelId="create-version"
            id="create-version-select"
            value={minor === -1 ? "" : minor}
            label="Version"
            onChange={(e) => {
              setMinor(e.target.value);
              setPatchData(getPatches(versions, major, e.target.value));
              onMinorChange(e.target.value);
            }}
            autoWidth
            disabled={major === -1}
          >
            {minorData.sorted.map((v) => (
              <MenuItem value={v}>{v}</MenuItem>
            ))}
          </Select>
        </FormControl>
        <Typography fontSize={"2em"}>.</Typography>
        <FormControl fullWidth>
          <InputLabel id="create-version">Patch *</InputLabel>
          <Select
            labelId="create-version"
            id="create-version-select"
            value={patch === -1 ? "" : patch}
            label="Version"
            onChange={(e) => {
              setPatch(e.target.value);
              onPatchChange(e.target.value);
            }}
            autoWidth
            disabled={minor === -1}
          >
            {patchData.sorted.map((v) => (
              <MenuItem value={v}>{v}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </Stack>
    </Stack >
  )
}


export default BaseVersionSelector;
