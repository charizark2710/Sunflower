import { useEffect } from 'react';
import { AdminData, HeadCell } from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import './ListAdmin.scss';
import { setPage } from '../../../../redux/actions/page';
import { connect } from 'react-redux';

const ListAdmin = ({dispatch} : any) => {
  function navigateToDetailPage(detail: any) {
    return;
  }

  useEffect(() => {
    dispatch(setPage('List Admin'));
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
    <div className='list-container'>
      <TableAtom
        onRowClick={navigateToDetailPage}
        rows={userListData}
        deviceColumns={userColumns}
        title='List Admin'
        headCells={headCells}
      />
    </div>
  );
};

export default connect()(ListAdmin);
