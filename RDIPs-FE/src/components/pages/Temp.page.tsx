import React from 'react';
import { Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent } from '@mui/material';
import chartData from '../../lib/chartData.json';
import DatePickerCustom from '../../lib/datepicker/DatePickerCustom';
import { HighChartCustom } from '../../lib/highchart/HighChartCustom';
import { TypeChart } from '../../utils/enum';
import './Temp.page.scss';

function TempPage() {
  const [type, setType] = React.useState('1');

  const handleChange = (event: SelectChangeEvent) => {
    setType(event.target.value);
  };

  const onNextClick = () => {
    console.log('Xứng đáng là con trai của taaaaa next');
  };

  const onPrevClick = () => {
    console.log('Xứng đáng là con trai của taaaaa right');
  };

  const listTimeType = ['Day', 'Week', 'Month', 'Year', 'Decade'];

  React.useEffect(() => {
    console.log(type);
  }, [type]);

  return (
    // <div className="Sidebar" > Thấy thì mất 5 nghìn</div>

    <div style={{ padding: '15px' }}>
      <h5>Range Selection</h5>

      <Box
        sx={{
          display: 'grid',
          columnGap: 3,
          rowGap: 1,
          gridTemplateColumns: 'repeat(2, 1fr)',
        }}
      >
        <span>
          <FormControl fullWidth>
            <InputLabel id='demo-simple-select-label'>Filter By</InputLabel>
            <Select
              labelId='demo-simple-select-label'
              value={type}
              id='dateSelect'
              variant='outlined'
              sx={{ m: 2, maxWidth: 300 }}
              size='small'
              onChange={handleChange}
            >
              {listTimeType.map((type, i) => {
                return (
                  <MenuItem key={i} value={i + 1}>
                    {type}
                  </MenuItem>
                );
              })}
            </Select>
          </FormControl>
        </span>
      </Box>
      <Box>
        <DatePickerCustom
          fromDate={chartData[0].dateTime}
          onNextClick={onNextClick}
          onPrevClick={onPrevClick}
          disabled={false}
        ></DatePickerCustom>
      </Box>

      <h5>Case 1: Sline chart </h5>
      <HighChartCustom
        typeChart={TypeChart.sline}
        timeType={+type}
        chartData={chartData as any}
        titleChart={'Daily Power Chart'}
      ></HighChartCustom>
      <h5>Case 2: Bar chart </h5>
      <HighChartCustom
        typeChart={TypeChart.bar}
        timeType={+type}
        chartData={chartData as any}
        titleChart={'Daily Power Chart'}
      ></HighChartCustom>
    </div>
  );
}

export default TempPage;
