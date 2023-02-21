import { Stack } from "@mui/material";
import { useState } from "react";
import Button from "@mui/material/Button";
import { useParams, useNavigate } from "react-router-dom";
import { HotfixAdd } from "./HotfixAdd";
import HotfixSelector from "./HotfixSelector";

const HotfixPlane = ({
}) => {

  const navigate = useNavigate();
  const params = useParams();
  const hotfix = params.hotfix === undefined ? "none" : params.hotfix;
  const [createHotfix, setCreateHotfix] = useState(false);

  return (
    <>
      <Stack spacing={1}>
        <Stack direction={"row"} justifyContent={"space-between"}>
          <HotfixSelector
            hotfixProb={hotfix}
            onChange={(v) => {
              var queryString = window.location.search
              navigate(`/home/hotfix/${v}${queryString}`, { replace: true });
            }}
          />
          <Button
            variant="contained"
            onClick={() => {
              setCreateHotfix(true);
            }}
          >
            Apply for Hotfix
          </Button>
        </Stack>
        {
          createHotfix &&
          <HotfixAdd
            open={createHotfix}
            onClose={() => {
              setCreateHotfix(false);
            }}
            hotfixes={[]}
          />
        }
      </Stack>
    </>
  );
};

export default HotfixPlane;
