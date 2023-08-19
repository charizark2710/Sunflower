import { Box, Grid } from '@mui/material';
import { StraightAtom } from '../../../atoms/straight/Straight.atom';
import { NavbarMolecules } from '../../../molecules/navbar-molecules/Navbar.molecules';
import './Footer.organism.scss';

const dataFooterRender = {
  linkAdditionLeft: [
    { name: 'About us', link: '' },
    { name: 'Services', link: '' },
    { name: 'Features', link: '' },
    { name: 'Resources', link: '' },
  ],
  linkAdditionRight: [
    { name: 'Terms Of Use', link: '' },
    { name: 'Privacy Policy', link: '' },
  ],
  copyRight: '© 2022 Sunflower ProjectProject. All rights reserved',
};

const footerLeftLinks = [
  {
    to: dataFooterRender.linkAdditionLeft[0].link,
    className: 'link-item',
    children: dataFooterRender.linkAdditionLeft[0].name,
  },
  {
    to: dataFooterRender.linkAdditionLeft[1].link,
    className: 'link-item',
    children: dataFooterRender.linkAdditionLeft[1].name,
  },
  {
    to: dataFooterRender.linkAdditionLeft[2].link,
    className: 'link-item',
    children: dataFooterRender.linkAdditionLeft[2].name,
  },
  {
    to: dataFooterRender.linkAdditionLeft[3].link,
    className: 'link-item',
    children: dataFooterRender.linkAdditionLeft[3].name,
  },
];

const footerRightLinks = [
  {
    to: dataFooterRender.linkAdditionRight[0].link,
    className: 'link-item',
    children: dataFooterRender.linkAdditionRight[0].name,
  },
  {
    to: dataFooterRender.linkAdditionRight[1].link,
    className: 'link-item',
    children: dataFooterRender.linkAdditionRight[1].name,
  },
];

function FooterOrganism() {
  return (
    <footer className='footer-container'>
      <Grid container>
        <Grid item xs={5} sx={{ paddingLeft: '30px' }}>
          <NavbarMolecules links={footerLeftLinks} />
        </Grid>
      </Grid>
      <StraightAtom color='#9F9F9F' />
      <Grid container sx={{ padding: '10px 0px' }} className='flex-align-center'>
        <Grid item xs={11} sx={{ paddingLeft: '40px' }}>
          {dataFooterRender.copyRight}
        </Grid>
        <Grid item xs={1} className='flex-justify-end'>
          <NavbarMolecules links={footerRightLinks} />
        </Grid>
      </Grid>
    </footer>
  );
}

export default FooterOrganism;
