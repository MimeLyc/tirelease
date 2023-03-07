import * as React from "react";
import {
  Stack, TextField, Chip, Divider, Paper, Typography, Link
} from '@mui/material';

import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';

import Autocomplete from '@mui/material/Autocomplete';

function parseIssueUrl(url) {
  const pattern = /https:\/\/github\.com\/([\w-]+)\/([\w-]+)\/issues\/(\d+)/;
  const matches = pattern.exec(url);

  if (!matches) {
    return null;
  }

  const ownerName = matches[1];
  const repoName = matches[2];
  const issueNumber = matches[3];

  return {
    "owner": ownerName,
    "repo": repoName,
    "number": parseInt(issueNumber),
    "html_url": url
  }
}

function parsePRUrl(url) {
  const pattern = /^https:\/\/github\.com\/([^/]+)\/([^/]+)\/pull\/(\d+)/;
  const matches = pattern.exec(url);

  if (!matches) {
    return null;
  }

  const ownerName = matches[1];
  const repoName = matches[2];
  const prNumber = matches[3];

  return {
    "owner": ownerName,
    "repo": repoName,
    "number": parseInt(prNumber),
    "html_url": url
  }
}

export const HotfixTriageReleaseItem = ({ hotfixRelease = {}, onUpdate, releaseRepos = [] }) => {
  const [repoInfos, setRepoInfos] = React.useState({})
  const [issueAutoKey, setIssueAutoKey] = React.useState(0)
  const [prAutoKey, setPrAutoKey] = React.useState(0)
  const [cpAutoKey, setCpAutoKey] = React.useState(0)

  var repoChanged = false
  releaseRepos.forEach((repo) => {
    if (hotfixRelease.release_infos.filter((info) => info.repo == repo).length == 0) {
      repoChanged = true
      hotfixRelease.release_infos.push({
        repo: repo
      })
    }
  })
  hotfixRelease.release_infos = hotfixRelease.release_infos.filter((info) => {
    if (!releaseRepos.includes(info.repo)) {
      repoChanged = true
      return false
    }
    return true
  })
  if (repoChanged) {
    // onUpdate(hotfixRelease)
  }

  const handleUpdate = () => {
    var repos = []
    for (var k in repoInfos) {
      var v = repoInfos[k]
      repos.push({
        repo: k,
        branch: v.branch,
        based_release_version: v.releaseVersion,
        based_commit_sha: v.releaseCommit,
        issues: v.issues,
        master_prs: v.masterPrs,
        branch_prs: v.branchPrs
      })
    }

    hotfixRelease.release_infos = repos
    onUpdate(hotfixRelease)
  }

  return (
    <Stack direction="column" >
      {
        releaseRepos.map(
          (repo, index) => {
            repoInfos[repo] = hotfixRelease
              .release_infos
              .filter((info) => info.repo == repo).map((info) => {
                return {
                  owner: info.owner,
                  repo: info.repo,
                  releaseVersion: info.based_release_version,
                  releaseCommit: info.based_commit_sha,
                  issues: info.issues,
                  branch: info.branch,
                  masterPrs: info.master_prs,
                  branchPrs: info.branch_prs
                }
              })[0] || {}
            return (
              <Paper>
                <Divider orientation="horizontal" textAlign="left" >{repo.toUpperCase()} </Divider>
                <Table>
                  <TableRow>
                    <TableCell align="left" colSpan={2}>
                      <TextField
                        label="Based Release Version"
                        variant="standard"
                        disabled
                        sx={{ width: "100%" }}
                        value={repoInfos[repo].releaseVersion}
                        onChange={(event) => {
                          repoInfos[repo].releaseVersion = event.target.value;
                          setRepoInfos(repoInfos)
                          handleUpdate();
                        }}
                      />
                    </TableCell>

                    <TableCell align="left" colSpan={2}>
                      <TextField
                        label="Based Git Commit Hash"
                        value={repoInfos[repo].releaseCommit}
                        variant="standard"
                        disabled
                        sx={{ width: "100%" }}
                        onChange={(event) => {
                          repoInfos[repo].releaseCommit = event.target.value;
                          setRepoInfos(repoInfos)
                          handleUpdate();
                        }}
                      />
                    </TableCell>

                    <TableCell colSpan={2} align="left">
                      <Autocomplete
                        key={issueAutoKey}
                        multiple
                        defaultValue={null}
                        options={[]}
                        variant="standard"
                        disabled
                        freeSolo
                        sx={{ width: "100%" }}
                        value={repoInfos[repo]?.issues || []}
                        // newValue will be the array conaining all inputs
                        onChange={(event, newValue) => {
                          repoInfos[repo] = repoInfos[repo] || { "issues": [] }
                          repoInfos[repo].issues = newValue.map((issue) => {
                            if (typeof issue === 'string') {
                              return parseIssueUrl(issue);
                            } else {
                              return issue
                            }
                          })
                          setRepoInfos(
                            repoInfos,
                          );
                          setIssueAutoKey((prev) => prev + 1)
                          handleUpdate();
                        }}

                        renderTags={(value, getTagProps) =>
                          value.map((issue, index) => {
                            if (typeof issue === 'string') {
                              issue = parseIssueUrl(issue);
                            }
                            return <Chip
                              variant="outlined"
                              label={"#" + issue.number}
                              onClick={() => {
                                window.open(issue.html_url);
                              }}
                              size="small"
                              {...getTagProps({ index })} />
                          }
                          )
                        }
                        renderInput={(params) => (
                          <TextField
                            {...params}
                            variant="standard"
                            disabled
                            label="Please input issue url..."
                            placeholder="issue url..."
                          />
                        )}
                      />

                    </TableCell>

                    <TableCell colSpan={2} align="left">
                      <Autocomplete
                        key={prAutoKey}
                        multiple
                        defaultValue={null}
                        options={[]}
                        freeSolo
                        sx={{ width: "100%" }}
                        value={repoInfos[repo]?.masterPrs || []}
                        // newValue will be the array conaining all inputs
                        onChange={(event, newValue) => {
                          repoInfos[repo] = repoInfos[repo] || { "masterPrs": [] }
                          repoInfos[repo].masterPrs = newValue.map((pr) => {
                            if (typeof pr === 'string') {
                              return parsePRUrl(pr);
                            } else {
                              return pr
                            }
                          })
                          setRepoInfos(
                            repoInfos,
                          );
                          setPrAutoKey((prev) => prev + 1)
                          handleUpdate();
                        }}
                        renderTags={(value, getTagProps) =>
                          value.map((pr, index) => {
                            if (typeof pr === 'string') {
                              pr = parsePRUrl(pr);
                            }
                            return <Chip
                              variant="outlined"
                              label={"#" + pr.number}
                              onClick={() => {
                                window.open(pr.html_url);
                              }}
                              size="small"
                              {...getTagProps({ index })} />
                          }
                          )
                        }
                        renderInput={(params) => (
                          <TextField
                            {...params}
                            label="Please input master-pr url..."
                            placeholder="master-pr url..."
                          />
                        )}
                      />
                    </TableCell>

                  </TableRow>

                  <TableRow>
                    <TableCell align="left" colSpan={1}>
                      <TextField
                        disabled
                        label="Hotfix Branch"
                        variant="standard"
                        sx={{ width: 150 }}
                        defaultValue={repoInfos[repo].branch}
                      />
                    </TableCell>

                    <TableCell colSpan={6} align="left">
                      <Autocomplete
                        key={cpAutoKey}
                        multiple
                        defaultValue={null}
                        options={[]}
                        freeSolo
                        sx={{ width: "100%" }}
                        value={repoInfos[repo]?.branchPrs || []}
                        onChange={(event, newValue) => {
                          repoInfos[repo] = repoInfos[repo] || { "branchPrs": [] }
                          repoInfos[repo].branchPrs = newValue.map((pr) => {
                            if (typeof pr === 'string') {
                              return parsePRUrl(pr);
                            } else {
                              return pr
                            }
                          })
                          setRepoInfos(
                            repoInfos,
                          );
                          setCpAutoKey((prev) => prev + 1)
                          handleUpdate();
                        }}
                        renderTags={(value, getTagProps) =>
                          value.map((pr, index) => {
                            if (typeof pr === 'string') {
                              pr = parsePRUrl(pr);
                            }
                            return <Chip
                              variant="outlined"
                              label={"#" + pr.number}
                              onClick={() => {
                                window.open(pr.html_url);
                              }}
                              size="small"
                              {...getTagProps({ index })} />
                          }
                          )
                        }
                        renderInput={(params) => (
                          <TextField
                            {...params}
                            label="Please input cherry-pick-pr url..."
                            placeholder="cherry-pick-pr url url..."
                          />
                        )}
                      />
                    </TableCell>
                  </TableRow>

                  <TableRow>
                    <TableCell colSpan={1}>
                      <Stack direction="column">
                        <Stack >
                          <Typography variant="caption" gutterBottom>
                            {"Building Status"}
                          </Typography>
                        </Stack>

                        <Stack >
                          <Chip
                            color="primary"
                            sx={{
                              width: 150,
                            }}
                            label="Building" />
                        </Stack>
                      </Stack>
                    </TableCell>

                    <TableCell colSpan={6}>
                      <TextField
                        disabled
                        variant="standard"
                        label="Build Artifacts"
                        sx={{ width: "100%" }}
                        InputProps={{
                          startAdornment: (
                            <Link
                              sx={{ width: "100%" }}
                              target="_blank">
                              广告位招租
                            </Link>
                          ),
                        }}
                      />

                    </TableCell>

                  </TableRow>


                </Table>
              </Paper>
            )
          }
        )
      }
    </Stack >
  )
}
