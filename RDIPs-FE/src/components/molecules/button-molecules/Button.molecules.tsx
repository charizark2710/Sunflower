import { ButtonGroup } from "@mui/material"
import { ButtonAtom, ButtonAtomProps } from "../../atoms/button/Button.atom"

interface Props {
    variant: any;
    buttons: Button[];
}

interface Button extends ButtonAtomProps {
    key: string;
  }

export const ButtonMolecules : React.FC<Props> = (props) => {
    return (
        <ButtonGroup variant={props.variant || 'text'}>
            {props.buttons.map((button)=> {
                return (
                    <ButtonAtom buttonSize={button.buttonSize} buttonStyle={button.buttonStyle}
                        onClick={button.onClick} type={button.type || 'button'} key={button.key}>
                        {button.children}
                    </ButtonAtom>
                );
            })}
        </ButtonGroup>
    )
}