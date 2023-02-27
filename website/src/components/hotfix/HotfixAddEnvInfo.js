
import { BaseVersionSelector } from "./BaseVersionSelector";

import {
  Typography,
  Select, MenuItem, FormControl, InputLabel
} from '@mui/material';

import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

export const HotfixAddEnvInfo = ({ hotfixes = [], onUpdate, hotfixEnv = {} }) => {

  return (
    <TableContainer component={Paper}>
      <Table aria-label="spanning table">
        <TableRow>
          <Typography sx={{ fontWeight: 'bold' }} gutterBottom variant="h8" component="div">
            {"Environment Info"}
          </Typography>
        </TableRow>
        <TableRow>
          <TableCell colSpan={3}>
            <BaseVersionSelector
              onMajorChange={
                (major) => {
                  hotfixEnv.major = major
                  onUpdate(hotfixEnv)
                }
              }
              onMinorChange={
                (minor) => {
                  hotfixEnv.minor = minor
                  onUpdate(hotfixEnv)
                }
              }
              onPatchChange={
                (value) => {
                  hotfixEnv.patch = value
                  onUpdate(hotfixEnv)
                }
              }
            />
          </TableCell>

          {/* Automatically set owner */}
        </TableRow>

        <TableRow>
          <TableCell align="left">
            <FormControl>
              <InputLabel>x86 or arm?</InputLabel>
              <Select
                value={hotfixEnv.artifact_arch}
                label="Has fix control switch?"
                onChange={
                  (event) => {
                    hotfixEnv.artifact_arch = event.target.value;
                    onUpdate(hotfixEnv);
                  }
                }
                sx={{ width: 275 }}
              >
                <MenuItem value={"x86"}>x86</MenuItem>
                <MenuItem value={"arm"}>arm</MenuItem>
              </Select>
            </FormControl>
          </TableCell>

          <TableCell align="left">
            <FormControl>
              <InputLabel>Enterprise edition or community edition?</InputLabel>
              <Select
                value={hotfixEnv.artifact_edition}
                label="Enterprise edition or community edition?"
                onChange={
                  (event) => {
                    hotfixEnv.artifact_edition = event.target.value;
                    onUpdate(hotfixEnv);
                  }
                }
                sx={{ width: 275 }}
              >
                <MenuItem value={"enterprise"}>Enterprise</MenuItem>
                <MenuItem value={"community"}>Community</MenuItem>
              </Select>
            </FormControl>
          </TableCell>

          <TableCell align="left">
            <FormControl>
              <InputLabel>Delivering a TiUP offline package or image?</InputLabel>
              <Select
                value={hotfixEnv.artifact_type}
                onChange={
                  (event) => {
                    hotfixEnv.artifact_type = event.target.value;
                    onUpdate(hotfixEnv);
                  }
                }
                label="Delivering a TiUP offline package or image?"
                sx={{ width: 275 }}
              >
                <MenuItem value={"tiup"}>TiUP offline package</MenuItem>
                <MenuItem value={"image"}>Image</MenuItem>
              </Select>
            </FormControl>

          </TableCell>

        </TableRow>



      </Table>
    </TableContainer>
  )
}
