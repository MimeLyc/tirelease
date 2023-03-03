import * as React from "react";
import { Stack, } from '@mui/material';

import { Accordion, AccordionDetails, AccordionSummary } from "@mui/material";
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';
import FormHelperText from '@mui/material/FormHelperText';

import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
// Creators must check the box before they can submit a hotfix
export const HotfixAddPrecheck = ({ onChange }) => {
  const [masterCheck, setMasterCheck] = React.useState("");
  const [bugCheck, setBugCheck] = React.useState("");
  const [workaroundCheck, setWorkaroundCheck] = React.useState("");
  const [waitingCheck, setWaitingCheck] = React.useState("");
  const [prCheck, setPrCheck] = React.useState("");

  const [expanded, setExpanded] = React.useState(false);

  const allChecked = masterCheck == "true" && bugCheck == "true" && workaroundCheck == "true" && waitingCheck == "true" && prCheck == "true"
  onChange(allChecked)

  const handleExpanded = () => {
    setExpanded(!expanded)
  };

  return (
    <Accordion
      defaultExpanded={!allChecked}
      expanded={!allChecked || expanded}
      onChange={handleExpanded}
    >
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        {
          allChecked ? `You have passed all precheck for applying new hotfix!` : `You should check all rules below before applying for hotfix.`
        }
      </AccordionSummary>

      <AccordionDetails>
        <Stack direction="column">

          <FormLabel>
            How about oncall master suggest for the hotfix?
          </FormLabel>
          <RadioGroup
            value={masterCheck}
            onChange={(event) => {
              setMasterCheck(event.target.value);
            }
            }
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="Oncall master suggest to apply the hotfix." />
            <FormControlLabel
              value={"false"} control={<Radio color="default" />}
              label="Other." />
            {masterCheck == "false" ? <FormHelperText>You could feel free to contact with oncall for support.</FormHelperText> : <div />}
          </RadioGroup>

          <FormLabel>
            Is it a bug fix?
          </FormLabel>
          <RadioGroup
            value={bugCheck}
            onChange={(event) => {
              setBugCheck(event.target.value);
            }
            }
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="Yes, it's a bug fix." />
            <FormControlLabel
              value={"false"} control={<Radio color="default" />}
              label="No, it's a new feature or over the feature design." />
          </RadioGroup>

          <FormLabel>
            Do we have a workaround for user?
          </FormLabel>
          <RadioGroup
            value={workaroundCheck}
            onChange={(event) => {
              setWorkaroundCheck(event.target.value);
            }
            }
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="No, we don't have any workaround." />
            <FormControlLabel
              value={"false"} control={<Radio color="default" />}
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
            value={waitingCheck}
            onChange={(event) => {
              setWaitingCheck(event.target.value);
            }
            }
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="No, it's urgent." />
            <FormControlLabel
              value={"false"} control={<Radio color="default" />}
              label="Yes." />
            {waitingCheck == "false" ? <FormHelperText>You can feel free to write it here for patch release.
              <a
                href="https://pingcap.feishu.cn/sheets/shtcnBhDkzwOC9FBzY94Qo8jzmf"
                _target="blank"
                rel="noopener noreferrer"
                onClick={(e) => {
                  window.open("https://pingcap.feishu.cn/sheets/shtcnBhDkzwOC9FBzY94Qo8jzmf");
                  e.preventDefault();
                  e.stopPropagation();
                }}
              >
                待排期 Patch 版本情况
              </a>

            </FormHelperText> : <div />}

          </RadioGroup>


          <FormLabel>
            Do we have already fixed the issue and the master RP merged with test complete?
          </FormLabel>
          <RadioGroup
            value={prCheck}
            onChange={(event) => {
              setPrCheck(event.target.value);
            }
            }
          >
            <FormControlLabel
              value={"true"} control={<Radio />}
              label="Yes, master PR merged with test complete." />
            <FormControlLabel
              value={"false"} control={<Radio color="default" />}
              label="No, have not fixed yet." />

            {prCheck == "false" ? <FormHelperText>You could feel free to contact with oncall for support.</FormHelperText> : <div />}
          </RadioGroup>
        </Stack >

      </AccordionDetails>
    </Accordion >
  )
}
