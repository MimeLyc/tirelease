
import ListItemText from '@mui/material/ListItemText';

import { getPickTriageValue, renderPickTriage } from '../renderer/PickTriage'
import { renderBlockRelease } from '../renderer/BlockRelease'
import { renderChanged } from '../renderer/ChangedItem'
import { renderComment } from '../renderer/Comment'
import Divider from '@mui/material/Divider';

import { renderPullRequest, getPullRequest } from '../renderer/PullRequest'
import Box from '@mui/material/Box';
import {
  Chip, Button, Stack, Typography
} from "@mui/material";
import Paper from "@mui/material/Paper";
import { useMutation } from "react-query";
import axios from "axios";
import { url } from "../../../utils";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import Checkbox from '@mui/material/Checkbox';


const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
      width: 250,
    },
  },
};

const version = {
  field: "affected version",
  width: 80,
  headerName: "Version",
  valueGetter: (params) => params.row.minorVersion,
};

const prs = {
  field: "version prs",
  flex: 1,
  headerName: "Version PRs",
  valueGetter: (params) => getPullRequest("release-" + params.row.minorVersion)(params),
  renderCell: (params) => renderPullRequest("release-" + params.row.minorVersion)(params),
}

const block = {
  field: "block release",
  headerName: "Release Blocked",
  width: 120,
  valueGetter: (params) => params.row.version_triage.block_version_release,
  renderCell: (params) => renderBlockRelease(params),
};


const triageStatus = {
  field: "triage status",
  headerName: "Triage Status",
  width: 120,
  valueGetter: (params) => {
    return params.row.version_triage.merge_status
  },
};

const triage = {
  field: "triage",
  width: 120,
  headerName: "Pick Triage",
  valueGetter: (params) => getPickTriageValue(params.row.minorVersion)(params),
  renderCell: (params) => renderPickTriage(params.row.version, params.row.minorVersion)(params),
};

const changedItem = {
  field: "changed_item",
  headerName: "Changed Item",
  flex: 1,
  valueGetter: (params) => params.row.version_triage.changed_item,
  renderCell: (params) => renderChanged(params),
};

const comment = {
  field: "comment",
  headerName: "Comment",
  flex: 2,
  valueGetter: (params) => params.row.version_triage.comment,
  renderCell: (params) => renderComment(params),
};


export function IssueTriage({ issue, versionTriages, activeVersions, affectVersions, onAffect }) {

  const issueId = issue.issue_id
  const affectMutation = useMutation((newAffect) => {
    return axios.patch(url(`issue/${issueId}/affect/${newAffect.affect_version}`), newAffect);
  });


  const handleAffect = (event) => {
    const {
      target: { value },
    } = event;
    const values = typeof value === 'string' ? value.split(',') : value
    const addedAffection = values.filter(v => !affectVersions.includes(v))
    addedAffection.forEach(v => {
      affectMutation.mutate(
        {
          issue_id: issueId,
          affect_version: v,
          affect_result: "Yes",
        }
      )
    });

    const removedAffection = affectVersions.filter(v => !values.includes(v))
    removedAffection.forEach(v => {
      affectMutation.mutate(
        {
          issue_id: issueId,
          affect_version: v,
          affect_result: "No",
        }
      )
    });

    onAffect(values);
  };

  const triageRows = versionTriages?.filter(triage => {
    var version = triage.release_version
    var minorVersion = version.name.split(".").slice(0, 2).join(".");
    return activeVersions.includes(minorVersion)
  }).map(triage => {
    var version = triage.release_version
    var minorVersion = version.name.split(".").slice(0, 2).join(".");

    return {
      issue: issue,
      pull_requests: triage.version_prs,
      minorVersion: minorVersion,
      version: version,
      version_triage: triage,
      version_triages: versionTriages.map((t) => {
        return {
          ...t,
          version_name: t.release_version.name

        }
      }),
      id: minorVersion,
      issue_affects: [{
        affect_version: minorVersion,
        affect_result: triage.affect_result,
      }]
    }
  })

  const triageColumns = [
    version,
    prs,
    triageStatus,
    block,
    triage,
    changedItem,
    comment,
  ]


  return (
    <Paper sx={{ p: 2, width: "100%", flexDirection: "column" }}
    >
      <Stack
        spacing={2}
        width="100%"
      >
        <Typography sx={{ fontWeight: 'bold' }} gutterBottom variant="h8" component="div">
          {"Triage"}
        </Typography>

        <Stack
          direction="row"
          divider={<Divider orientation="vertical" flexItem />}
          spacing={10}
          width="100%"
        >
          <div>
            Affect Versions:&nbsp;
            <Select
              labelId="demo-multiple-checkbox-label"
              id="demo-multiple-checkbox"
              multiple
              value={affectVersions}
              onChange={handleAffect}
              // input={<OutlinedInput label="versions" />}
              renderValue={(selected) => (
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                  {selected.map((value) => (
                    <Chip key={value} label={value} />
                  ))}
                </Box>
              )}
              MenuProps={MenuProps}
            >
              {activeVersions.map((version) => (
                <MenuItem key={version} value={version}>
                  <Checkbox checked={(affectVersions || []).includes(version)} />
                  <ListItemText primary={version} />
                </MenuItem>
              ))}
            </Select>
          </div>
        </Stack>

        <Stack
          direction="row"
          divider={<Divider orientation="vertical" />}
          spacing={10}
        >
          <div style={{ height: 500, width: "100%" }}>
            <DataGrid
              density="compact"
              columns={triageColumns}
              rows={triageRows}
              components={{ Toolbar: GridToolbar }}
              showCellRightBorder={true}
              showColumnRightBorder={false}
            >
            </DataGrid>
          </div>
        </Stack>
      </Stack>
    </Paper>
  )
}

export default IssueTriage;
