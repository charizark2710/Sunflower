import React from 'react';
import { Typography } from '@mui/material';
import { LinkAtom } from '../../atoms/link/Link.atom';

interface CopyrightMoleculesProps {
  children?: React.ReactNode;
}

export const CopyrightMolecules: React.FC<CopyrightMoleculesProps> = () => {
  return (
    <Typography variant="body2" color="text.secondary">
    {'Copyright © '}
    <LinkAtom color="inherit" target="_blank" href="https://sunflower.com/">
     Sunflower
    </LinkAtom>{' '}
    {new Date().getFullYear()}
    {'.'}
  </Typography>
  );
};
