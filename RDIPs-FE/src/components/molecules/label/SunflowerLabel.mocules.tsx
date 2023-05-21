import AcUnitIcon from '@mui/icons-material/AcUnit';
import SunFlowerIcon from '../../atoms/icon/SunflowerIcon.atom';
import { LinkAtom, LinkAtomProps } from '../../atoms/link/Link.atom';
import './SunflowerLabel.mocules.scss';

interface SunLabelProps {
  link?: LinkAtomProps;
  icon?: string;
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
  onClick?: (args: any) => void;
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
  icon = '',
  style = defaultStyle,
  iconPos = 0,
  height = '37px',
  size = 'md',
  children,
  onClick,
}) => {
  function iconItem() {
    switch (icon) {
      case 'toggle':
        return <SunFlowerIcon />;
      default:
        return <AcUnitIcon fontSize='inherit' />;
    }
  }
  return (
    <>
      <div style={style} onClick={onClick}>
        <div className='flex-align-center flex-justify-center side-bar-item' style={{ height: height }}>
          <div>
            {size === 'md' ? (
              <>
                <LinkAtom to={link.to} className={link.className}>
                  {iconPos === 0 ? iconItem() : ''}
                  &nbsp;
                  {link.children}
                </LinkAtom>
              </>
            ) : (
              <>
                <LinkAtom to={link.to} className={link.className}>
                  {iconPos === 0 ? iconItem() : ''}
                  &nbsp;
                </LinkAtom>
              </>
            )}
          </div>
        </div>
      </div>
      {size === 'md' && children}
    </>
  );
};

export default SunflowerLabel;
