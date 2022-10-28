import { useState } from "react";
import { useQuery, useQueryClient } from "react-query";
import { fetchIssue } from "./fetcher/fetchIssue";
import { Stack } from "@mui/material";

export function IssueDetail({ issueId }) {


  return (
    <Stack spacing={1}>
      <div style={{ height: 600, width: "100%" }}>
        {issueId}
      </div>
    </Stack  >
  );
}

