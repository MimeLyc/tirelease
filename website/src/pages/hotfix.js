import * as React from "react";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";
import Layout from "../layout/Layout";
import { useParams } from "react-router-dom";
import { useSearchParams } from "react-router-dom";
import { fetchActiveVersions } from "../components/issues/fetcher/fetchVersion"
import { useQuery } from "react-query";
import HotfixPlane from "../components/hotfix/HotfixPlane"

const Hotfix = () => {
  const params = useParams();
  const [searchParams, setSearchParams] = useSearchParams();

  const hotfixName = params.hotfxi === undefined ? "none" : params.hotfix;

  const hotfixQuery = useQuery(["hotfix", hotfixName], () => fetchActiveVersions({ versionName: "none" }));
  if (hotfixQuery.isLoading) {
    return (
      <div>
        <p>Loading...</p>
      </div>
    );
  }

  if (hotfixQuery.error) {
    return (
      <div>
        <p>Error: {hotfixQuery.error}</p>
      </div>
    );
  }


  return (
    <Layout>
      <Container maxWidth="xxl" sx={{ mt: 4, mb: 4 }}>
        <Paper sx={{ p: 2, display: "flex", flexDirection: "column" }}>
          <HotfixPlane />
        </Paper>
      </Container>
    </Layout>
  );
};

export default Hotfix;
