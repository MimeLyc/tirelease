import * as React from 'react';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Typography from '@mui/material/Typography';

const steps = [
  '[TS]Submit Application',
  '[Admin]Approve Application',
  '[RD]Dev & Triage',
  '[AUTO]Building',
  '[QA]Testing',
  '[RM]Release',
];

export const HotfixTriageStepper = ({ hotfixBase = {} }) => {
  const isStepFailed = (step) => {
    return step === 1;
  };


  const isStepSuccess = (step) => {
    return step === 0;
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Stepper activeStep={1}>
        {steps.map((label, index) => {
          const labelProps = {};
          if (isStepFailed(index)) {
            labelProps.optional = (
              <Typography variant="caption" color="error">
                Denied
              </Typography>
            );

            labelProps.error = true;
          }

          return (
            <Step key={label}>
              <StepLabel {...labelProps}>{label}</StepLabel>
            </Step>
          );
        })}
      </Stepper>
    </Box>
  );
}
