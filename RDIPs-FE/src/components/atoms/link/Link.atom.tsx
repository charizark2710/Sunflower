import React from 'react';
import './Link.atom.scss';
import { Link } from '@mui/material';

export interface LinkAtomProps {
  onClick?: (args: any) => void;
  to?: string;
  className?: string;
  children?: React.ReactNode;
  color?: string;
  target?: string;
  href? : string;
}

export const LinkAtom: React.FC<LinkAtomProps>  = ({ onClick, to, children, className = '', color, target }) => {
  return (
    <Link underline='none' href={to} onClick={onClick} className={className} color={color} target={target}>
      {children}
    </Link>
  );
};
