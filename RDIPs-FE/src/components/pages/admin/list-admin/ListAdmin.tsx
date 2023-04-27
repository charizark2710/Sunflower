import { HeadCell, TypeUserEnum, UserData } from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import './ListAdmin.scss';

const ListAdmin = () => {
  function navigateToDetailPage(detail: any) {
    return;
  }

  const userListData = [
    createData('U001', 'Ly Nguyen', '123 Thien Duong', '09876543212', 'thienduong@gmail.com', TypeUserEnum.Regular),
    createData('U002', 'Anh Phan', '502 Thien Duong', '09876543213', 'thienduong1@gmail.com', TypeUserEnum.Industrial),
    createData('U003', 'Canh Ngo', '503 Thien Duong', '098765432132', 'thienduong2@gmail.com', TypeUserEnum.Regular),
    createData('U004', 'Thanh Bui', '504 Thien Duong', '09876543215', 'thienduong3@gmail.com', TypeUserEnum.Regular),
    createData('U005', 'Minh Hung', '505 Thien Duong', '09876543212', 'thienduon4g@gmail.com', TypeUserEnum.Regular),
    createData('U006', 'Huong Nguyen', '506 Thien Duong', '09876543212', 'thienduong5@gmail.com', TypeUserEnum.Industrial),
    createData('U007', 'Huy Doan', '507 Thien Duong', '09876543212', 'thienduong6@gmail.com', TypeUserEnum.Regular),
  ];

  function createData(
    user_id: string,
    user_name: string,
    address: string,
    phone_num: string,
    email: string,
    type: TypeUserEnum,
  ): UserData {
    return {
      user_id,
      user_name,
      address,
      phone_num,
      email,
      type
    };
  }

  const headCells: HeadCell[] = [
    {
      numeric: false,
      disablePadding: false,
      label: 'STT',
    },
    {
      id: 'user_name',
      numeric: false,
      disablePadding: false,
      label: 'User Name',
    },
    {
      id: 'address',
      numeric: true,
      disablePadding: false,
      label: 'Address',
    },
    {
      id: 'phone_num',
      numeric: true,
      disablePadding: false,
      label: 'Phone number',
    },
    {
      id: 'email',
      numeric: true,
      disablePadding: false,
      label: 'Email',
    },
    {
      id: 'type',
      numeric: true,
      disablePadding: false,
      label: 'Type',
    }
  ];

  const userColumns = ['user_name', 'address', 'phone_num', 'email', 'type'];

  return (
    <div className='list-admin-container'>
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

export default ListAdmin;