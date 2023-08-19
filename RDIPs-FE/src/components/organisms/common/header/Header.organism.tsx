import { AppBar, Box, Grid } from '@mui/material';

import { useState } from 'react';
import avatar from '../../../../assets/Avatar.png';
import Image from '../../../atoms/image/Image';
import { StraightAtom } from '../../../atoms/straight/Straight.atom';
import { ClutterButtonMolecules } from '../../../molecules/clutter-button-mocules/ClutterButton.molecules';
import { NavbarMolecules } from '../../../molecules/navbar-molecules/Navbar.molecules';
import './Header.organism.scss';

function HeaderOrganism() {
  // const navbarTitle = useSelector((state: any) => state.navbarTitle);
  const [click, setClick] = useState(false);
  const isLogin = true;

  const appBarStyle = {
    bgcolor: 'white',
    color: '#000000',
    padding: '5px 24px 0px 24px',
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
            icon
          </Grid>
          <Grid item xs={10} sx={{ ...leftSideStyle, alignItems: 'center' }}>
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
