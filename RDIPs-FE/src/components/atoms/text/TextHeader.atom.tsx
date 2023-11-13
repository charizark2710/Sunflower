import './TextHeader.atom.scss';

export interface TextAtomProps {
  text: string
}

const TextAtomHeader : React.FC<TextAtomProps> = ({ text }) => {
  return (<span className='header-text'>
    {text}
  </span>)
}

export default TextAtomHeader;