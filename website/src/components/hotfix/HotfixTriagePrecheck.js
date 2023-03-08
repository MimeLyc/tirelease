import * as React from "react";
import { Stack, } from '@mui/material';

import { Accordion, AccordionDetails, AccordionSummary } from "@mui/material";
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';

import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
// Creators must check the box before they can submit a hotfix
export const HotfixTriagePrecheck = ({ }) => {
  const [expanded, setExpanded] = React.useState(false);

  const handleExpanded = () => {
    setExpanded(!expanded)
  };

  return (
    <Accordion
      defaultExpanded={expanded}
      expanded={expanded}
      onChange={handleExpanded}
    >
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        You have passed all precheck for applying new hotfix!` : `You should check all rules below before applying for hotfix.
      </AccordionSummary>

      <AccordionDetails>
        <Stack direction="column">
          <FormLabel>
            How about oncall master suggest for the hotfix?
          </FormLabel>
          <RadioGroup
            value={"true"}
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="Oncall master suggest to apply the hotfix." />
            <FormControlLabel
              value={"false"} control={<Radio disabled color="default" />}
              label="Other." />
          </RadioGroup>

          <FormLabel>
            Is it a bug fix?
          </FormLabel>
          <RadioGroup
            value={"true"}
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="Yes, it's a bug fix." />
            <FormControlLabel
              value={"false"} control={<Radio disabled color="default" />}
              label="No, it's a new feature or over the feature design." />
          </RadioGroup>

          <FormLabel>
            Do we have a workaround for user?
          </FormLabel>
          <RadioGroup
            value={"true"}
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="No, we don't have any workaround." />
            <FormControlLabel
              value={"false"} control={<Radio disabled color="default" />}
              label="Yes we have a workaround, but we still need a hotfix." />
          </RadioGroup>

          <FormLabel>
            Waiting for the next TiDB patch version?
            <a
              href="https://github.com/pingcap/tidb/projects/63"
              _target="blank"
              rel="noopener noreferrer"
              onClick={(e) => {
                window.open("https://github.com/pingcap/tidb/projects/63");
                e.preventDefault();
                e.stopPropagation();
              }}
            >
              TiDB release plan
            </a>
          </FormLabel>
          <RadioGroup
            value={"true"}
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="No, it's urgent." />
            <FormControlLabel
              value={"false"} control={<Radio disabled color="default" />}
              label="Yes." />
          </RadioGroup>

          <FormLabel>
            Do we have already fixed the issue and the master RP merged with test complete?
          </FormLabel>
          <RadioGroup
            value={"true"}
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="Yes, master PR merged with test complete." />
            <FormControlLabel
              value={"false"} control={<Radio disabled color="default" />}
              label="No, have not fixed yet." />
          </RadioGroup>
        </Stack >
      </AccordionDetails>
    </Accordion >
  )
}
