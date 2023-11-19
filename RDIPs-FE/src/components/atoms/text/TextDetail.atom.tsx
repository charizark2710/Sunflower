import './TextHeader.atom.scss';

export interface TextDetailAtomProps {
  title: string;
  children?: React.ReactNode;
}

const TextAtomDetail: React.FC<TextDetailAtomProps> = ({ title, children }) => {
  return (
    <div className='text-detail'>
      <span>{title}</span>
      <div>{children}</div>
    </div>
  );
};

export default TextAtomDetail;
