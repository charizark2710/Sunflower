import { connect } from 'react-redux';
import './Breakcrumb.mocules.scss';
import { ButtonAtom } from '../../atoms/button/Button.atom';

interface BreakCrumbProps {
  title: string;
  icon?: React.ReactNode;
}

const defaultStyle = {
  fontWeight: '400',
  fontSize: '12px',
  lineHeight: '15px',
  padding: '10px',
};

const defaultLink = { to: '/', className: 'link-item', children: 'Home' };

const BreakCrumbMocule: React.FC<BreakCrumbProps> = ({
  title,
  icon = <></>,
}) => {
  return (
    <div className='breakcrumb-container'>
      <div className='breakcrumb-big flex-justify-space-between'>
        {`${title}`}
        <ButtonAtom children={`Add ${title}`}/>
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
