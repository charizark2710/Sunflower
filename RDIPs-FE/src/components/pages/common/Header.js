import React, { useState } from 'react';
import NavbarOrganism from '../../organisms/common/Navbar.organism';
import HeaderTemplate from '../../templates/common/Header.template';

function Header() {
    return (
        <HeaderTemplate navbar={<NavbarOrganism />} />
    )
}

export default Header