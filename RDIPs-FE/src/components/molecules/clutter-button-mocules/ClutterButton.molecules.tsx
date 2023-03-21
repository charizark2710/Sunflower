import React from 'react';
import { ButtonAtom } from '../../atoms/button/Button.atom';
import './ClutterButton.molecules.scss';

interface ClutterButtonMoleculesProps {}

export const ClutterButtonMolecules: React.FC<ClutterButtonMoleculesProps> = () => {
  return (
    <div style={{marginRight: '10px', display: 'flex', alignItems: 'center'}}>
      <ButtonAtom
        buttonStyle='btn--register'
        buttonSize='btn--small'
      >
        Log In
      </ButtonAtom>
      <ButtonAtom
        buttonStyle='btn--login'
        buttonSize='btn--small'
      >
        Register
      </ButtonAtom>
    </div>
  );
};
