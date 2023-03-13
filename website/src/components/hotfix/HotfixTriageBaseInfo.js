import NativeSelect from '@mui/material/NativeSelect';
import {
  TextField, Typography, Chip, Autocomplete, Link,
  InputAdornment, Select, MenuItem, FormControl, InputLabel
} from '@mui/material';
import SyncIcon from '@mui/icons-material/Sync'
import Table from '@mui/material/Table';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';

function StatusBar(props) {
  const status = props.status.toUpperCase();
  if (status == "UPCOMING") {
    return <Chip color="primary" label={status} />
  }
  return <Chip>status</Chip>
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
              variant="standard"
              value={hotfixBase.customer}
              label="Customer" />
          </TableCell>

          <TableCell colSpan={1}>
            <TextField
              disabled
              variant="standard"
              value={`${hotfixBase.customer}`}
              sx={{ width: 150 }}
              label="Submitor"
            />
          </TableCell>

          <TableCell>
            <TextField
              disabled
              variant="standard"
              label="Related Oncall"
              sx={{ width: 150 }}
              InputProps={{
                startAdornment: (
                  <Link
                    sx={{ width: 275 }}
                    href={hotfixBase.oncall_url} target="_blank">
                    {`${hotfixBase.oncall_prefix}-${hotfixBase.oncall_id}`}
                  </Link>
                ),
              }}
            />

          </TableCell>
          <TableCell align="left">
            <FormControl
            >
              <InputLabel id="demo-simple-select-label">Is for debug?</InputLabel>
              <Select
                value={hotfixBase.is_debug ? "Yes" : "No"}
                label="Is it a debug hotfix?"
                variant="standard"
                disabled
                sx={{ width: 150 }}
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
                variant="standard"
                label="Whether the user has applied any hotfix before?"
                disabled
                sx={{ width: 150 }}
              >
                <MenuItem value={"OP"}>OP</MenuItem>
                <MenuItem value={"TiDB Cloud"}>TiDB Cloud</MenuItem>
              </Select>
            </FormControl>
          </TableCell>

          <TableCell colSpan={1}>
            <Typography gutterBottom variant="h8" component="div">
              <StatusBar status={hotfixBase.status} />
            </Typography>
          </TableCell>

        </TableRow>


        <TableRow>
          <TableCell colSpan={4} align="left">
            <TextField
              value={hotfixBase.trigger_reason}
              multiline
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
