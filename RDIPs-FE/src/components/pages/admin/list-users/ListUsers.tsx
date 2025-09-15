import { Box } from '@mui/material';
import { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { getAllUsers } from '../../../../axios/api';
import { setNavbarTitle } from '../../../../redux/slice/pageSlice';
import config from '../../../../utils/en.json';
import {
  HeadCell,
  TypeUserEnum,
  UserData,
  UserResponse,
} from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import BreakcrumbMocules from '../../../molecules/breakcrumb/Breakcrumb.mocules';
import { FormCreateUserMolecules } from '../../../molecules/form/user-create/FormCreateUser.molecules';
import './ListUsers.scss';

export function createUserData(data: UserResponse): UserData {
  const {
    id,
    username,
    address,
    phone_num,
    email,
    type,
    emailVerified,
    enabled,
    firstName,
    lastName,
  } = data;
  return {
    user_id: id,
    user_name: username,
    address: address ? address : '',
    phone_num: phone_num ? phone_num : '',
    email: email ? email : '',
    type: type ? type : TypeUserEnum.Regular,
    emailVerified,
    enabled,
    firstName,
    lastName,
  };
}

const ListUsers = ({ dispatch }: any) => {
  const [userListData, setUserListData] = useState([]);
  const [popupStatus, setPopupStatus] = useState('');
  const navigate = useNavigate();
  function navigateToDetailPage(detail: any) {
    navigate('/detail-user', { replace: false, state: detail });
  }

  useEffect(() => {
    dispatch(setNavbarTitle(config['usersList.title']));
  }, [dispatch]);

  useEffect(() => {
    getListUser();
  }, []);

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

  const getListUser = () => {
    getAllUsers()
      .then((data: { data: any }) => {
        let users = data.data;
        setUserListData(
          users
            .filter((user: UserData) => user.enabled)
            .reverse()
            .map((user: UserResponse) => createUserData(user))
        );
      })
      .catch(() => {
        setUserListData([]);
      });
  };

  // const userListData = [
  //   createUserData('U001', 'Ly Nguyen', '123 Thien Duong', '09876543212', 'thienduong@gmail.com', TypeUserEnum.Regular),
  //   createUserData('U002', 'Anh Phan', '502 Thien Duong', '09876543213', 'thienduong1@gmail.com', TypeUserEnum.Industrial),
  //   createUserData('U003', 'Canh Ngo', '503 Thien Duong', '098765432132', 'thienduong2@gmail.com', TypeUserEnum.Regular),
  //   createUserData('U004', 'Thanh Bui', '504 Thien Duong', '09876543215', 'thienduong3@gmail.com', TypeUserEnum.Regular),
  //   createUserData('U005', 'Minh Hung', '505 Thien Duong', '09876543212', 'thienduon4g@gmail.com', TypeUserEnum.Regular),
  //   createUserData('U006', 'Huong Nguyen', '506 Thien Duong', '09876543212', 'thienduong5@gmail.com', TypeUserEnum.Industrial),
  //   createUserData('U007', 'Huy Doan', '507 Thien Duong', '09876543212', 'thienduong6@gmail.com', TypeUserEnum.Regular),
  // ];

  function onClosePopUp() {
    getListUser();
    setPopupStatus('closed');
  }

  return (
    <Box className='list-container'>
      <Box className='card-container'>
        <BreakcrumbMocules
          status={popupStatus}
          modal={<FormCreateUserMolecules onClosePopUp={onClosePopUp} />}
          title={config['usersList.name']}
          icon={''}
          link={config['usersList.pathLink']}
        />
        <TableAtom
          onRowClick={navigateToDetailPage}
          rows={userListData}
          deviceColumns={userColumns}
          title={config['usersList.title']}
          headCells={headCells}
        />
      </Box>
    </Box>
  );
};

export default connect()(ListUsers);
