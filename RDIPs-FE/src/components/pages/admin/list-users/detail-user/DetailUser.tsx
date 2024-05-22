import { Box, Button, Grid } from '@mui/material';
import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import chartData from '../../../../../lib/chartData.json';
import { HighChartCustom } from '../../../../../lib/highchart/HighChartCustom';
import { TypeChart } from '../../../../../utils/enum';
import { DeviceChangeHistoryData, DeviceLogHistoryData, HeadCell, StatusEnum } from '../../../../../utils/interface';
import CollapseAtom from '../../../../atoms/collapse/Collapse';
import TableAtom from '../../../../atoms/table/Table.atom';
import TitlePageAtom from '../../../../atoms/text/TitlePgae.atom';
import DetailDeviceUser from './DetailDeviceUser';
import './DetailUser.scss';
import Expense from './history/Expense';
import Feedback from './history/Feedback';
import Maintaince from './history/Maintainace';
import Receipt from './history/Receipt';
// import { DatePicker } from '@mui/x-date-pickers/DatePicker';
// import dayjs, { Dayjs } from 'dayjs';

function DatepickerByType(type: string) {
  switch (type) {
    case 'Day':
      return (
        <Box>
          Date here
          {/* <DatePicker label="From" defaultValue={dayjs('2023-09-01T00:00:00.000Z')} /> */}
        </Box>
      );
    case 'Month':
      return <Box>Month here</Box>;
    case 'Year':
      return <Box>Year here</Box>;
    case 'Decade':
      return <Box>Decade here</Box>;
  }
  return <Box>DatePicker here</Box>;
}

export const HighChartInDevice = () => {
  const [type, setType] = React.useState('1');
  const listTimeType = ['Day', 'Month', 'Year', 'Decade'];

  return (
    <>
      <Box>
        {listTimeType.map((t, i) => {
          return (
            <Button style={{ fontWeight: t === type ? 'bold' : '' }} key={i} onClick={() => setType(t)}>
              {t}
            </Button>
          );
        })}
        {DatepickerByType(type)}
      </Box>
      <HighChartCustom
        typeChart={TypeChart.sline}
        timeType={+type}
        chartData={chartData as any}
        titleChart={'Performance statistics'}
      ></HighChartCustom>
    </>
  );
};

export const HistoryLogTableInDevice = () => {
  const navigate = useNavigate();

  function navigateToDetailPage(detail: any) {
    navigate('/detail-history-log', { replace: false, state: detail });
  }

  const historyListData = [
    createData('2023-02-18', StatusEnum.Warning, 'great'),
    createData('2023-02-20', StatusEnum.Error, 'bad'),
    createData('2023-02-25', StatusEnum.Warning, 'not found'),
  ];

  function createData(datetime: string, status: StatusEnum, message: string): DeviceLogHistoryData {
    return {
      datetime,
      status,
      message,
    };
  }

  const headCells: HeadCell[] = [
    {
      numeric: undefined,
      label: 'STT',
    },
    {
      id: 'datetime',
      numeric: false,

      label: 'Datetime',
    },
    {
      id: 'status',
      numeric: true,

      label: 'Status',
    },
    {
      id: 'message',
      numeric: true,
      label: 'Message',
    },
  ];

  const logHistoryColumns = ['datetime', 'status', 'message'];

  // const [type, setType] = React.useState('1');
  return (
    <TableAtom
      onRowClick={navigateToDetailPage}
      rows={historyListData}
      deviceColumns={logHistoryColumns}
      title='Log History'
      headCells={headCells}
    />
  );
};

export const HistoryChangeTableInDevice = () => {
  const navigate = useNavigate();

  function navigateToDetailPage(detail: any) {
    navigate('/detail-history-change', { replace: false, state: detail });
  }

  const changeHistoryListData = [
    createData('2023-02-18', 'A', 'great'),
    createData('2023-02-20', 'B', 'bad'),
    createData('2023-02-25', 'C', 'not found'),
  ];

  function createData(datetime: string, type: string, description: string): DeviceChangeHistoryData {
    return {
      datetime,
      type,
      description,
    };
  }

  const headCells: HeadCell[] = [
    {
      numeric: false,
      label: 'STT',
    },
    {
      id: 'datetime',
      numeric: false,
      label: 'Datetime',
    },
    {
      id: 'type',
      numeric: true,
      label: 'Type',
    },
    {
      id: 'description',
      numeric: true,
      label: 'Description',
    },
  ];

  const changeHistoryColumns = ['datetime', 'type', 'description'];

  // const [type, setType] = React.useState('1');
  return (
    <TableAtom
      onRowClick={navigateToDetailPage}
      rows={changeHistoryListData}
      deviceColumns={changeHistoryColumns}
      title='Change History'
      headCells={headCells}
    />
  );
};

const DetailUser = () => {
  let { state } = useLocation();
  const detailUser: any = state;

  return (
    <Box className='list-container'>
      <Box className='card-container'>
        <Box>
          <TitlePageAtom title='User information'></TitlePageAtom>
          <br></br>
          <Box sx={{ flexGrow: 1 }}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Box>UserId: {detailUser.user_id}</Box>
              </Grid>
              <Grid item xs={6}>
                <Box>User Name: {detailUser.user_name}</Box>
              </Grid>
              <Grid item xs={6}>
                <Box>Firmware version: {detailUser.address}</Box>
              </Grid>
              <Grid item xs={6}>
                <Box>App version: {detailUser.phone_num}</Box>
              </Grid>
              <Grid item xs={6}>
                <Box>Type: {detailUser.email}</Box>
              </Grid>
              <Grid item xs={6}>
                <Box>Status: {detailUser.type}</Box>
              </Grid>
            </Grid>
          </Box>
        </Box>
        <br></br>
        <TitlePageAtom title='List Devices Base User'></TitlePageAtom>
        <br></br>
        <DetailDeviceUser />
        <Box className='history'>
          <TitlePageAtom title='History'></TitlePageAtom>
          {/* need to confirm more to complete */}
          <CollapseAtom buttonTitle='Maintaince' children={<Maintaince />} />
          <CollapseAtom buttonTitle='Feedback' children={<Feedback />} />
          <CollapseAtom buttonTitle='Receipt' children={<Receipt />} />
          <CollapseAtom buttonTitle='Expense' children={<Expense />} />
        </Box>
      </Box>
    </Box>
  );
};

export default DetailUser;
