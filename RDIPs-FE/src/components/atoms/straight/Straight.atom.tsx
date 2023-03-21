import React from 'react';
import './Straight.atom.scss';

export interface StraightAtomProps {
  color?: string;
  className?: string;
}
export const StraightAtom: React.FC<StraightAtomProps> = ({ className = '', color = '#000000' }) => {
  return (
    <div className='flex-justify-center' style={{padding: '0', margin: '0'}}>
      <div style={{ border: `1px solid ${color}`, width: '90%' }}></div>
    </div>
  );
};
