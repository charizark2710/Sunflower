import { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { getAllDevices } from '../../../../axios/api';
import { setPage } from '../../../../redux/actions/page';
import { DeviceData, HeadCell } from '../../../../utils/interface';
import { DeviceListIcon } from '../../../atoms/icon/ListIcon.atom';
import TableAtom from '../../../atoms/table/Table.atom';
import BreakcrumbMocules from '../../../molecules/breakcrumb/Breakcrumb.mocules';
import { FormCreateDeviceMolecules } from '../../../molecules/form/device-create/FormCreateDevice.molecules';
import config from '../../../../utils/en.json';
import './ListDevices.scss';

interface ListDevicesProps {
  dispatch: any;
}

interface DeviceResponse {
  id: string;
  name: string;
  firmware_ver: number;
  app_ver: number;
  type: string;
  status: string;
  life_time: string;
}

export function createData(data: DeviceResponse): DeviceData {
  const { id, name, firmware_ver, app_ver, type, status, life_time } = data;
  return {
    device_id: id,
    device_name: name,
    firmware_ver,
    app_ver,
    type,
    status,
    lifetime: life_time,
  };
}

const ListDevices: React.FC<ListDevicesProps> = ({ dispatch }) => {
  const [deviceListData, setDeviceListData] = useState([]);
  const [popupStatus, setPopupStatus] = useState('');

  const navigate = useNavigate();
  useEffect(() => {
    dispatch(setPage(config['deviceList.title']));
  }, [dispatch]);

  useEffect(() => {
    getListDevice();
  }, []);

  const getListDevice = () => {
    getAllDevices()
      .then((data) => {
        let devices = data.data;
        setDeviceListData(devices.reverse().map((device: DeviceResponse) => createData(device)));
      })
      .catch(() => {
        setDeviceListData([]);
      });
  };

  const headCells: HeadCell[] = [
    {
      numeric: false,
      label: 'STT',
    },
    {
      id: 'device_name',
      numeric: false,
      label: config['deviceDetail.device.name'],
    },
    {
      id: 'firmware_ver',
      numeric: false,
      label: config['deviceDetail.device.firm'],
    },
    {
      id: 'app_ver',
      numeric: false,
      label: config['deviceDetail.device.app'],
    },
    {
      id: 'type',
      numeric: false,
      label: config['deviceDetail.device.type'],
    },
    {
      id: 'status',
      numeric: false,
      label: config['deviceDetail.device.status'],
    },
    {
      id: 'lifetime',
      numeric: false,
      label: config['deviceDetail.device.lifetime'],
    },
  ];

  const deviceColumns = ['device_name', 'firmware_ver', 'app_ver', 'type', 'status', 'lifetime'];

  function onClosePopUp() {
    getListDevice();
    setPopupStatus('closed');
  }

  function navigateToDetailPage(detail: any) {
    navigate('/detail-device', { replace: false, state: detail });
  }

  return (
    <div className='list-container'>
      <div className='card-container'>
        <BreakcrumbMocules
          title={config['deviceList.name']}
          icon={<DeviceListIcon />}
          modal={<FormCreateDeviceMolecules onClosePopUp={onClosePopUp} />}
          status={popupStatus}
        />
        <TableAtom
          onRowClick={navigateToDetailPage}
          rows={deviceListData}
          deviceColumns={deviceColumns}
          title={config['deviceList.title']}
          headCells={headCells}
        />
      </div>
    </div>
  );
};

export default connect()(ListDevices);
