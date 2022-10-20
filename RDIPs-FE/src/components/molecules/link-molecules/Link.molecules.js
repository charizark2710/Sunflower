import React from 'react';
import { LinkAtom } from '../../atoms/link/Link.atom';
import './Link.molecules.scss'

export const LinkMolecules = ({ links, isClick, children }) => {
    return (
        <div className={isClick ? 'link-menu active' : 'link-menu'}>
            {links.map((link) => {
                return (
                    <div key={link.to} style={{
                        display: 'flex',
                        width: `calc(100%/${links.length})`,
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}>
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
        </div >
    )
};