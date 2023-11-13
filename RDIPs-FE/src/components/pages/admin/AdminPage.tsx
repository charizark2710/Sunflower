import { Grid } from '@mui/material';
import React, { useState } from 'react';
import { connect } from 'react-redux';
import FooterOrganism from '../../organisms/common/footer/Footer.organism';
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
        <Grid item xs={12} display={{ xs: 'block', md: 'none' }}>
          <HeaderTemplate header={<HeaderOrganism />} />
        </Grid>
        <Grid item xs={12} md={1.4} display={{ xs: 'none', md: 'block' }}>
          <SidebarOrganism size={collapse ? 'sm' : 'md'} />
        </Grid>
        <Grid item xs={12} md={10.6} className='right-side'>
          <Grid display={{ xs: 'none', md: 'block' }}>
            <HeaderTemplate header={<HeaderOrganism />} />
          </Grid>
          <div className='body-container'>{props.children}</div>
          <FooterOrganism />
        </Grid>
      </Grid>
    </>
  );
}

export default connect()(AdminPage);
