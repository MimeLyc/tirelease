import * as React from "react";
import Container from "@mui/material/Container";
import Layout from "../layout/Layout";

import Typography from '@mui/material/Typography';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import Link from '@mui/material/Link';
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import Box from "@mui/material/Box";

import { useQuery } from "react-query";
import { IssueGrid } from "../components/issues/IssueGrid";
import Columns from "../components/issues/GridColumns";
import { fetchActiveVersions } from "../components/issues/fetcher/fetchVersion";
import { Filters } from "../components/issues/filter/FilterDialog";
import { useSearchParams } from "react-router-dom";
import { IssueDetail } from "../components/issues/IssueDetail"

const SingleIssue = () => {
  // Duplicate with VersionTriage plane.
  // Because the "useSearchParams" must be used in component function.
  const [searchParams, setSearchParams] = useSearchParams();
  const issueNum = searchParams.get("issue_num")
  const issueId = searchParams.get("issue_id")
  const title = issueNum ? issueNum : issueId

  return (
    <Layout>
      <Container maxWidth="xxl" sx={{ mt: 4, mb: 4 }}>
        <Breadcrumbs aria-label="breadcrumb">
          <Link underline="hover" color="inherit" href="${redirectFromUrl}">
            {redirectFromUrl}
          </Link>
          <Typography color="text.primary">{title}</Typography>
        </Breadcrumbs>
        <Box sx={{ width: "100%" }}>
          <IssueDetail issueId={issueId}></IssueDetail>
        </Box>

      </Container>
    </Layout>
  );
};

export default SingleIssue;
