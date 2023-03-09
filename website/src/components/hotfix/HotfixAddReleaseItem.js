import * as React from "react";
import NativeSelect from '@mui/material/NativeSelect';
import {
  Stack, TextField, Typography, Chip, Divider,
  InputAdornment, Select, MenuItem, FormControl, InputLabel
} from '@mui/material';


import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
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

export const HotfixAddBuildItem = ({ hotfixRelease = {}, onUpdate, releaseRepos = [] }) => {
  const [repoInfos, setRepoInfos] = React.useState({})
  const [issueAutoKey, setIssueAutoKey] = React.useState(0)
  const [prAutoKey, setPrAutoKey] = React.useState(0)

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
    onUpdate(hotfixRelease)
  }

  const handleUpdate = () => {
    var repos = []
    for (var k in repoInfos) {
      var v = repoInfos[k]
      repos.push({
        repo: k,
        // git_ref_type: v.gitRefType,
        based_release_version: v.releaseVersion,
        based_commit_sha: v.releaseCommit,
        issues: v.issues,
        master_prs: v.masterPrs
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
            repoInfos[repo] = repoInfos[repo] || {}

            return (
              <Stack >
                <Divider orientation="horizontal" textAlign="left" >{repo.toUpperCase()} </Divider>
                <Table>
                  <TableRow>
                    <TableCell align="left" colSpan={1}>
                      <TextField
                        disabled
                        label="Repo"
                        sx={{ width: 150 }}
                        defaultValue={repo}
                      />
                    </TableCell>

                    <TableCell align="left" colSpan={1}>
                      <TextField
                        label="Based Release Version"
                        // sx={{ width: 150 }}
                        sx={{ width: "100%" }}
                        onChange={(event) => {
                          repoInfos[repo].releaseVersion = event.target.value;
                          setRepoInfos(repoInfos)
                          handleUpdate();
                        }}
                      // InputProps={{
                      //   startAdornment:
                      //     <InputAdornment position="start">
                      //       <FormControl fullWidth>
                      //         <NativeSelect
                      //           sx={{ width: 100 }}
                      //           defaultValue={""}
                      //           value={repoInfos[repo].gitRefType || ""}
                      //           onChange={(event) => {
                      //             repoInfos[repo].gitRefType = event.target.value;
                      //             setRepoInfos(repoInfos)
                      //             handleUpdate();
                      //           }}
                      //           inputProps={{
                      //             name: 'git_ref_type',
                      //             id: 'uncontrolled-native',
                      //           }}
                      //         >
                      //           <option value={"branch"}>branchs</option>
                      //           <option value={"tag"}>tags</option>
                      //         </NativeSelect>
                      //       </FormControl>
                      //       {"/"}
                      //     </InputAdornment>,
                      // }}
                      />
                    </TableCell>

                    <TableCell align="left" colSpan={2}>
                      <TextField
                        label="Based Git Commit Hash"
                        sx={{ width: 300 }}
                        onChange={(event) => {
                          repoInfos[repo].releaseCommit = event.target.value;
                          setRepoInfos(repoInfos)
                          handleUpdate();
                        }}

                      />
                    </TableCell>

                  </TableRow>

                  <TableRow>
                    <TableCell colSpan={2} align="left">
                      <Autocomplete
                        key={issueAutoKey}
                        multiple
                        defaultValue={null}
                        options={[]}
                        freeSolo
                        sx={{ width: 550 }}
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
                        sx={{ width: 550 }}
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
                </Table>
              </Stack>
            )
          }
        )
      }
    </Stack >
  )
}
