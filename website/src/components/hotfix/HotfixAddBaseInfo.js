import NativeSelect from '@mui/material/NativeSelect';
import {
  TextField, Typography,
  InputAdornment, Select, MenuItem, FormControl, InputLabel
} from '@mui/material';

import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import Autocomplete from '@mui/material/Autocomplete';

import storage from '../common/LocalStorage';

export const HotfixAddBaseInfo = ({ hotfixes = [], onUpdate, hotfixBase = {} }) => {
  let user = storage.getUser();

  return (
    <TableContainer component={Paper}>
      <Table aria-label="spanning table">
        <TableRow>
          <Typography sx={{ fontWeight: 'bold' }} gutterBottom variant="h8" component="div">
            {"Basic Info"}
          </Typography>
        </TableRow>

        <TableRow>
          <TableCell colSpan={1}>
            <Autocomplete
              id="free-solo-demo"
              freeSolo
              sx={{ width: 275 }}
              options={hotfixes.map((option) => option.customer)}
              onChange={(event, newValue) => {
                hotfixBase.customer = newValue;
                onUpdate(hotfixBase);
              }}
              renderInput={(params) => <TextField  {...params} label="Customer" />}
            />
          </TableCell>

          <TableCell colSpan={1}>
            <TextField
              disabled
              id="outlined-disabled"
              label="Submitor"
              sx={{ width: 275 }}
              defaultValue={user ? `${user.name}(${user.git_login})` : "You haven't logged in yet."}
            />
          </TableCell>
        </TableRow>

        <TableRow>
          <TableCell align="left">
            <TextField
              label="Oncall ID"
              sx={{ width: 275 }}
              onChange={(event) => {
                hotfixBase.oncall_id = event.target.value;
                onUpdate(hotfixBase);
              }}
              InputProps={{
                startAdornment:
                  <InputAdornment position="start">
                    <FormControl fullWidth>
                      <NativeSelect
                        sx={{ width: 100 }}
                        defaultValue={hotfixBase.oncall_prefix}
                        value={hotfixBase.oncall_prefix}
                        onChange={(event) => {
                          hotfixBase.oncall_prefix = event.target.value;
                          onUpdate(hotfixBase);
                        }}
                        inputProps={{
                          name: 'oncall_platform',
                          id: 'uncontrolled-native',
                        }}
                      >
                        <option value={"oncall"}>ONCALL</option>
                        <option value={"ticket"}>TICKET</option>
                      </NativeSelect>
                    </FormControl>
                    {" - "}
                  </InputAdornment>,
              }}
            />
          </TableCell>
          <TableCell colSpan={3} align="left">
            <TextField
              sx={{ width: "100%" }}
              label="Oncall URL"
              onChange={
                (event) => {
                  hotfixBase.oncall_url = event.target.value;
                  onUpdate(hotfixBase);
                }
              }
            />
          </TableCell>
        </TableRow>
        <TableRow>
          <TableCell align="left">
            <FormControl
            >
              <InputLabel id="demo-simple-select-label">Is for debug?</InputLabel>
              <Select
                value={hotfixBase.is_debug ? "Yes" : "No"}
                label="Is it a debug hotfix?"
                onChange={
                  (event) => {
                    hotfixBase.is_debug = event.target.value == "Yes";
                    onUpdate(hotfixBase);
                  }
                }
                sx={{ width: 275 }}
              >
                <MenuItem value={"Yes"}>Yes</MenuItem>
                <MenuItem value={"No"}>No</MenuItem>
              </Select>
            </FormControl>
          </TableCell>

          <TableCell align="left">
            <FormControl>
              <InputLabel >OP or on TiDB Cloud?</InputLabel>
              <Select
                value={hotfixBase.platform}
                label="Whether the user has applied any hotfix before?"
                onChange={
                  (event) => {
                    hotfixBase.platform = event.target.value;
                    onUpdate(hotfixBase);
                  }
                }
                sx={{ width: 275 }}
              >
                <MenuItem value={"OP"}>OP</MenuItem>
                <MenuItem value={"TiDB Cloud"}>TiDB Cloud</MenuItem>
              </Select>
            </FormControl>
          </TableCell>


          {/* <TableCell align="left">
            <FormControl>
              <InputLabel >Has fix control switch?</InputLabel>
              <Select
                value={hotfixBase.has_control_switch ? "Yes" : "No"}
                label="Has fix control switch?"
                onChange={
                  (event) => {
                    hotfixBase.has_control_switch = event.target.value == "Yes";
                    onUpdate(hotfixBase);
                  }
                }

                sx={{ width: 275 }}
              >
                <MenuItem value={"Yes"}>Yes</MenuItem>
                <MenuItem value={"No"}>No</MenuItem>
              </Select>
            </FormControl>
          </TableCell> */}

        </TableRow>

        {/* hotfixBase.has_control_switch ? <div /> :
          <TableRow>
            <TableCell colSpan={4} align="left">
              <TextField
                id="outlined-disabled"
                label="How to roll back?"
                // value={hotfixBase.roleback_method}
                onChange={
                  (event) => {
                    hotfixBase.rollback_method = event.target.value;
                    onUpdate(hotfixBase);
                  }
                }

                sx={{ width: "100%" }}
              />
            </TableCell>
          </TableRow>
        */}

        {/*
        <TableRow>
          <TableCell colSpan={4} align="left">
            <TextField
              value={hotfixBase.trigger_reason}
              id="outlined-disabled"
              label="Reason for triggering hotfix (why was it not discovered earlier)?"
              onChange={
                (event) => {
                  hotfixBase.trigger_reason = event.target.value;
                  onUpdate(hotfixBase);
                }
              }
              sx={{ width: "100%" }}
            />
          </TableCell>
        </TableRow>
*/}
      </Table>
    </TableContainer >
  )
}
