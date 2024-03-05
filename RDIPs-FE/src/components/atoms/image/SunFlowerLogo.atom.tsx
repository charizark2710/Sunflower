import './SunFlowerLogo.atom.scss';
import Image from './Image';
import Sunflower from '../../../assets/Sunflower.svg';
import SunflowerLabel from '../../../assets/Sunflower-label.svg';

const SunFlowerLogo = ({ label = true, w="150px" }) => {
  return (
    <div className="sunflower-logo">
      <Image url={Sunflower} w={w}/>
      { label ?
        <Image url={SunflowerLabel} w={w}/>  : ''
      }
    </div>
  );
};

export default SunFlowerLogo;