import { Box, Button, Grid } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { getAllDevices } from '../../../../../axios/api';
import chartData from '../../../../../lib/chartData.json';
import { HighChartCustom } from '../../../../../lib/highchart/HighChartCustom';
import config from '../../../../../utils/en.json';
import { TypeChart } from '../../../../../utils/enum';
import {
  DeviceChangeHistoryData,
  DeviceLogHistoryData,
  HeadCell,
  StatusEnum,
} from '../../../../../utils/interface';
import CollapseAtom from '../../../../atoms/collapse/Collapse';
import TableAtom from '../../../../atoms/table/Table.atom';
import CardMocules from '../../../../molecules/card/Card.mocules';
import { FormCreateDeviceMolecules } from '../../../../molecules/form/device-create/FormCreateDevice.molecules';
import { createDeviceData } from '../ListDevices';

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
            <Button
              style={{ fontWeight: t === type ? 'bold' : '' }}
              key={i}
              onClick={() => setType(t)}
            >
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
    createDeviceData('2023-02-18', StatusEnum.Warning, 'great'),
    createDeviceData('2023-02-20', StatusEnum.Error, 'bad'),
    createDeviceData('2023-02-25', StatusEnum.Warning, 'not found'),
  ];

  function createDeviceData(
    datetime: string,
    status: StatusEnum,
    message: string
  ): DeviceLogHistoryData {
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
    createDeviceData('2023-02-18', 'A', 'great'),
    createDeviceData('2023-02-20', 'B', 'bad'),
    createDeviceData('2023-02-25', 'C', 'not found'),
  ];

  function createDeviceData(
    datetime: string,
    type: string,
    description: string
  ): DeviceChangeHistoryData {
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

const DetailDevice = () => {
  let { state } = useLocation();
  const [detailDevice, setDetailDevice]: any = useState(state);
  const [popupStatus, setPopupStatus] = useState('');

  useEffect(() => {
    getDeviceById();
  }, [popupStatus]);

  const getDeviceById = () => {
    let id = (state as any).device_id as string;
    getAllDevices(id)
      .then((data: { data: any }) => {
        setDetailDevice(createDeviceData(data.data));
      })
      .catch(() => {
        setDetailDevice(state);
      });
  };

  return (
    <Box className='list-container'>
      <Box className='card-container'>
        <Box>
          <CardMocules
            title={config['deviceDetail.infoTitle']}
            status={popupStatus}
            modal={
              <FormCreateDeviceMolecules
                state='update'
                onClosePopUp={() => setPopupStatus('closed')}
                data={detailDevice}
              />
            }
          />
          <br></br>
          <Box sx={{ flexGrow: 1 }}>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Box>
                  {config['deviceDetail.device.id']}: {detailDevice.device_id}
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box>
                  {config['deviceDetail.device.device_name']}:{' '}
                  {detailDevice.device_name}
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box>
                  {config['deviceDetail.device.firm']}:{' '}
                  {detailDevice.firmware_ver}
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box>
                  {config['deviceDetail.device.app']}: {detailDevice.app_ver}
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box>
                  {config['deviceDetail.device.type']}: {detailDevice.type}
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box>
                  {config['deviceDetail.device.status']}: {detailDevice.status}
                </Box>
              </Grid>
              <Grid item xs={6}>
                <Box>
                  {config['deviceDetail.device.lifetime']}:{' '}
                  {detailDevice.life_time}
                </Box>
              </Grid>
            </Grid>
          </Box>
        </Box>
        <br></br>
        <Box className='performance-statistics'>
          <CollapseAtom
            buttonTitle={config['deviceDetail.performance.buttonTitle']}
            children={<HighChartInDevice />}
          />
          <br></br>
        </Box>
        <Box className='log-history'>
          <CollapseAtom
            buttonTitle={config['deviceDetail.logHistory.buttonTitle']}
            children={<HistoryLogTableInDevice />}
          />
          <br></br>
        </Box>
        <Box className='change-history'>
          <CollapseAtom
            buttonTitle={config['deviceDetail.changeHistory.buttonTitle']}
            children={<HistoryChangeTableInDevice />}
          />
          <br></br>
        </Box>
      </Box>
    </Box>
  );
};

export default DetailDevice;
