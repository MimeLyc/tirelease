import { BaseVersionSelector } from "./BaseVersionSelector";
import {
  Table, TableCell, TableRow, Stack
} from '@mui/material';

import DialogContentText from "@mui/material/DialogContentText";

import storage from '../common/LocalStorage';

export const HotfixAddBaseInfo = ({ onUpdate, hotfixBase = {} }) => {
  let user = storage.getUser();

  return (
    <Table>
      <TableRow>
        <TableCell colSpan={3}>
          <BaseVersionSelector
            onMajorChange={
              (major) => {
                hotfixBase.major = major
                onUpdate(hotfixBase)
              }
            }
            onMinorChange={
              (minor) => {
                hotfixBase.minor = minor
                onUpdate(hotfixBase)
              }
            }
            onPatchChange={
              (value) => {
                hotfixBase.patch = value
                onUpdate(hotfixBase)
              }
            }
          />
        </TableCell>

        {/* Automatically set owner */}
        <TableCell colSpan={1}>
          <Stack direction="column" spacing={2} alignItems="top">

            <DialogContentText>Creator</DialogContentText>
            <Stack direction="row" spacing={2} alignItems="flex-top">
              {`${user?.name}(${user?.git_login})`}
            </Stack>
          </Stack>

        </TableCell>
      </TableRow>
    </Table>
  )
}
