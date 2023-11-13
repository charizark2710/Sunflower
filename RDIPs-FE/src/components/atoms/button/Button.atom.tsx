import React from 'react';
import './Button.atom.scss';
import { Button } from '@mui/material';

export interface ButtonAtomProps  {
  type?: any;
  onClick?: (args: any) => void;
  buttonStyle?: string | '';
  buttonSize?: string | '';
  children: React.ReactNode;
}
export const ButtonAtom: React.FC<ButtonAtomProps>  = ({ type, onClick, buttonStyle = '', buttonSize = '', children }) => {
  return (
    <Button className={`btn ${buttonStyle} ${buttonSize}`} onClick={onClick} type={type}>
      {children}
    </Button>
  );
};
