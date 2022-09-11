import React, { useState } from 'react';
import NavbarOrganism from '../organisms/common/Navbar.organism';
import HomepageTemplate from '../templates/homepage.template';

function Homepage() {
    return (
        <HomepageTemplate navbar={<NavbarOrganism />} content={<div>ABCDEFGH</div>} footer={<div>FOOTER</div>} />
    )
}

export default Homepage