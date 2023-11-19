import { connect } from 'react-redux';
import './Breakcrumb.mocules.scss';
import { ButtonAtom } from '../../atoms/button/Button.atom';
import { SimpleDialog } from '../../atoms/dialog/Dialog.atom';
import { useEffect, useState } from 'react';

interface BreakCrumbProps {
  title: string;
  modal?: React.ReactNode;
  icon?: React.ReactNode;
  status?: string;
}

const BreakCrumbMocule: React.FC<BreakCrumbProps> = ({ title, icon = <></>, modal, status }) => {
  let addTitle = `Add ${title}`;
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if(status === "closed") {
      handleClose();
    }
  }, [status]);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = (value?: string) => {
    setOpen(false);
  };

  return (
    <div className='breakcrumb-container'>
      <SimpleDialog title={addTitle} children={modal} open={open} onClose={handleClose} />
      <div className='breakcrumb-big flex-justify-space-between'>
        {title}
        <ButtonAtom children={'+ ' + addTitle} buttonStyle='add-button' onClick={handleClickOpen} />
      </div>
      <div className='breakcrumb-small'>
        <span className='font-grey'>Home / ... /</span>{' '}
        <span>
          {icon} {title}
        </span>
      </div>
    </div>
  );
};

const mapPropToState = ({ navbarTitle }: any) => {
  return {
    navbarTitle,
  };
};

export default connect(mapPropToState)(BreakCrumbMocule);
