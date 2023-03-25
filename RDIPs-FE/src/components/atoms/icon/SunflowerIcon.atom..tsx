import ArrowBackIosIcon from '@mui/icons-material/ArrowBackIos';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';

const SunFlowerIcon = ({ fontSize = 16 }) => {
  return (
    <>
      <ArrowBackIosIcon style={{ fontSize: fontSize }} />
      <ArrowForwardIosIcon style={{ fontSize: fontSize }} />
    </>
  );
};

export default SunFlowerIcon;
