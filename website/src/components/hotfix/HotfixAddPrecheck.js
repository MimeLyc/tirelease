import * as React from "react";
import { Link, Stack, } from '@mui/material';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';
import FormLabel from '@mui/material/FormLabel';

// Creators must check the box before they can submit a hotfix
export const HotfixAddPrecheck = ({ }) => {
  const [masterCheck, setMasterCheck] = React.useState(false);
  const [bugCheck, setBugCheck] = React.useState(false);
  const [workaroundCheck, setWorkaroundCheck] = React.useState(false);
  const [waitingCheck, setWaitingCheck] = React.useState(false);
  const [prCheck, setPrCheck] = React.useState(false);

  return (
    // cc https://mui.com/material-ui/react-radio-button/
    <Stack direction="column">

      <FormLabel>
        How about oncall master suggest for the hotfix?
      </FormLabel>
      <RadioGroup >
        <FormControlLabel
          value={true} control={<Radio />}
          label="Oncall master suggest to apply the hotfix." />
        <FormControlLabel
          value={false} control={<Radio color="default" />}
          label="Other." />
      </RadioGroup>

      <FormLabel>
        Is it a bug fix?
      </FormLabel>
      <RadioGroup>
        <FormControlLabel
          value={true} control={<Radio />}
          label="Yes, it's a bug fix." />
        <FormControlLabel
          value={false} control={<Radio color="default" />}
          label="No, it's a new feature or over the feature design." />
      </RadioGroup>

      <FormLabel>
        Do we have a workaround for user?
      </FormLabel>
      <RadioGroup>
        <FormControlLabel
          value={true} control={<Radio />}
          label="No, we don't have any workaround." />
        <FormControlLabel
          value={false} control={<Radio color="default" />}
          label="Yes we have a workaround, but we still need a hotfix." />
      </RadioGroup>

      <FormLabel>
        Waiting for the next TiDB patch version?
        <Link underline="hover" color="inherit" href="https://github.com/pingcap/tidb/projects/63">
          TiDB release plan
        </Link>
      </FormLabel>
      <RadioGroup>
        <FormControlLabel
          value={true} control={<Radio />}
          label="No, it's urgent." />
        <FormControlLabel
          value={false} control={<Radio color="default" />}
          label="Yes." />
      </RadioGroup>


      <FormLabel>
        Do we have already fixed the issue and the master RP merged with test complete?
      </FormLabel>
      <RadioGroup>
        <FormControlLabel
          value={true} control={<Radio />}
          label="Yes, master PR merged with test complete." />
        <FormControlLabel u
          value={false} control={<Radio color="default" />}
          label="No, have not fixed yet." />
      </RadioGroup>
    </Stack>
  )
}
