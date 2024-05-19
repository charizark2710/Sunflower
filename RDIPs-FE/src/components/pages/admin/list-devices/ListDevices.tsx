import { Box } from '@mui/material';
import { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { getAllDevices } from '../../../../axios/api';
import { setNavbarTitle } from '../../../../redux/slice/pageSlice';
import config from '../../../../utils/en.json';
import { DeviceData, HeadCell } from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import BreakcrumbMocules from '../../../molecules/breakcrumb/Breakcrumb.mocules';
import { FormCreateDeviceMolecules } from '../../../molecules/form/device-create/FormCreateDevice.molecules';
import './ListDevices.scss';

interface ListDevicesProps {
  dispatch: any;
  showTableOnly: boolean;
}

interface DeviceResponse {
  id: string;
  name: string;
  firmware_ver: number;
  app_ver: number;
  type: string;
  status: string;
  life_time: string;
  region: string;
}

export function createData(data: DeviceResponse): DeviceData {
  const { id, name, firmware_ver, app_ver, type, status, life_time, region } = data;
  return {
    device_id: id,
    device_name: name,
    firmware_ver,
    app_ver,
    type,
    status,
    lifetime: life_time,
    region,
  };
}

const ListDevices: React.FC<ListDevicesProps> = ({ dispatch, showTableOnly = false }) => {
  const [deviceListData, setDeviceListData] = useState([]);
  const [popupStatus, setPopupStatus] = useState('');

  const navigate = useNavigate();
  useEffect(() => {
    dispatch(setNavbarTitle(config['deviceList.title']));
  }, [dispatch]);

  useEffect(() => {
    getListDevice();
  }, []);

  const getListDevice = () => {
    getAllDevices()
      .then((data: { data: any }) => {
        let devices = data.data;
        setDeviceListData(devices.reverse().map((device: DeviceResponse) => createData(device)));
      })
      .catch(() => {
        setDeviceListData([]);
      });
  };

  const headCells: HeadCell[] = [
    {
      label: 'No',
      numeric: undefined,
    },
    {
      id: 'device_name',
      numeric: false,
      label: config['deviceDetail.device.device_name'],
    },
    {
      id: 'user_name',
      numeric: false,
      label: config['deviceDetail.device.user_name'],
    },
    {
      id: 'region',
      numeric: false,
      label: config['deviceDetail.device.region'],
    },
    {
      id: 'lifetime',
      numeric: false,
      label: config['deviceDetail.device.lifetime'],
    },
    {
      id: 'status',
      numeric: false,
      label: config['deviceDetail.device.status'],
    },
    {
      id: '',
      numeric: false,
      label: '',
    },
  ];

  const deviceColumns = ['device_name', 'user_name', 'region', 'lifetime', 'status', ''];

  function onClosePopUp() {
    getListDevice();
    setPopupStatus('closed');
  }

  function navigateToDetailPage(detail: any) {
    navigate('/detail-device', { replace: false, state: detail });
  }

  return showTableOnly ? (
    <TableAtom
      onRowClick={navigateToDetailPage}
      rows={deviceListData}
      deviceColumns={deviceColumns}
      title={config['deviceList.title']}
      headCells={headCells}
    />
  ) : (
    <Box className='list-container'>
      <Box className='card-container'>
        <BreakcrumbMocules
          title={config['deviceList.name']}
          icon={''}
          modal={<FormCreateDeviceMolecules onClosePopUp={onClosePopUp} />}
          status={popupStatus}
          link={config['deviceList.pathLink']}
        />
        <TableAtom
          onRowClick={navigateToDetailPage}
          rows={deviceListData}
          deviceColumns={deviceColumns}
          title={config['deviceList.title']}
          headCells={headCells}
        />
      </Box>
    </Box>
  );
};

export default connect()(ListDevices);
