import NativeSelect from '@mui/material/NativeSelect';
import {
  TextField, Typography, Chip,
  InputAdornment, Select, MenuItem, FormControl, InputLabel
} from '@mui/material';
import SyncIcon from '@mui/icons-material/Sync'
import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

function StatusBar(props) {
  if (props.status == "UPCOMING") {
    return <SyncIcon color="primary" />
  }
  return <div > props.status</div>
}

export const HotfixTriageBaseInfo = ({ onUpdate, hotfixBase = {} }) => {
  const user = hotfixBase.creator

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
            <TextField
              disabled
              value={hotfixBase.customer}
              label="Customer" />
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

          <TableCell colSpan={1}>
            <Typography gutterBottom variant="h8" component="div">
              <StatusBar status={hotfixBase.status} />
            </Typography>
          </TableCell>

        </TableRow>

        <TableRow>
          <TableCell align="left">
            <TextField
              label="Oncall ID"
              sx={{ width: 275 }}
              value={hotfixBase.oncall_id}
              disabled
              // onChange={(event) => {
              //   hotfixBase.oncall_id = event.target.value;
              //   onUpdate(hotfixBase);
              // }}
              InputProps={{
                startAdornment:
                  <InputAdornment position="start">
                    <FormControl fullWidth>
                      <NativeSelect
                        sx={{ width: 100 }}
                        disabled
                        defaultValue={hotfixBase.oncall_prefix}
                        value={hotfixBase.oncall_prefix}
                        // onChange={(event) => {
                        //   hotfixBase.oncall_prefix = event.target.value;
                        //   onUpdate(hotfixBase);
                        // }}
                        inputProps={{
                          name: 'oncall_platform',
                          id: 'uncontrolled-native',
                        }}
                      >
                        // <option value={"oncall"}>ONCALL</option>
                        // <option value={"ticket"}>TICKET</option>
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
              disabled
              value={hotfixBase.oncall_url}
            // onChange={
            //   (event) => {
            //     hotfixBase.oncall_url = event.target.value;
            //     onUpdate(hotfixBase);
            //   }
            // }
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
                disabled
                // onChange={
                //   (event) => {
                //     hotfixBase.is_debug = event.target.value == "Yes";
                //     onUpdate(hotfixBase);
                //   }
                // }
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
                disabled
                // onChange={
                //   (event) => {
                //     hotfixBase.platform = event.target.value;
                //     onUpdate(hotfixBase);
                //   }
                // }
                sx={{ width: 275 }}
              >
                <MenuItem value={"OP"}>OP</MenuItem>
                <MenuItem value={"TiDB Cloud"}>TiDB Cloud</MenuItem>
              </Select>
            </FormControl>
          </TableCell>

        </TableRow>

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

      </Table>
    </TableContainer >
  )
}
