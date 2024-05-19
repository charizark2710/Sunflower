import { Box, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { ButtonAtom } from '../../atoms/button/Button.atom';
import { SimpleDialog } from '../../atoms/dialog/Dialog.atom';
import { SearchAtom } from '../../atoms/search/Search.atom';
import TitlePageAtom from '../../atoms/text/TitlePgae.atom';
import './Breakcrumb.mocules.scss';

interface BreakCrumbProps {
  title: string;
  modal?: React.ReactNode;
  icon?: React.ReactNode;
  status?: string;
  level?: 1 | 2 | 3 | 4 | 5 | 6 | 7;
  link?: string;
}

const BreakCrumbMocule: React.FC<BreakCrumbProps> = ({ title, icon = <></>, modal, status, level = 2, link = '/admin' }) => {
  const navigate = useNavigate();
  let addTitle = `Add ${title}`;
  let path = level === 2 ? 'Home' : 'Home/...';

  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (status === 'closed') {
      handleClose();
    }
  }, [status]);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = (value?: string) => {
    setOpen(false);
  };

  const redirectToHome = () => {
    navigate('/admin');
  };

  const redirectToLink = () => {
    navigate(link, { replace: true });
  };

  return (
    <Box className='breakcrumb-container'>
      <SimpleDialog title={addTitle} children={modal} open={open} onClose={handleClose} />
      <Box className='flex-justify-space-between'>
        <Box>
          <Box className='breakcrumb-small'>
            <Typography component={'span'} className='underline' onClick={redirectToHome}>
              {path}
            </Typography>
            {'/'}
            <Typography component={'span'} className='underline' onClick={redirectToLink}>
              {icon} {title}
            </Typography>
          </Box>
          <TitlePageAtom title={title}></TitlePageAtom>
        </Box>
        <Box className='flex-align-center w-max'>
          <SearchAtom />
          <ButtonAtom children={'Add New'} buttonStyle='add-button' onClick={handleClickOpen} />
        </Box>
      </Box>
    </Box>
  );
};

const mapPropToState = ({ navbarTitle }: any) => {
  return {
    navbarTitle,
  };
};

export default connect(mapPropToState)(BreakCrumbMocule);
