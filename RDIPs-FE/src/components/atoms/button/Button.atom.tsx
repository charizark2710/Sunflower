import React from 'react';
import './Button.atom.scss';
import Image from '../image/Image';
import addIcon from '../../../assets/icons/add-icon.svg';
import { Button } from '@mui/material';

export interface ButtonAtomProps  {
  type?: any;
  onClick?: (args: any) => void;
  buttonStyle?: string | '';
  buttonSize?: string | '';
  children?: React.ReactNode;
}
export const ButtonAtom: React.FC<ButtonAtomProps>  = ({ type, onClick, buttonStyle = '', buttonSize = '', children }) => {
  return (
    <Button className={`btn ${buttonStyle} ${buttonSize}`} onClick={onClick} type={type}>
      <Image url={addIcon} w="24" ></Image>
      {children}
    </Button>
  );
};
