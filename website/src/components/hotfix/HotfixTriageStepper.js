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

function getStepInfo(hotfixBase) {
  switch (hotfixBase.status.toUpperCase()) {
    case "INIT":
      return {
        step: 0,
        is_success: true,
      }
    case "PENDING_APPROVAL":
      return {
        step: 1,
        is_success: true,
        message: "Waiting for approval",
      }
    case "DENIED":
      return {
        step: 1,
        is_success: false,
        message: "Application Denied!",
      }
    case "UPCOMING":
      return {
        step: 2,
        is_success: true,
        message: "",
      }
    case "BUILDING":
      return {
        step: 3,
        is_success: true,
        message: "Waiting for building result",
      }
    case "QA_TESTING":
      return {
        step: 4,
        is_success: true,
        message: "Waiting for qa testing",
      }
  }
}

export const HotfixTriageStepper = ({ hotfixBase = {}, hotfixRelease = {} }) => {

  const stepInfo = getStepInfo(hotfixBase);
  return (
    <Box sx={{ width: '100%' }}>
      <Stepper activeStep={stepInfo?.step}>
        {steps.map((label, index) => {
          const labelProps = {};
          if (!stepInfo.is_success) {
            labelProps.optional = (
              <Typography variant="caption" color="error">
                {stepInfo.message}
              </Typography>
            );

            labelProps.error = true;
          } else if (stepInfo.message || "" != "") {
            labelProps.optional = (
              <Typography variant="caption" color="default">
                {stepInfo.message}
              </Typography>
            );
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
