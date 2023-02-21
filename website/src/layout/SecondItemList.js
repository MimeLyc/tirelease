import * as React from "react";
import { Link } from "react-router-dom";

import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import ListSubheader from "@mui/material/ListSubheader";
import ColorizeIcon from "@mui/icons-material/Colorize";
import SendTimeExtensionIcon from '@mui/icons-material/SendTimeExtension';
import AdUnitsIcon from "@mui/icons-material/AdUnits";
import { useSearchParams } from "react-router-dom";

const SecondItemList = () => {
  // Feature flag for control the visibility.
  const [searchParams, setSearchParams] = useSearchParams();
  var isInDev = searchParams.get('isInDev') || 'false';

  return (
    <div>
      <ListSubheader inset>Triage Management</ListSubheader>
      <ListItem button component={Link} to="/home/triage">
        <ListItemIcon>
          <ColorizeIcon />
        </ListItemIcon>
        <ListItemText primary="Version Triage" />
      </ListItem>
      {
        isInDev === 'true' ?
          <ListItem button component={Link} to="/home/hotfix">
            <ListItemIcon>
              <SendTimeExtensionIcon />
            </ListItemIcon>
            <ListItemText primary="Hotfix Triage" />
          </ListItem> : <dev />
      }
      <ListItem button component={Link} to="/home/affection">
        <ListItemIcon>
          <AdUnitsIcon />
        </ListItemIcon>
        <ListItemText primary="Affects Triage" />
      </ListItem>
    </div >
  )
};

export default SecondItemList;
