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
        {/* <Grid item xs={12} display={{ xs: 'block', md: 'none' }}>
          <HeaderTemplate header={<HeaderOrganism />} />
        </Grid> */}
        <div className='left-side'>
          <SidebarOrganism size={collapse ? 'sm' : 'md'} />
        </div>
        <div className='right-side'>
          <HeaderTemplate header={<HeaderOrganism />} />
          <div className='body-container'>{props.children}</div>
        </div>
      </Grid>
    </>
  );
}

export default connect()(AdminPage);
