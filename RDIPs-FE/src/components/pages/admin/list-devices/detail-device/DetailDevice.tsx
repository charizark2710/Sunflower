import { Box } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router';
import chartData from '../../../../../lib/chartData.json';
import { HighChartCustom } from '../../../../../lib/highchart/HighChartCustom';
import { TypeChart } from '../../../../../utils/enum';
import { DeviceChangeHistoryData, DeviceLogHistoryData, HeadCell, StatusEnum } from '../../../../../utils/interface';
import CollapseAtom from '../../../../atoms/collapse/Collapse';
import TableAtom from '../../../../atoms/table/Table.atom';
import TextAtomDetail from '../../../../atoms/text/TextDetail.atom';
import CardMocules from '../../../../molecules/card/Card.mocules';
import { FormCreateDeviceMolecules } from '../../../../molecules/form/device-create/FormCreateDevice.molecules';
import config from '../../../../../utils/en.json';
import { getAllDevices } from '../../../../../axios/api';
import { createData } from '../ListDevices';
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
      return <div>Month here</div>;
    case 'Year':
      return <div>Year here</div>;
    case 'Decade':
      return <div>Decade here</div>;
  }
  return <div>DatePicker here</div>;
}

export const HighChartInDevice = () => {
  const [type, setType] = React.useState('1');
  const listTimeType = ['Day', 'Month', 'Year', 'Decade'];

  return (
    <>
      <div>
        {listTimeType.map((t, i) => {
          return (
            <button style={{ fontWeight: t === type ? 'bold' : '' }} key={i} onClick={() => setType(t)}>
              {t}
            </button>
          );
        })}
        {DatepickerByType(type)}
      </div>
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

const DetailDevice = () => {
  let { state } = useLocation();
  const [detailDevice, setDetailDevice]: any = useState(state);
  const [popupStatus, setPopupStatus] = useState("");

  useEffect(() => {
    getDeviceById();
  }, [popupStatus]);

  const getDeviceById = () => {
    console.log("calling me twice");
    
    let id = (state as any).device_id as string;
    getAllDevices(id)
    .then((data) => {
      setDetailDevice(createData(data.data));
    })
    .catch(() => {
      setDetailDevice(state);
    });
  }

  return (
    <div className='list-container'>
      <CardMocules title={config['deviceDetail.infoTitle']} status={popupStatus} modal={<FormCreateDeviceMolecules state="update" onClosePopUp={()=> setPopupStatus("closed")} data={detailDevice}/>}>
        <TextAtomDetail title={config['deviceDetail.device.device_name']}> {detailDevice.device_name} </TextAtomDetail>
        <TextAtomDetail title= {config['deviceDetail.device.id']}> {detailDevice.device_id} </TextAtomDetail>
        <TextAtomDetail title= {config['deviceDetail.device.firm']}> {detailDevice.firmware_ver} </TextAtomDetail>
        <TextAtomDetail title={config['deviceDetail.device.app']}> {detailDevice.app_ver} </TextAtomDetail>
        <TextAtomDetail title={config['deviceDetail.device.type']}> {detailDevice.type} </TextAtomDetail>
        <TextAtomDetail title={config['deviceDetail.device.status']} > {detailDevice.status} </TextAtomDetail>
        <TextAtomDetail title={config['deviceDetail.device.lifetime']}> {detailDevice.life_time} </TextAtomDetail>
      </CardMocules>
      <section className='performance-statistics'>
        <CollapseAtom buttonTitle={config['deviceDetail.performance.buttonTitle']} children={<HighChartInDevice />} />
      </section>
      <section className='log-history'>
        <CollapseAtom buttonTitle={config['deviceDetail.logHistory.buttonTitle']} children={<HistoryLogTableInDevice />} />
      </section>
      <section className='change-history'>
        <CollapseAtom buttonTitle={config['deviceDetail.changeHistory.buttonTitle']} children={<HistoryChangeTableInDevice />} />
      </section>
    </div>
  );
};

export default DetailDevice;
