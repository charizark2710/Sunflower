import { Box } from '@mui/material';
import { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { getAllUsers } from '../../../../axios/api';
import { setNavbarTitle } from '../../../../redux/slice/pageSlice';
import { useDispatch } from '../../../../redux/store';
import config from '../../../../utils/en.json';
import {
  HeadCell,
  TypeUserEnum,
  UserData,
  UserResponse,
} from '../../../../utils/interface';
import TableAtom from '../../../atoms/table/Table.atom';
import BreakcrumbMocules from '../../../molecules/breakcrumb/Breakcrumb.mocules';
import './Campaign.scss';

const Campaign = () => {
  const [userListData, setUserListData] = useState([]);
  function navigateToDetailPage() {
    return;
  }

  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(setNavbarTitle(config['campaignList.title']));
  }, [dispatch]);

  useEffect(() => {
    getListUser();
  }, []);

  const getListUser = () => {
    getAllUsers()
      .then((data: { data: any }) => {
        let users = data.data;
        setUserListData(
          users
            .filter((user: UserData) => user.enabled)
            .reverse()
            .map((user: UserResponse) => createData(user))
        );
      })
      .catch(() => {
        setUserListData([]);
      });
  };

  function createData(data: UserResponse): UserData {
    const {
      id,
      username,
      address,
      phone_num,
      email,
      type,
      emailVerified,
      enabled,
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
        <BreakcrumbMocules
          title={config['campaignList.name']}
          icon={''}
          link={config['campaignList.pathLink']}
        />
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
