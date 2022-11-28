import { useQuery } from "react-query";
import TiDialogTitle from "../common/TiDialogTitle";
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import Divider from '@mui/material/Divider';
import Box from '@mui/material/Box';
import * as React from "react";
import { labelFilter } from "./GridColumns"
import { renderLabel } from './renderer/Label'
import { renderIssueState } from './renderer/IssueState'
import { getPickTriageValue, renderPickTriage } from './renderer/PickTriage'
import { renderBlockRelease } from './renderer/BlockRelease'
import { renderChanged } from './renderer/ChangedItem'
import { renderComment } from './renderer/Comment'

import { renderPullRequest, getPullRequest } from './renderer/PullRequest'
import { renderAssignee } from './renderer/Assignee'
import {
  Accordion, AccordionDetails, AccordionSummary, Chip, Button, Dialog, Stack, Typography,
  Table, TableCell, TableRow
} from "@mui/material";
import Paper from "@mui/material/Paper";
import { useMutation } from "react-query";
import axios from "axios";
import { url } from "../../utils";

import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

import { DataGrid, GridToolbar } from "@mui/x-data-grid";
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import Checkbox from '@mui/material/Checkbox';

import { fetchActiveVersions } from "../../components/issues/fetcher/fetchVersion";
import { fetchSingleIssue } from "../../components/issues/fetcher/fetchIssue";

import dayjs from "dayjs";

import DialogActions from "@mui/material/DialogActions";

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

