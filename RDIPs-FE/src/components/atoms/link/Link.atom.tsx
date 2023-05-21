import React from 'react';
import Link from '@mui/material/Link';
import './Link.atom.scss';

export interface LinkAtomProps {
  onClick?: (args: any) => void;
  to: string;
  className?: string;
  children: React.ReactNode;
}
export const LinkAtom: React.FC<LinkAtomProps> = ({ onClick, to, children, className = '' }) => {
  return (
    <Link underline='none' href={to} onClick={onClick} className={className}>
      {children}
    </Link>
  );
};
