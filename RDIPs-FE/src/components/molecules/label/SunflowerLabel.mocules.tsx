import AcUnitIcon from '@mui/icons-material/AcUnit';
import SunFlowerIcon from '../../atoms/icon/SunflowerIcon.atom';
import { StraightAtom } from '../../atoms/straight/Straight.atom';
import './SunflowerLabel.mocules.scss';

interface SunLabelProps {
  labelName: React.ReactNode;
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

const SunflowerLabel: React.FC<SunLabelProps> = ({
  labelName = 'List Devices',
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
        return <AcUnitIcon fontSize='inherit' />;
    }
  }
  return (
    <div style={style} onClick={onClick}>
      <div className='flex-align-center flex-justify-center' style={{height: height}}>
        <div>
          {iconPos === 0 ? iconItem() : ''}
          &nbsp;
          {size==='md' ? labelName : ""}
          &nbsp;
          {iconPos === 1 ? iconItem() : ''}
        </div>
      </div>

      <StraightAtom width='80%' thick='0.1px' color='#3E3E3E' />
    </div>
  );
};

export default SunflowerLabel;
