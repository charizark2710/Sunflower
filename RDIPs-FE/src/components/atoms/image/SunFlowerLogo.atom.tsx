import './SunFlowerLogo.atom.scss';
import Image from './Image';
import Sunflower from '../../../assets/Sunflower.svg';
import SunflowerLabel from '../../../assets/Sunflower-label.svg';
import { Box } from '@mui/system';

const SunFlowerLogo = ({ label = true}) => {
  return (
    <div className="sunflower-logo">
      <Image url={Sunflower} w="150px"/>
      { label ?
        <Image url={SunflowerLabel} w="150px"/>  : ''
      }
    </div>
  );
};

export default SunFlowerLogo;