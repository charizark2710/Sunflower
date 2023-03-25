import React from 'react';
import './Straight.atom.scss';

export interface StraightAtomProps {
  color?: string;
  width?: string;
  thick?: string;
  className?: string;
}
export const StraightAtom: React.FC<StraightAtomProps> = ({ className = '', color = '#000000', width = '90%', thick= '1px' }) => {
  return (
    <div className='flex-justify-center' style={{padding: '0', margin: '0'}}>
      <div style={{ borderBottom: `${thick} solid ${color}`, width: width}}></div>
    </div>
  );
};
