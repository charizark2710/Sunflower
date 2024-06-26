import { Box } from '@mui/material';
import { useEffect } from 'react';
import { connect } from 'react-redux';
import { setNavbarTitle } from '../../../../redux/slice/pageSlice';
import { useDispatch } from '../../../../redux/store';
import config from '../../../../utils/en.json';
import { HeadCell, TypeUserEnum, UserData } from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import BreakcrumbMocules from '../../../molecules/breakcrumb/Breakcrumb.mocules';
import './Campaign.scss';

const Campaign = () => {
  function navigateToDetailPage() {
    return;
  }

  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(setNavbarTitle(config['campaignList.title']));
  }, [dispatch]);

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
    type: TypeUserEnum
  ): UserData {
    return {
      user_id,
      user_name,
      address,
      phone_num,
      email,
      type,
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
      label: 'User Name',
    },
    {
      id: 'address',
      numeric: true,
      label: 'Address',
    },
    {
      id: 'phone_num',
      numeric: true,
      label: 'Phone number',
    },
    {
      id: 'email',
      numeric: true,
      label: 'Email',
    },
    {
      id: 'type',
      numeric: true,
      label: 'Type',
    },
  ];

  const userColumns = ['user_name', 'address', 'phone_num', 'email', 'type'];

  return (
    <Box className='list-container'>
      <Box className='card-container'>
        <BreakcrumbMocules title={config['campaignList.name']} icon={''} link={config['campaignList.pathLink']} />
        <TableAtom
          onRowClick={navigateToDetailPage}
          rows={userListData}
          deviceColumns={userColumns}
          title={config['campaignList.title']}
          headCells={headCells}
        />
      </Box>
    </Box>
  );
};

export default connect()(Campaign);
