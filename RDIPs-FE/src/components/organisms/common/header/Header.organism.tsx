import { AppBar, Box, Grid } from '@mui/material';

import { useState } from 'react';
import avatar from '../../../../assets/Avatar.png';
import notiOn from '../../../../assets/icons/noti-on.svg';
import notiOff from '../../../../assets/icons/noti-off.svg';
import Image from '../../../atoms/image/Image';
import { StraightAtom } from '../../../atoms/straight/Straight.atom';
import { ClutterButtonMolecules } from '../../../molecules/clutter-button-mocules/ClutterButton.molecules';
import { NavbarMolecules } from '../../../molecules/navbar-molecules/Navbar.molecules';
import './Header.organism.scss';

function HeaderOrganism() {
  // const navbarTitle = useSelector((state: any) => state.navbarTitle);
  const [click, setClick] = useState(false);
  const [notifyOn, setNotifyOn] = useState(true);
  const isLogin = true;

  const appBarStyle = {
    bgcolor: 'white',
    color: '#25205B',
    padding: '5px 32px 0px 32px',
    borderBottom: '1px solid #e4e4e4',
    boxShadow: 'unset',
  };

  const leftSideStyle = {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
  };

  return (
    <AppBar position='static' sx={appBarStyle}>
      {isLogin ? (
        <Grid container sx={{ pb: '5px', height: '64px', alignItems: 'center'}}>
          <Grid item xs={2} fontWeight='bold'>
            {/* icon */}
          </Grid>
          <Grid item xs={10} sx={{ ...leftSideStyle, alignItems: 'center' }}>
             <Box mr={1}>
              <Image url={notifyOn ? notiOn : notiOff} w={'24px'} />
            </Box>
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
          {/* <StraightAtom /> */}
        </>
      )}
    </AppBar>
  );
}

export default HeaderOrganism;
