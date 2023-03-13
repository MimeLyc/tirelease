import * as React from "react";
import {
  Typography, Stack, TextField,
  Select, MenuItem, FormControl, InputLabel,
  Accordion, AccordionDetails, AccordionSummary
} from '@mui/material';

import ExpandMoreIcon from "@mui/icons-material/ExpandMore";

export const HotfixTriageQAResultInfo = ({ onUpdate, hotfixQAResult = {} }) => {
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
        <Typography sx={{ fontWeight: 'bold' }} gutterBottom variant="h8" component="div">
          {"QA Testing Result"}
        </Typography>
      </AccordionSummary>

      <AccordionDetails>
        <Stack direction="row" spacing={2}>
          <Stack>
            <FormControl>
              <InputLabel >QA Testing Result</InputLabel>
              <Select
                value={hotfixQAResult.pass_qa_test || "UNTESTED"}
                label="QA Testing Result"
                onChange={
                  (event) => {
                    hotfixQAResult.pass_qa_test = event.target.value;
                    onUpdate(hotfixQAResult);
                  }
                }
                sx={{ width: 275 }}
              >
                <MenuItem value={"UNTESTED"}>UNTESTED</MenuItem>
                <MenuItem value={"PASSED"}>PASSED</MenuItem>
                <MenuItem value={"FAILED"}>FAILED</MenuItem>
              </Select>
            </FormControl>
          </Stack>

          <Stack style={{ width: "100%" }}>
            <TextField
              label="QA Testing Report"
              multiline
              maxRows={4}
              sx={{ width: "100%" }}
              variant="outlined"
              value={hotfixQAResult.qa_test_report}
              onChange={
                (event) => {
                  hotfixQAResult.qa_test_report = event.target.value;
                  onUpdate(hotfixQAResult);
                }
              }
            />
          </Stack>
        </Stack>
      </AccordionDetails>
    </Accordion >

  )
}