export function IssueDetail({ id, onClose, open }) {

  const [scroll, setScroll] = React.useState('paper');
  const [maxWidth, setMaxWidth] = React.useState('lg');
  const [affectVersions, setAffectVersions] = React.useState(undefined);
  const affectMutation = useMutation((newAffect) => {
    return axios.patch(url(`issue/${issueId}/affect/${newAffect.affect_version}`), newAffect);
  });
  const [issueId, setIssueId] = React.useState(id)

  const issueQuery = useQuery(
    ["single_issue", affectVersions, issueId],
    () => {
      return fetchSingleIssue({ issueId: issueId })
    },
    {
      keepPreviousData: true,
      staleTime: 500,
    }
  );
  const versionQuery = useQuery(["open", "version", "maintained"], fetchActiveVersions);

  if (issueId == undefined) {
    return <div />
  }

  if (issueQuery.isLoading) {
    return (
      <div>
        <p>Loading...</p>
      </div>
    );
  }
  if (issueQuery.isError) {
    return (
      <div>
        <p>error: {issueQuery.error}</p>
      </div>
    );
  }

  const data = issueQuery.data
  const issue = data?.data?.issue
  const masterPrs = data?.data?.master_prs
  const versionTriages = data?.data?.version_triages

  if (versionQuery.isLoading) {
    return (
      <div>
        <p>Loading...</p>
      </div>
    );
  }

  if (versionQuery.error) {
    return (
      <div>
        <p>Error: {versionQuery.error}</p>
      </div>
    );
  }

  // Get current active versions
  var minorVersions = []
  for (const version of versionQuery.data) {
    const versionName = version.name
    const minorVersion = versionName == undefined ? "none" : versionName.split(".").slice(0, 2).join(".");
    minorVersions.push(minorVersion)
  }

  if (issueId !== undefined && affectVersions == undefined) {
    setAffectVersions(versionTriages?.filter(t =>
      t.is_affect == true).filter(t => {
        var version = t.release_version
        var minorVersion = version.name.split(".").slice(0, 2).join(".");
        return minorVersions.includes(minorVersion)
      }).map((triage) => {
        return triage.release_version.name.split(".").slice(0, 2).join(".");
      }))
  }

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

    setAffectVersions(values);
  };

  const triageRows = versionTriages?.filter(triage => {
    var version = triage.release_version
    var minorVersion = version.name.split(".").slice(0, 2).join(".");
    return minorVersions.includes(minorVersion)
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
    <div>
      <Dialog
        onClose={onClose}
        open={open}
        sx={{ overflow: "visible" }}
        scroll={scroll}
        fullWidth={true}
        maxWidth={maxWidth}
        aria-labelledby="scroll-dialog-title"
        aria-describedby="scroll-dialog-description"
      >
        <Stack padding={2}>
          {(() => {
            if (issue !== undefined) {
              return <>
                <TiDialogTitle id="scroll-dialog-title" onClose={onClose}>
                  Issue Info: {issue.repo}#{issue.number}
                </TiDialogTitle>

                <List sx={{
                  width: '100%',
                  bgcolor: 'background.paper',
                }} aria-label="mailbox folders">

                  <ListItem >
                    <Typography gutterBottom variant="h6" component="div">
                      <a
                        href={issue.html_url}
                        _target="blank"
                        rel="noopener noreferrer"
                        onClick={(e) => {
                          window.open(issue.html_url);
                          e.preventDefault();
                          e.stopPropagation();
                        }}
                      >
                        {issue.title}
                      </a>
                    </Typography>
                  </ListItem>
                  <Divider />

                  <ListItem divider>
                    <Stack
                      spacing={0}
                      width="100%"
                    >
                      <Accordion defaultExpanded={true} width={"100%"}>
                        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                          <Typography sx={{ fontWeight: 'bold' }} gutterBottom variant="h8" component="div">
                            {"Detail"}
                          </Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                          <Box sx={{ width: "100%" }}>
                            <Table>
                              <TableRow>
                                <TableCell>
                                  <div>
                                    Repo:&nbsp;
                                    {issue.owner}/{issue.repo}
                                  </div>
                                </TableCell>
                                <TableCell>
                                  <div>
                                    Components:&nbsp;
                                    {issue.components.join(", ")}
                                  </div>
                                </TableCell>
                                <TableCell>
                                  <div />
                                </TableCell>
                              </TableRow>

                              <TableRow>
                                <TableCell>
                                  <div>
                                    State:&nbsp;
                                    {renderIssueState({ row: { issue: issue } })}
                                  </div>
                                </TableCell>
                                <TableCell>
                                  <div>
                                    CreateTime:&nbsp;
                                    {
                                      dayjs(issue.create_time).format(
                                        "YYYY-MM-DD HH:mm:ss"
                                      )
                                    }
                                  </div>
                                </TableCell>
                                <TableCell>
                                  {issue.close_time !== undefined && (<div>
                                    CloseTime:&nbsp;{
                                      dayjs(issue.close_time).format(
                                        "YYYY-MM-DD HH:mm:ss"
                                      )
                                    }
                                  </div>)}
                                </TableCell>
                              </TableRow>

                              <TableRow>

                                <TableCell>
                                  <div>
                                    Severity:&nbsp;
                                    {renderLabel((label) => label.name.startsWith("severity/"),
                                      (label) => label.replace("severity/", "")
                                    )({ row: { issue: issue } })}
                                  </div>

                                </TableCell>

                                <TableCell>
                                  <div>
                                    Type:&nbsp;
                                    {renderLabel((label) => label.name.startsWith("type/"),
                                      (label) => label.replace("type/", "")
                                    )({ row: { issue: issue } })}
                                  </div>

                                </TableCell>

                                <TableCell>
                                  <div>
                                    Assignees:&nbsp;
                                    {renderAssignee({ row: { issue: issue } })}
                                  </div>
                                </TableCell>
                              </TableRow>
                              <TableRow>
                                <TableCell colSpan={3}>
                                  <div>
                                    Other Labels:&nbsp;
                                    {renderLabel(labelFilter,
                                      (label) => label)({ row: { issue: issue } })
                                    }
                                  </div>

                                </TableCell>
                              </TableRow>

                              <TableRow>
                                <TableCell colSpan={3}>
                                  <div>
                                    Master PRs:&nbsp;
                                    {renderPullRequest("master")({ row: { pull_requests: masterPrs } })
                                    }
                                  </div>

                                </TableCell>
                              </TableRow>
                            </Table>
                          </Box>
                        </AccordionDetails>
                      </Accordion>
                    </Stack>
                  </ListItem>
                  <ListItem
                  >
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
                              {minorVersions.map((version) => (
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
                  </ListItem>
                </List>
              </>
            } else {
              return <div />
            }
          })()}

        </Stack>
        <DialogActions>
          <Button autoFocus onClick={onClose}>
            Close
          </Button>
        </DialogActions>
      </Dialog>
    </div >
  );
};
export default IssueDetail;
