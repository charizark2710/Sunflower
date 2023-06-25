import { useNavigate } from 'react-router-dom';
import { DeviceData, HeadCell } from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import './ListDevices.scss';
import { connect } from 'react-redux';
import { useEffect } from 'react';
import { setPage } from '../../../../redux/actions/page';

interface ListDevicesProps {
  dispatch: any;
}

const ListDevices: React.FC<ListDevicesProps> = ({ dispatch }) => {
  const navigate = useNavigate();
  useEffect(() => {
    dispatch(setPage('List Devices'));
  }, [dispatch]);

  const deviceListData = [
    createData('D001', 'Device01', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D002', 'Device02', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D003', 'Device03', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D004', 'Device04', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D005', 'Device05', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D006', 'Device06', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D007', 'Device07', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D008', 'Device08', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D009', 'Device09', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D0010', 'Device021', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D0011', 'Device013', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D0012', 'Device014', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D0013', 'Device015', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D0014', 'Device016', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D0015', 'Device017', 1, 1, 'ABC', 'XYZ', 'lifetime'),
    createData('D0016', 'Device018', 1, 1, 'ABC', 'XYZ', 'lifetime'),
  ];

  function createData(
    device_id: string,
    device_name: string,
    firmware_ver: number,
    app_ver: number,
    type: string,
    status: string,
    lifetime: string
  ): DeviceData {
    return {
      device_id,
      device_name,
      firmware_ver,
      app_ver,
      type,
      status,
      lifetime,
    };
  }

  const headCells: HeadCell[] = [
    {
      numeric: false,
      label: 'STT',
    },
    {
      id: 'device_name',
      numeric: false,
      label: 'Device Name',
    },
    {
      id: 'firmware_ver',
      numeric: false,
      label: 'Firmware version',
    },
    {
      id: 'app_ver',
      numeric: false,
      label: 'App version',
    },
    {
      id: 'type',
      numeric: false,
      label: 'Type',
    },
    {
      id: 'status',
      numeric: false,
      label: 'Status',
    },
    {
      id: 'lifetime',
      numeric: false,
      label: 'Lifetime',
    },
  ];

  const deviceColumns = ['device_name', 'firmware_ver', 'app_ver', 'type', 'status', 'lifetime'];

  function navigateToDetailPage(detail: any) {
    navigate('/detail-device', { replace: false, state: detail });
  }

  return (
    <div className='list-container'>
      <TableAtom
        onRowClick={navigateToDetailPage}
        rows={deviceListData}
        deviceColumns={deviceColumns}
        title='List Devices'
        headCells={headCells}
      />
    </div>
  );
};

export default connect()(ListDevices);
