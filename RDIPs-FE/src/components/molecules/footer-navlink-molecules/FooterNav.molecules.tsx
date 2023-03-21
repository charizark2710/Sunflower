import React from 'react';
import { LinkAtom, LinkAtomProps } from '../../atoms/link/Link.atom';
import './FooterNav.molecules.scss';

interface NavbarMoleculesProps {
  links?: Link[];
  isClick?: boolean;
  target?: string;
  children?: React.ReactNode;
  to?: string;
  className?: string;
}

interface Link extends LinkAtomProps {
  key?: string;
}

const defaultLinks = [
  { to: '/', className: 'link-item', children: 'Home' },
  { to: '/products', className: 'link-item', children: 'Products' },
  { to: '/about-us', className: 'link-item', children: 'About Us' },
  { to: '/button-name', className: 'link-item', children: 'Button Name' },
];

const defaultStyle = {
  display: 'flex',
  justifyContent: 'center',
  width: '200px',
  padding: '20px 0px',
};
export const NavbarMolecules: React.FC<NavbarMoleculesProps> = ({ links = defaultLinks, isClick, children }) => {
  return (
    <div className={isClick ? 'link-menu active' : 'link-menu'}>
      {links.map((link) => {
        return (
          <div key={link.to} style={defaultStyle}>
            <LinkAtom to={link.to} className={link.className}>
              {link.children}
            </LinkAtom>
          </div>
        );
      })}
      {children}
    </div>
  );
};
