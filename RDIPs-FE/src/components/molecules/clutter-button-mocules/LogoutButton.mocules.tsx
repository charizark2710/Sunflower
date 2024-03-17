import { LogoutIcon } from '../../atoms/icon/ListIcon.atom';
import './ClutterButton.molecules.scss';

export const LogoutButtonMolecules = () => {
  return (
    <div className='logout-area'>
      <LogoutIcon />
      <span> Logout </span>
    </div>
  );
};
