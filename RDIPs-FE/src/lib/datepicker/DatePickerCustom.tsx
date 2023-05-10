import * as React from 'react';
import dayjs, { Dayjs } from 'dayjs';
import { DemoContainer } from '@mui/x-date-pickers/internals/demo';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { Grid, Box, Button } from '@mui/material';

export default function DatePickerCustom(props: {fromDate: string, onNextClick?: any, onPrevClick?: any, disabled?: any}) {
  const [value, setValue] = React.useState<Dayjs | null>(dayjs(props.fromDate));
  return (
    <LocalizationProvider dateAdapter={AdapterDayjs}>
      <DemoContainer components={['DatePicker', 'DatePicker']}>
        <Grid container >
        <Button variant="text" onClick={props.onPrevClick}>Previous</Button>
          <Grid item xs={4}>
            <DatePicker label="From" 
            disabled={props.disabled}
            defaultValue={dayjs(props.fromDate)} />
          </Grid>
          <Grid item xs={4}>
          <DatePicker
            label="To"
            disabled={props.disabled}
            value={value}
            onChange={(newValue) => setValue(newValue)}
          />
          </Grid>
          <Button variant="text" onClick={props.onNextClick}>Next</Button>
        </Grid>
        <Box></Box>
      </DemoContainer>
    </LocalizationProvider>
  );
}