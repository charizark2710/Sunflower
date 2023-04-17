import AcUnitIcon from '@mui/icons-material/AcUnit';
import SunFlowerIcon from '../../atoms/icon/SunflowerIcon.atom';
import { StraightAtom } from '../../atoms/straight/Straight.atom';
import './SunflowerLabel.mocules.scss';
import { LinkAtom, LinkAtomProps } from '../../atoms/link/Link.atom';

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
  size?: string //md: medium- show full sidebar and sm: small- only show icon
  onClick?: (args: any) => void;
}

const defaultStyle = {
  fontWeight: '400',
  fontSize: '12px',
  lineHeight: '15px',
};

const defaultLink = { to: '/', className: 'link-item', children: 'Home' };

const SunflowerLabel: React.FC<SunLabelProps> = ({
  link = defaultLink,
  icon = '',
  style = defaultStyle,
  iconPos = 0,
  height = '37px',
  size="md",
  onClick
}) => {
  function iconItem() {
    switch (icon) {
      case 'toggle':
        return <SunFlowerIcon />;
      default:
        return <AcUnitIcon fontSize='inherit'/>;
    }
  }
  return (
    <div style={style} >
      <div className='flex-align-center flex-justify-center side-bar-item' style={{height: height}}>
        <div>
          {iconPos === 0 ? iconItem() : ''}
          &nbsp;
          {size==='md' ?
          <LinkAtom to={link.to} className={link.className}>
              {link.children}
          </LinkAtom> : ""}
          &nbsp;
          <span onClick={onClick}> {iconPos === 1 ? iconItem() : ''}</span>
        </div>
      </div>

      <StraightAtom width='80%' thick='0.1px' color='#3E3E3E' />
    </div>
  );
};

export default SunflowerLabel;
