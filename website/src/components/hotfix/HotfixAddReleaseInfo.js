import * as React from "react";
import {
  Checkbox, ListItemText, Typography, Box, Chip,
  Select, MenuItem, FormControl, InputLabel
} from '@mui/material';
import Table from '@mui/material/Table';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

import { HotfixAddBuildItem } from "./HotfixAddReleaseItem"

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

// TODO, add api to fetch repos
const allRepos = ["tidb", "tiflash", "tikv", "pd", "tiflow", "tidb-binlog", "tidb-tools"];

export const HotfixAddBuildInfo = ({ hotfixes = [], onUpdate, hotfixRelease = {} }) => {
  const [repos, setRepos] = React.useState([]);
  const handleSelectRepos = (event) => {
    const {
      target: { value },
    } = event;
    const values = typeof value === 'string' ? value.split(',') : value
    setRepos(values)
  };

  return (
    <TableContainer component={Paper}>
      <Table aria-label="spanning table">
        <TableRow>
          <Typography sx={{ fontWeight: 'bold' }} gutterBottom variant="h8" component="div">
            {"Build Info"}
          </Typography>
        </TableRow>

        <TableRow>
          <FormControl sx={{ m: 1, minWidth: 200 }}>
            <InputLabel >Related Repos</InputLabel>
            <Select
              label="Related Repos"
              autoWidth
              multiple
              value={repos}
              onChange={handleSelectRepos}
              renderValue={(selected) => (
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                  {selected.map((value) => (
                    <Chip key={value} label={value} />
                  ))}
                </Box>
              )}
              MenuProps={MenuProps}
            >
              {allRepos.map((repo) => (
                <MenuItem key={repo} value={repo}>
                  <Checkbox checked={(repos || []).includes(repo)} />
                  <ListItemText primary={repo} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </TableRow>

        <TableRow>
          <HotfixAddBuildItem
            hotfixRelease={hotfixRelease}
            onUpdate={onUpdate}
            releaseRepos={repos} />
        </TableRow>
      </Table >
    </TableContainer >
  )
}
