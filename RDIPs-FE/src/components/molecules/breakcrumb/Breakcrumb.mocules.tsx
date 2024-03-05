import { connect } from 'react-redux';
import './Breakcrumb.mocules.scss';
import { ButtonAtom } from '../../atoms/button/Button.atom';
import { SimpleDialog } from '../../atoms/dialog/Dialog.atom';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { SearchAtom } from '../../atoms/search/Search.atom';

interface BreakCrumbProps {
  title: string;
  modal?: React.ReactNode;
  icon?: React.ReactNode;
  status?: string;
  level?: 1 | 2 | 3 | 4 | 5 | 6 | 7;
  link?: string;
}

const BreakCrumbMocule: React.FC<BreakCrumbProps> = ({
  title,
  icon = <></>,
  modal,
  status,
  level = 2,
  link = '/admin',
}) => {
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
    <div className='breakcrumb-container'>
      <SimpleDialog
        title={addTitle}
        children={modal}
        open={open}
        onClose={handleClose}
      />
      <div className='flex-justify-space-between'>
        <div>
          <div className='breakcrumb-small'>
            <span className='underline' onClick={redirectToHome}>
              {path}
            </span>
            {'/'}
            <span className='underline' onClick={redirectToLink}>
              {icon} {title}
            </span>
          </div>
          <div className='breakcrumb-big'>{title}</div>
        </div>
        <div className='flex-align-center w-max'>
          <SearchAtom />
          <ButtonAtom
            children={'Add New'}
            buttonStyle='add-button'
            onClick={handleClickOpen}
          />
        </div>
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