import { useQuery } from "react-query";
import TiDialogTitle from "../../common/TiDialogTitle";
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import Divider from '@mui/material/Divider';
import * as React from "react";
import { Button, Dialog, Stack, Typography } from "@mui/material";

import { fetchSingleIssue } from "../../../components/issues/fetcher/fetchIssue";
import DialogActions from "@mui/material/DialogActions";
import IssueMeta from "./IssueMeta";
import IssueTriage from "./IssueTriage";
import { fetchActiveVersions } from "../../../components/issues/fetcher/fetchVersion";

export function IssueDetail({ id, onClose, open }) {

  const [scroll, setScroll] = React.useState('paper');
  const [maxWidth, setMaxWidth] = React.useState('lg');
  const [issueId, setIssueId] = React.useState(id)
  const [affectVersions, setAffectVersions] = React.useState(undefined);
  const [triages, setTriages] = React.useState(undefined);
  const [trigger, setTrigger] = React.useState(0);

  const issueQuery = useQuery(
    ["single_issue", affectVersions, issueId, trigger],
    () => {
      return fetchSingleIssue({ issueId: issueId })
    },
    {
      onSuccess: (data) => {
        setTriages(undefined)
      },
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

  const data = issueQuery.data
  const issue = data?.data?.issue
  const masterPrs = data?.data?.master_prs

  if (triages == undefined || triages != data?.data?.version_triages) {
    setTriages(data.data.version_triages)
  }

  if (triages !== undefined && affectVersions == undefined) {
    setAffectVersions(triages?.filter(t =>
      t.is_affect == true).filter(t => {
        var version = t.release_version
        var minorVersion = version.name.split(".").slice(0, 2).join(".");
        return minorVersions.includes(minorVersion)
      }).map((triage) => {
        return triage.release_version.name.split(".").slice(0, 2).join(".");
      }))
  }

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

                  <ListItem>
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
                    <IssueMeta issue={issue} masterPrs={masterPrs} />
                  </ListItem>
                  <ListItem>
                    <IssueTriage
                      issue={issue}
                      versionTriages={triages}
                      activeVersions={minorVersions}
                      affectVersions={affectVersions}
                      onAffect={(values) => { setAffectVersions(values) }}
                      onBatchApprove={() => {
                        setTrigger(trigger + 1)
                      }}
                    />
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
