import React from 'react';
import './Search.atom.scss';
import Image from '../image/Image';
import searchIcon from '../../../assets/icons/search-icon.svg';

export interface SearchAtomProps {
  onSearch?: (args: any) => void;
  children?: React.ReactNode;
}
export const SearchAtom: React.FC<SearchAtomProps> = ({}) => {
  return (
    <div className='search-bar'>
      <div className='search-input'>
        <Image w='24' url={searchIcon} />
        <input name='keyword' id='keyword' type='text' placeholder='Search here'/>
      </div>
    </div>
  );
};
