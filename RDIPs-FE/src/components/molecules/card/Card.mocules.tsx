import EditIcon from '@mui/icons-material/Edit';
import { Box, Divider } from '@mui/material';
import { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { SimpleDialog } from '../../atoms/dialog/Dialog.atom';
import TextAtomHeader from '../../atoms/text/TextHeader.atom';
import './Card.mocules.scss';

interface CardProps {
  title: string;
  children?: React.ReactNode;
  modal?: React.ReactNode;
  status?: string;
}

const CardMocules: React.FC<CardProps> = ({ title, children, modal, status }) => {
  const addTitle = 'Update ' + title;
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if(status === "closed") {
      handleClose();
    }
  }, [status]);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };
  return (
    <Divider>
      <SimpleDialog title={addTitle} children={modal} open={open} onClose={handleClose} />
      <Box className='card-container'>
        <Box className='mb-10 flex-justify-space-between'>
          <TextAtomHeader text={title} />
          <EditIcon onClick={handleClickOpen} />
        </Box>
        <Box className='card-body'>{children}</Box>
      </Box>
    </Divider>
  );
};

export default connect()(CardMocules);
