import React from 'react';
import { ButtonAtom } from '../../atoms/button/Button.atom';
import './ClutterButton.molecules.scss';

interface ClutterButtonMoleculesProps {
  onLogin?: (args: any) => void;
  onRegister?: (args: any) => void;
}

export const ClutterButtonMolecules: React.FC<ClutterButtonMoleculesProps> = (props) => {
  const { onLogin, onRegister } = props;
  return (
    <div style={{marginRight: '10px', display: 'flex', alignItems: 'center'}}>
      <ButtonAtom
        buttonStyle='btn--login'
        buttonSize='btn--small'
        onClick={onLogin}
      >
        Log In
      </ButtonAtom>
      <ButtonAtom
        buttonStyle='btn--register'
        buttonSize='btn--small'
        onClick={onRegister}
      >
        Register
      </ButtonAtom>
    </div>
  );
};
