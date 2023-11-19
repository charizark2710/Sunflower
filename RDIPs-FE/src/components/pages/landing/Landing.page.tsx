import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import Typography from '@mui/material/Typography';
import { useNavigate } from 'react-router';
import config from '../../../utils/en.json';
import { ClutterButtonMolecules } from '../../molecules/clutter-button-mocules/ClutterButton.molecules';
import { CopyrightMolecules } from '../../molecules/copyright/Copyright.mocules';

const defaultTheme = createTheme();

export default function LandingPage() {
  const navigate = useNavigate();
  const onLogin = () => {
    navigate('/login');
  }

  const onRegister = () => {
    navigate('/register');
  }
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
        <Container component="main" sx={{ mt: 8, mb: 2 }} maxWidth="md">
          <Typography variant="h2" component="h1" gutterBottom>
            {config['landingPage.title']}
          </Typography>
          <Typography variant="h5" component="h2" gutterBottom>
            {config['landingPage.description']}
          </Typography>
          <ClutterButtonMolecules onLogin={onLogin} onRegister={onRegister}></ClutterButtonMolecules>
        </Container>
        <Box
          component="footer"
          sx={{
            py: 3,
            px: 2,
            mt: 'auto',
            backgroundColor: (theme) =>
              theme.palette.mode === 'light'
                ? theme.palette.grey[200]
                : theme.palette.grey[800],
          }}
        >
          <Container maxWidth="sm">
            <Typography variant="body1">
              My sticky footer can be found here.
            </Typography>
            <CopyrightMolecules />
          </Container>
        </Box>
      </Box>
    </ThemeProvider>
  );
}