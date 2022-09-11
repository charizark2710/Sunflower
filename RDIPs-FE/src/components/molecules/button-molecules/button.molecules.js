import { ButtonGroup } from "@mui/material"
import { ButtonAtom } from "../../atoms/button/Button.atom"

export const ButtonMolecules = (props) => {
    return (
        <ButtonGroup variant={props.variant || 'text'}>
            {props.buttons.map(button => {
                return (
                    <ButtonAtom buttonSize={button.size} buttonStyle={button.style}
                        onClick={button.onClick} type={button.type || 'button'} key={button.key}>
                        {button.children}
                    </ButtonAtom>
                );
            })}
        </ButtonGroup>
    )
}