import { AppBar, Grid } from '@mui/material';

import { useState } from 'react';
import { StraightAtom } from '../../../atoms/straight/Straight.atom';
import { ClutterButtonMolecules } from '../../../molecules/clutter-button-mocules/ClutterButton.molecules';
import { NavbarMolecules } from '../../../molecules/navbar-molecules/Navbar.molecules';
import './Header.organism.scss';

function HeaderOrganism() {
  const [click, setClick] = useState(false);

  const appBarStyle = {
    bgcolor: '#FFFFFF',
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
      <Grid container sx={{paddingBottom: '5px'}}>
        <Grid item xs={6}>
          <NavbarMolecules isClick={click} />
        </Grid>
        <Grid item xs={6} sx={leftSideStyle}>
          <ClutterButtonMolecules />
        </Grid>
      </Grid>
      <StraightAtom />
    </AppBar>
  );
}

export default HeaderOrganism;
