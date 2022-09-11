import React from 'react';
import './Button.atom.scss';
import { Button } from '@mui/material';

const STYLES = ['btn--primary', 'btn--outline'];

const SIZES = ['btn--medium', 'btn--large'];

export const ButtonAtom = ({ type, onClick, buttonStyle, buttonSize, children }) => {
    const checkButtonStyle = STYLES.includes(buttonStyle)
        ? buttonStyle
        : STYLES[0];

    const checkButtonSize = SIZES.includes(buttonSize) ? buttonSize : SIZES[0];
    return (
        <Button
            className={`btn ${checkButtonStyle} ${checkButtonSize}`}
            onClick={onClick}
            type={type}
        >
            {children}
        </Button>
    )
};