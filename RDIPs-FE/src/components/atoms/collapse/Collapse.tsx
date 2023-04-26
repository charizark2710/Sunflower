import React from 'react';
import { Button } from '@mui/material';

interface CollapseAtomProps {
  children: React.ReactNode;
  buttonTitle?: string;
}

const CollapseAtom: React.FC<CollapseAtomProps> = ({ children, buttonTitle = 'Collapse here' }) => {
  const [open, setOpen] = React.useState(false);
  return (
    <div>
      <Button onClick={() => setOpen(!open)}>{buttonTitle}</Button>
      {open ? children : ''}
    </div>
  );
};

export default CollapseAtom;
