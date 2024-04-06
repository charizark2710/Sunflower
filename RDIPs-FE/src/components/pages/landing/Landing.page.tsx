import { Box, Container, CssBaseline, ThemeProvider, Typography, createTheme } from '@mui/material';
import config from '../../../utils/en.json';
import { ClutterButtonMolecules } from '../../molecules/clutter-button-mocules/ClutterButton.molecules';
import { CopyrightMolecules } from '../../molecules/copyright/Copyright.mocules';
import { useNavigate } from 'react-router-dom';

const defaultTheme = createTheme();

export default function LandingPage() {
  const navigate = useNavigate();
  const onLogin = () => {
    navigate('/login');
  };

  const onRegister = () => {
    navigate('/register');
  };
  return (
    <ThemeProvider theme={defaultTheme}>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          minHeight: '100vh',
        }}
      >
        <CssBaseline />
        <Container component='main' sx={{ mt: 8, mb: 2 }} maxWidth='md'>
          <Typography component={'h1'} variant='h2' gutterBottom>
            {config['landingPage.title']}
          </Typography>
          <Typography component={'h2'} variant='h5' gutterBottom>
            {config['landingPage.description']}
          </Typography>
          <ClutterButtonMolecules onLogin={onLogin} onRegister={onRegister}></ClutterButtonMolecules>
        </Container>
        <Box
          component='footer'
          sx={{
            py: 3,
            px: 2,
            mt: 'auto',
            backgroundColor: (theme) => (theme.palette.mode === 'light' ? theme.palette.grey[200] : theme.palette.grey[800]),
          }}
        >
          <Container maxWidth='sm'>
            <Typography component={'span'} variant='body1'>
              My sticky footer can be found here.
            </Typography>
            <CopyrightMolecules />
          </Container>
        </Box>
      </Box>
    </ThemeProvider>
  );
}
