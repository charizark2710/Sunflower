import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select, { SelectChangeEvent } from '@mui/material/Select';
import * as React from 'react';
import './Search.atom.scss';

export default function SearchAtom({ label = 'Select Option', className = '' }) {
  const [age, setAge] = React.useState('');

  const handleChange = (event: SelectChangeEvent) => {
    setAge(event.target.value);
  };
  return (
    <div className='search-container'>
      <FormControl variant='standard' sx={{ m: 1, minWidth: 120 }}>
        <InputLabel id='demo-simple-select-standard-label'>Search...</InputLabel>
        <Select id='demo-simple-select-standard' value={age} onChange={handleChange}>
          <MenuItem value=''>
            <em>None</em>
          </MenuItem>
          <MenuItem value={10}>Ten</MenuItem>
          <MenuItem value={20}>Twenty</MenuItem>
          <MenuItem value={30}>Thirty</MenuItem>
        </Select>
      </FormControl>
    </div>
  );
}
