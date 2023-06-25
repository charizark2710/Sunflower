import { AppBar, Box, Grid } from '@mui/material';

import NotificationsIcon from '@mui/icons-material/Notifications';
import Badge from '@mui/material/Badge';
import IconButton from '@mui/material/IconButton';
import { useState } from 'react';
import { useSelector } from 'react-redux';
import avatar from '../../../../assets/Avatar.png';
import gear from '../../../../assets/Gear.svg';
import Image from '../../../atoms/image/Image';
import { StraightAtom } from '../../../atoms/straight/Straight.atom';
import { ClutterButtonMolecules } from '../../../molecules/clutter-button-mocules/ClutterButton.molecules';
import { NavbarMolecules } from '../../../molecules/navbar-molecules/Navbar.molecules';
import './Header.organism.scss';

function HeaderOrganism() {
  const navbarTitle = useSelector((state: any) => state.navbarTitle);
  const [click, setClick] = useState(false);
  const isLogin = true;

  const appBarStyle = {
    bgcolor: 'white',
    color: '#000000',
    padding: '5px 24px 0px 10px',
    borderBottom: '1px solid #e4e4e4',  
    boxShadow: 'unset'
  };

  const leftSideStyle = {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
  };

  return (
    <AppBar position='static' sx={appBarStyle}>
      {isLogin ? (
        <Grid container sx={{ pb: '5px', height: '64px', alignItems: 'center', pl: '70px' }}>
          <Grid item xs={5} fontWeight='bold'>
            {navbarTitle ? navbarTitle : ''}
          </Grid>
          <Grid item xs={7} sx={{ ...leftSideStyle, alignItems: 'center' }}>
            {/* <div style={{ paddingRight: '15px', fontWeight: 'bold' }}>
              <div>Admin: PBTAnh</div>
              <div>#696969</div>
            </div>
            <Box>
              <Image url={avatar} w={'40px'} />
              <IconButton sx={{ transform: 'translate(-50%, -50%)', ml: '20px' }}>
                <Badge badgeContent={100} color='secondary'>
                  <NotificationsIcon />
                </Badge>
              </IconButton>
              <IconButton sx={{ transform: 'translate(-50%, -50%)', ml: '5px' }}>
                <Image url={gear} w={'30px'} />
              </IconButton>
            </Box> */}
            <Box>
              <Image url={avatar} w={'40px'} />
            </Box>
          </Grid>
        </Grid>
      ) : (
        <>
          <Grid container sx={{ paddingBottom: '5px' }}>
            <Grid item xs={6}>
              <NavbarMolecules isClick={click} />
            </Grid>
            <Grid item xs={6} sx={leftSideStyle}>
              <ClutterButtonMolecules />
            </Grid>
          </Grid>
          <StraightAtom />
        </>
      )}
    </AppBar>
  );
}

export default HeaderOrganism;
