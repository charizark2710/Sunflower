import React from 'react';
import './Button.atom.scss';
import { Button } from '@mui/material';

const STYLES = ['btn--primary', 'btn--outline', 'btn--login'];

const SIZES = ['btn--medium', 'btn--large'];

export interface ButtonAtomProps {
  type?: any;
  onClick?: (args: any) => void;
  buttonStyle?: string | '';
  buttonSize?: string | '';
  children: React.ReactNode;
}
export const ButtonAtom: React.FC<ButtonAtomProps> = ({ type, onClick, buttonStyle = '', buttonSize = '', children }) => {
  const checkStyle: string = STYLES.includes(buttonStyle) ? buttonStyle : STYLES[0];

  const checkSize: string = SIZES.includes(buttonSize) ? buttonSize : SIZES[0];
  return (
    <Button className={`btn ${checkStyle} ${checkSize}`} onClick={onClick} type={type}>
      {children}
    </Button>
  );
};
