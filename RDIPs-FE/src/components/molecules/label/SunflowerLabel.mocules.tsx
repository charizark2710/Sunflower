import AcUnitIcon from '@mui/icons-material/AcUnit';
import { Box, Typography } from '@mui/material';
import { connect } from 'react-redux';
import { Page } from '../../../model/page';
import { LinkAtom, LinkAtomProps } from '../../atoms/link/Link.atom';
import './SunflowerLabel.mocules.scss';

interface SunLabelProps {
  link?: LinkAtomProps;
  isHomepage?: boolean;
  icon?: string;
  key: any;
  style?: {
    fontWeight?: string;
    fontSize: string;
    lineHeight: string;
    height?: string;
  };
  iconPos?: number; // iconPos = 0 : icon before label name
  height?: string;
  size?: string; //md: medium- show full sidebar and sm: small- only show icon
  children?: React.ReactNode;
  specialIcon?: React.ReactNode;
  state?: boolean; // true = collapse, false = expand
  onClick?: (args: any) => void;
  navbarTitle?: string;
}

const defaultStyle = {
  fontWeight: '400',
  fontSize: '16px',
  lineHeight: '24px',
  padding: '4px 16px',
  letterSpacing: '0.15px',
};

const defaultLink = { to: '/', className: 'link-item', children: 'Home' };

const SunflowerLabel: React.FC<SunLabelProps> = ({
  link = defaultLink,
  isHomepage = false,
  specialIcon = <AcUnitIcon />,
  style = defaultStyle,
  iconPos = 0,
  height = '32px',
  state,
  navbarTitle,
  onClick,
}) => {
  function iconItem() {
    return <>{specialIcon}</>;
  }

  return (
    <LinkAtom to={link?.to} className={link?.className}>
      <Box style={style} onClick={onClick}>
        <Box className={state ? 'side-bar-item flex-justify-center' : 'side-bar-item'} style={{ height: height }}>
          <Box>
            <Box
              className={
                isHomepage
                  ? 'flex-align-center home-page-item'
                  : link != null && link?.children === navbarTitle
                  ? 'flex-align-center bg-sidebar-item active'
                  : 'flex-align-center  bg-sidebar-item'
              }
            >
              <Box
                className='flex-align-center'
                style={{
                  height: '30px',
                  width: 'fix-content',
                  padding: iconPos !== 0 ? '0' : '0 10px',
                }}
              >
                {iconPos === 0 ? iconItem() : ''}{' '}
              </Box>
              {link != null ? (
                <Typography component={'span'} className={link?.children === navbarTitle ? 'active link-text' : 'link-text'}>
                  {state ? '' : link?.children}
                </Typography>
              ) : (
                ''
              )}
            </Box>
          </Box>
        </Box>
      </Box>
    </LinkAtom>
  );
};

const mapPropToState = ({ page }: { page: Page }) => {
  return {
    navbarTitle: page.navbarTitle,
  };
};

export default connect(mapPropToState)(SunflowerLabel);
