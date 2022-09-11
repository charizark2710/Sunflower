import React from 'react';
import Link from '@mui/material/Link';
import './Link.atom.scss'
             
export const LinkAtom = ({ onClick, to, children, className }) => {
    return (
        <Link
            underline='none'
            href={to}
            onClick={onClick}
            className={className}
        >
            {children}
        </Link>
    )
};