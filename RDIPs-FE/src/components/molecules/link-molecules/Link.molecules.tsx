import React from 'react';
import { LinkAtom, LinkAtomProps } from '../../atoms/link/Link.atom';
import './Link.molecules.scss'

interface LinkMoleculesProps {
    links?: Link[];
    isClick?: boolean;
    target?: string;
    children?: React.ReactNode;
    to?: string;
    className?:string;
}

interface Link extends LinkAtomProps {
    key?: string;
  }
export const LinkMolecules: React.FC<LinkMoleculesProps> = ({ links = [], isClick, children }) => {
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