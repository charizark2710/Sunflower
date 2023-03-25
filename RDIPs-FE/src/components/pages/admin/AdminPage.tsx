import React, { useState } from 'react';
import { Grid } from '@mui/material';
import FooterOrganism from '../../organisms/common/footer/Footer.organism';
import HeaderOrganism from '../../organisms/common/header/Header.organism';
import SidebarOrganism from '../../organisms/common/sidebar/Sidebar.organism';
import HeaderTemplate from '../../templates/common/Header.template';
import './AdminPage.scss';

interface AdminPageProps {
  children: React.ReactNode;
}

function AdminPage(props: AdminPageProps) {
  const [collapse, setCollapse] = useState(false);

  function toggleCollapse() {
    setCollapse(!collapse);
  }

  // React.useEffect(() => {
  //   console.log(collapse);
  // },[collapse])

  return (
    <>
      <Grid container className='admin-container'>
        <Grid item xs={collapse? 0.5: 2}>
          <SidebarOrganism onClick={toggleCollapse} size={collapse? 'sm': 'md'}/>
        </Grid>
        <Grid item xs={collapse? 11.5: 10} className='right-side'>
          <HeaderTemplate header={<HeaderOrganism />} />
          <div className='body-container'>{props.children}</div>
          <FooterOrganism />
        </Grid>
      </Grid>
    </>
  );
}

export default AdminPage;
