import * as React from "react";

import { BaseVersionSelector } from "./BaseVersionSelector";

import {
  Typography, Box, Chip, DialogContentText,
  Select, MenuItem, FormControl, InputLabel,
  Accordion, AccordionSummary, AccordionDetails
} from '@mui/material';

import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
      width: 250,
    },
  },
};

export const HotfixTriageEnvInfo = ({ onUpdate, hotfixEnv = {} }) => {
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
          {"Environment Info"}
        </Typography>
      </AccordionSummary>


      <AccordionDetails>
        <TableContainer component={Paper}>
          <Table aria-label="spanning table">
            <TableRow>
              <TableCell colSpan={3}>
                <DialogContentText>
                  Base Version {`${hotfixEnv.base_version}`}
                </DialogContentText>
              </TableCell>

              <TableCell align="left">
                <FormControl sx={{ m: 1, minWidth: 200 }}>
                  <InputLabel >x86 or arm?</InputLabel>
                  <Select
                    label="x86 or arm?"
                    autoWidth
                    multiple
                    variant="standard"
                    disabled
                    value={hotfixEnv.artifact_archs}
                    onChange={
                      (event) => {
                        const {
                          target: { value },
                        } = event;
                        const values = typeof value === 'string' ? value.split(',') : value
                        hotfixEnv.artifact_archs = values;
                        onUpdate(hotfixEnv);
                      }
                    }
                    sx={{ width: 275 }}
                    renderValue={(selected) => (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selected.map((value) => (
                          <Chip key={value} label={value} />
                        ))}
                      </Box>
                    )}
                    MenuProps={MenuProps}
                  >
                    <MenuItem value={"x86"}>x86</MenuItem>
                    <MenuItem value={"arm"}>arm</MenuItem>
                  </Select>
                </FormControl>

              </TableCell>

              <TableCell align="left">

                <FormControl sx={{ m: 1, minWidth: 200 }}>
                  <InputLabel >Enterprise edition or community edition?</InputLabel>
                  <Select
                    label="Enterprise edition or community edition?"
                    autoWidth
                    multiple
                    value={hotfixEnv.artifact_editions}
                    variant="standard"
                    disabled
                    onChange={
                      (event) => {
                        const {
                          target: { value },
                        } = event;
                        const values = typeof value === 'string' ? value.split(',') : value
                        hotfixEnv.artifact_editions = values;
                        onUpdate(hotfixEnv);
                      }
                    }
                    sx={{ width: 275 }}
                    renderValue={(selected) => (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selected.map((value) => (
                          <Chip key={value} label={value} />
                        ))}
                      </Box>
                    )}
                    MenuProps={MenuProps}
                  >
                    <MenuItem value={"enterprise"}>enterprise</MenuItem>
                    <MenuItem value={"community"}>community</MenuItem>
                  </Select>
                </FormControl>
              </TableCell>

              <TableCell align="left">

                <FormControl sx={{ m: 1, minWidth: 200 }}>
                  <InputLabel >Delivering a TiUP offline package or image?</InputLabel>
                  <Select
                    label="Delivering a TiUP offline package or image?"
                    autoWidth
                    multiple
                    variant="standard"
                    disabled
                    value={hotfixEnv.artifact_types}
                    onChange={
                      (event) => {
                        const {
                          target: { value },
                        } = event;
                        const values = typeof value === 'string' ? value.split(',') : value
                        hotfixEnv.artifact_types = values;
                        onUpdate(hotfixEnv);
                      }
                    }
                    sx={{ width: 275 }}
                    renderValue={(selected) => (
                      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selected.map((value) => (
                          <Chip key={value} label={value} />
                        ))}
                      </Box>
                    )}
                    MenuProps={MenuProps}
                  >
                    <MenuItem value={"TiUP offline package"}>TiUP offline package</MenuItem>
                    <MenuItem value={"image"}>image</MenuItem>
                  </Select>
                </FormControl>
              </TableCell>
            </TableRow>
          </Table>
        </TableContainer>
      </AccordionDetails>
    </Accordion>
  )
}
