import { Grid } from '@mui/material';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import HeaderOrganism from '../../organisms/common/header/Header.organism';
import SidebarOrganism from '../../organisms/common/sidebar/Sidebar.organism';
import HeaderTemplate from '../../templates/common/Header.template';
import './AdminPage.scss';

interface AdminPageProps {
  children?: React.ReactNode;
}

function AdminPage(props: AdminPageProps) {
  const [collapse, setCollapse] = useState(false);

  return (
    <>
      <Grid container className='admin-container'>
        <Grid item xs={collapse ? 1.2 : 2.25} display={{ xs: 'block', md: 'block' }}>
          <SidebarOrganism
            collapse={collapse}
            size={collapse ? 'sm' : 'md'}
            changeStateSideBar={() => setCollapse(!collapse)}
          />
        </Grid>
        <Grid item xs={collapse ? 10.7 : 9.75} display={{ xs: 'block', md: 'block' }} className='right-side'>
          <HeaderTemplate header={<HeaderOrganism />} />
          <div className='body-container'>{props.children}</div>
        </Grid>
      </Grid>
    </>
  );
}

export default connect()(AdminPage);
