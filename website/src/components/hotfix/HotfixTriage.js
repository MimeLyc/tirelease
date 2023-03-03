
import * as React from "react";
import { HotfixTriagePrecheck } from "./HotfixTriagePrecheck";
import { Stack } from "@mui/material";
import { HotfixTriageBaseInfo } from "./HotfixTriageBaseInfo";

export const HotfixTriage = ({ hotfixName }) => {
  const [hotfixBase, setHotfixBase] = React.useState(
    {
      customer: "PingCAP",
      oncall_url: "https://www.google.com",
      creator: {
        name: "Yuchao Li",
        git_login: "mimelyc",
      },
      oncall_prefix: "oncall",
      oncall_id: "12345",
      oncall_url: "",
      is_debug: false,
      platform: "OP",
      status: "UPCOMING",
    }
  )

  return (
    <Stack direction="column">
      <Stack style={{ width: "100%" }}>
        <HotfixTriagePrecheck />
      </Stack>
      <Stack style={{ width: "100%" }}>
        <HotfixTriageBaseInfo hotfixBase={hotfixBase} />

      </Stack>
    </Stack>
  )
}
