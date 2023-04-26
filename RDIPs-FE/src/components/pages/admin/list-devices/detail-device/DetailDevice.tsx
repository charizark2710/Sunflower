import React from 'react';
import { useLocation, useNavigate } from 'react-router';
import chartData from '../../../../../lib/chartData.json';
import { HighChartCustom } from '../../../../../lib/highchart/HighChartCustom';
import { TypeChart } from '../../../../../utils/enum';
import { DeviceChangeHistoryData, DeviceLogHistoryData, HeadCell, StatusEnum } from '../../../../../utils/interface';
import CollapseAtom from '../../../../atoms/collapse/Collapse';
import TableAtom from '../../../../atoms/table/Table.atom';

export const HighChartInDevice = () => {
  const [type, setType] = React.useState('1');
  return (
    <HighChartCustom
      typeChart={TypeChart.sline}
      timeType={+type}
      chartData={chartData as any}
      titleChart={'Performance statistics'}
    ></HighChartCustom>
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
      numeric: false,
      disablePadding: false,
      label: 'STT',
    },
    {
      id: 'datetime',
      numeric: false,
      disablePadding: false,
      label: 'Datetime',
    },
    {
      id: 'status',
      numeric: true,
      disablePadding: false,
      label: 'Status',
    },
    {
      id: 'message',
      numeric: true,
      disablePadding: false,
      label: 'Message',
    },
  ];

  const logHistoryColumns = ['datetime', 'status', 'message'];

  const [type, setType] = React.useState('1');
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
      disablePadding: false,
      label: 'STT',
    },
    {
      id: 'datetime',
      numeric: false,
      disablePadding: false,
      label: 'Datetime',
    },
    {
      id: 'type',
      numeric: true,
      disablePadding: false,
      label: 'Type',
    },
    {
      id: 'description',
      numeric: true,
      disablePadding: false,
      label: 'Description',
    },
  ];

  const changeHistoryColumns = ['datetime', 'type', 'description'];

  const [type, setType] = React.useState('1');
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
  const detailDevice: any = state;

  return (
    <div style={{ padding: '0 30px', backgroundColor: 'white', minHeight: '80vh' }}>
      <section>
        <h3>Device information</h3>
        <div>DeviceId: {detailDevice.device_id}</div>
        <div>Device Name: {detailDevice.device_name}</div>
        <div>Firmware version: {detailDevice.firmware_ver}</div>
        <div>App version: {detailDevice.app_ver}</div>
        <div>Type: {detailDevice.type}</div>
        <div>Status: {detailDevice.status}</div>
        <div>Life time: {detailDevice.lifetime}</div>
      </section>
      <section className='performance-statistics'>
        <CollapseAtom buttonTitle='Performance statistics' children={<HighChartInDevice />} />
      </section>
      <section className='log-history'>
        <CollapseAtom buttonTitle='Log history' children={<HistoryLogTableInDevice />} />
      </section>
      <section className='change-history'>
        <CollapseAtom buttonTitle='Change history' children={<HistoryChangeTableInDevice />} />
      </section>
    </div>
  );
};

export default DetailDevice;
