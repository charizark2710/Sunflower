import AcUnitIcon from '@mui/icons-material/AcUnit';
import SunFlowerIcon from '../../atoms/icon/SunflowerIcon.atom';
import { LinkAtom, LinkAtomProps } from '../../atoms/link/Link.atom';
import './SunflowerLabel.mocules.scss';
import { connect } from 'react-redux';

interface SunLabelProps {
  link?: LinkAtomProps;
  icon?: string;
  key: any;
  style?: {
    fontWeight?: string;
    fontSize: string;
    lineHeight: string;
    height?: string;
  };
  iconPos?: number; // iconPos = 0 : icon before label name
  height?: string;
  size?: string; //md: medium- show full sidebar and sm: small- only show icon
  children?: React.ReactNode;
  specialIcon?: React.ReactNode;
  onClick?: (args: any) => void;
  navbarTitle?: string;
}

const defaultStyle = {
  fontWeight: '400',
  fontSize: '12px',
  lineHeight: '15px',
  padding: '10px',
};

const defaultLink = { to: '/', className: 'link-item', children: 'Home' };

const SunflowerLabel: React.FC<SunLabelProps> = ({
  link = defaultLink,
  key = '',
  icon = '',
  specialIcon = <AcUnitIcon />,
  style = defaultStyle,
  iconPos = 0,
  height = '37px',
  size = 'md',
  navbarTitle,
  children,
  onClick,
}) => {
  function iconItem() {
    switch (icon) {
      case 'toggle':
        return <SunFlowerIcon />;
      default:
        return <>{specialIcon}</>;
    }
  }

  return (
    <>
      <div style={style} onClick={onClick}>
        <div className='side-bar-item' style={{ height: height }}>
          <div>
            {size === 'md' ? (
              <div className='flex-align-center'>
                <div className='flex-align-center' style={{ height: '30px', width: 'fix-content', paddingLeft: iconPos !== 0 ? '20px': '10px' }}>
                  {iconPos === 0 ? iconItem() : ''}{' '}
                </div>
                <LinkAtom to={link.to} className={link.className}>
                  &nbsp;
                  <span className={link.children === navbarTitle ? 'active link-text' : 'link-text'}>{link.children}</span>
                </LinkAtom>
              </div>
            ) : (
              <>
                <div className='flex-align-center flex-justify-center'>
                  <LinkAtom to={link.to} className={link.className}>
                    <div className='flex-align-center flex-justify-center' style={{ height: '30px', width: 'fix-content' }}>
                      {iconPos === 0 ? iconItem() : ''}{' '}
                    </div>
                    &nbsp;
                  </LinkAtom>
                </div>
              </>
            )}
          </div>
        </div>
      </div>
      {size === 'md' && children}
    </>
  );
};

const mapPropToState = ({ navbarTitle }: any) => {
  return {
    navbarTitle,
  };
};

export default connect(mapPropToState)(SunflowerLabel);
