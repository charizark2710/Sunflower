import { AppBar, IconButton, Toolbar, Typography } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';

import React, { useState } from 'react';
import { ButtonAtom } from '../../atoms/button/Button.atom';
import './Navbar.organism.scss';
import { LinkMolecules } from '../../molecules/link-molecules/Link.molecules';

function NavbarOrganism() {
    const [click, setClick] = useState(false);
    const [button, setButton] = useState(true);
    const handleClick = () => setClick(!click);
    const closeMobileMenu = () => setClick(false);

    const showButton = () => {
        if (window.innerWidth = 960) {
            setButton(false);
        } else {
            setButton(true);
        }
    };
    window.addEventListener('resize', showButton);


    // return (
    //     <nav className='navbar'>
    //         <div className="navbar-container">
    //             <Link to="/" className="navbar-logo">
    //                 RDIPS <i className='fab fa-typo3' />
    //             </Link>
    //             <div className='menu-icon' onClick={handleClick}>
    //                 <i className={click ? 'fas fa-times' : 'fas fa-bars'} />
    //             </div>
    //             <ul className={click ? 'nav-menu active' : 'nav-menu'}>
    //                 {/* <li className='nav-item'>
    //                         <Link to='/' className='nav-Links' onClick={closeMobileMenu}>
    //                             home
    //                         </Link>
    //                     </li>

    //                     <li className='nav-item'>
    //                         <Link to='/Services' className='nav-Links' onClick={closeMobileMenu} >
    //                             Services
    //                         </Link>
    //                     </li>
    //                     <li className='nav-item'>
    //                         <Link to='/Products' className='nav-Links' onClick={closeMobileMenu}>
    //                             Products
    //                         </Link>
    //                     </li>
    //                     <li className='nav-item'>
    //                         <Link to='/Sign-up' className='nav-Links-mobile' onClick={closeMobileMenu}>
    //                             Sign up
    //                         </Link>
    //                     </li> */}
    //                 <LinkMolecules links={[
    //                     { to: '/', className: 'nav-Links', onClick: closeMobileMenu, children: 'Home' },
    //                     { to: '/#discover', className: 'nav-Links', onClick: closeMobileMenu, children: 'Discover' },
    //                     { to: '/#about-us', className: 'nav-Links', onClick: closeMobileMenu, children: 'About Us' },
    //                     { to: '/#login', className: 'nav-Links-mobile', onClick: closeMobileMenu, children: 'Login' },
    //                 ]} isClick={click} />
    //             </ul>
    //             {button && <ButtonAtom onClick={() => { alert("TEST") }} buttonStyle='btn--outline'>SIGN UP</ButtonAtom>}
    //         </div>
    //     </nav>
    // )

    return (
        <AppBar position="static" className='navbar'>
            <Toolbar variant="dense">
                <IconButton edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }}>
                    <MenuIcon />
                </IconButton>
                {/* <IconButton edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }}>
                    <MenuIcon />
                </IconButton>
                <Typography variant="h6" color="inherit" component="div">
                    Photos
                </Typography> */}
            </Toolbar>
            <LinkMolecules links={[
                { to: '/', className: 'link-item', onClick: closeMobileMenu, children: 'Home' },
                { to: '/#discover', className: 'link-item', onClick: closeMobileMenu, children: 'Discover' },
                { to: '/#about-us', className: 'link-item', onClick: closeMobileMenu, children: 'About Us' },
               
            ]} isClick={click} />
            <div>
            {button && <ButtonAtom onClick={() => { alert("TEST") }} buttonStyle='btn--outline'>
              <Typography> SIGN UP </Typography>
            </ButtonAtom>}
            {button && <ButtonAtom onClick={() =>{ alert("go") }} buttonStyle='btn--login'>
            <Typography> Login </Typography>
            </ButtonAtom>}</div>
           
        </AppBar>
    )
}

export default NavbarOrganism