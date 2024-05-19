import { Box } from '@mui/material';
import { useEffect } from 'react';
import { connect } from 'react-redux';
import { setNavbarTitle } from '../../../../redux/slice/pageSlice';
import config from '../../../../utils/en.json';
import { AdminData, HeadCell } from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import BreakcrumbMocules from '../../../molecules/breakcrumb/Breakcrumb.mocules';
import './ListAdmin.scss';

const ListAdmin = ({ dispatch }: any) => {
  function navigateToDetailPage(detail: any) {
    return;
  }

  useEffect(() => {
    dispatch(setNavbarTitle(config['adminList.title']));
  }, [dispatch]);

  const userListData = [
    createData('AD001', 'Ly Nguyen', 'single', 'dev', 'authernication'),
    createData('AD002', 'Anh Phan', 'single', 'dev', 'authernication'),
    createData('AD003', 'Canh Ngo', 'single', 'dev', 'authernication'),
  ];

  function createData(admin_id: string, admin_name: string, status: string, role: string, auhentication: string): AdminData {
    return {
      admin_id,
      admin_name,
      status,
      role,
      auhentication,
    };
  }

  const headCells: HeadCell[] = [
    {
      numeric: undefined,
      label: 'STT',
    },
    {
      id: 'user_name',
      numeric: false,
      label: 'Admin Name',
    },
    {
      id: 'address',
      numeric: false,
      label: 'Staus',
    },
    {
      id: 'phone_num',
      numeric: false,
      label: 'Role',
    },
    {
      id: 'email',
      numeric: false,
      label: 'Authentication',
    },
  ];

  const userColumns = ['admin_name', 'status', 'role', 'authentication'];

  return (
    <Box className='list-container'>
      <Box className='card-container'>
        <BreakcrumbMocules title={config['adminList.name']} link={config['adminList.pathLink']} icon={''} />
        <TableAtom
          onRowClick={navigateToDetailPage}
          rows={userListData}
          deviceColumns={userColumns}
          title={config['adminList.title']}
          headCells={headCells}
        />
      </Box>
    </Box>
  );
};

export default connect()(ListAdmin);
