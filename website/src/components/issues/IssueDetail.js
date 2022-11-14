import { useQuery, useQueryClient } from "react-query";
import { fetchSingleIssue } from "./fetcher/fetchIssue";
import { Stack } from "@mui/material";

export function IssueDetail({ issueId }) {

  const issueQuery = useQuery(
    ["single_issue", issueId],
    () => {
      return fetchSingleIssue({ issueId: issueId })
    });

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

  const data = issueQuery.data?.data
  console.log(data)

  return (
    <Stack spacing={1}>
      <div style={{ height: 600, width: "100%" }}>
        {data.issue.html_url}
      </div>
    </Stack  >
  );
}

