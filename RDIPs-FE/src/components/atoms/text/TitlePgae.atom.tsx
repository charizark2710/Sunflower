import { Box } from '@mui/material';
import './TitlePage.atom.scss';

export interface TitlePageProps {
  title: string;
}

const TitlePageAtom: React.FC<TitlePageProps> = ({ title }) => {
  return <Box className='breakcrumb-big'>{title}</Box>;
};

export default TitlePageAtom ;
