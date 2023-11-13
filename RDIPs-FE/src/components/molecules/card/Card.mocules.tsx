import './Card.mocules.scss';
import { connect } from 'react-redux';
import TextAtomHeader from '../../atoms/text/TextHeader.atom';
import EditIcon from '@mui/icons-material/Edit';
import { SimpleDialog } from '../../atoms/dialog/Dialog.atom';
import { useEffect, useState } from 'react';

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
    <section>
      <SimpleDialog title={addTitle} children={modal} open={open} onClose={handleClose} />
      <div className='card-container'>
        <div className='mb-10 flex-justify-space-between'>
          <TextAtomHeader text={title} />
          <EditIcon onClick={handleClickOpen} />
        </div>
        <div className='card-body'>{children}</div>
      </div>
    </section>
  );
};

export default connect()(CardMocules);
