import { useEffect } from 'react';
import { AdminData, HeadCell } from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import './ListAdmin.scss';
import { setPage } from '../../../../redux/actions/page';
import config from '../../../../utils/en.json';
import { connect } from 'react-redux';
import { AdminListIcon } from '../../../atoms/icon/ListIcon.atom';
import BreakcrumbMocules from '../../../molecules/breakcrumb/Breakcrumb.mocules';

const ListAdmin = ({ dispatch }: any) => {
  function navigateToDetailPage(detail: any) {
    return;
  }

  useEffect(() => {
    dispatch(setPage(config['adminList.title']));
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
      <div className='card-container'>
        <BreakcrumbMocules title={config['adminList.name']}
         link={config['adminList.pathLink']}
         icon={''} />
        <TableAtom
          onRowClick={navigateToDetailPage}
          rows={userListData}
          deviceColumns={userColumns}
          title={config['adminList.title']}
          headCells={headCells}
        />
      </div>
    </div>
  );
};

export default connect()(ListAdmin);
