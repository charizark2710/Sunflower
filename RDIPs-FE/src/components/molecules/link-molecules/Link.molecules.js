import React from 'react';
import { LinkAtom } from '../../atoms/link/Link.atom';
import './Link.molecules.scss'

export const LinkMolecules = ({ links, isClick, children }) => {
    return (
        <div className={isClick ? 'nav-menu active' : 'nav-menu'}>
            {links.map((link) => {
                return (
                    <div className={'nav-item'}>
                        <LinkAtom
                            to={link.to}
                            onClick={link.onClick}
                            className={link.className}
                        >
                            {link.children}
                        </LinkAtom>
                    </div>)
            })}
            {children}
        </div>
    )
};