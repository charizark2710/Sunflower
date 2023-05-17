import { AppBar, Grid } from '@mui/material';

import { useState } from 'react';
import { StraightAtom } from '../../../atoms/straight/Straight.atom';
import { ClutterButtonMolecules } from '../../../molecules/clutter-button-mocules/ClutterButton.molecules';
import { NavbarMolecules } from '../../../molecules/navbar-molecules/Navbar.molecules';
import './Header.organism.scss';
import { useSelector } from 'react-redux';

function HeaderOrganism() {
  const navbarTitle = useSelector((state: any) => state.navbarTitle);
  const [click, setClick] = useState(false);
  const isLogin = true;

  const appBarStyle = {
    bgcolor: '#F3EDB5',
    color: '#000000',
    padding: '5px 10px 0px 10px',
  };

  const leftSideStyle = {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
  };

  return (
    <AppBar position='static' sx={appBarStyle}>
      {isLogin ? (
        <Grid container sx={{ pb: '5px', minHeight: '64px', px: '70px' }}>
          <Grid item xs={6}>
            {navbarTitle}
          </Grid>
          <Grid item xs={6} sx={leftSideStyle}>
            Personal Detail Frame Here
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
