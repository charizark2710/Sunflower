import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Avatar from '@mui/material/Avatar';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Checkbox from '@mui/material/Checkbox';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import FormControlLabel from '@mui/material/FormControlLabel';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import { Field } from 'formik';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router';
import * as Yup from 'yup';
import { login } from '../../../redux/actions/authentication';
import ErrorMessageAtom from '../../atoms/error-message/ErrorMessageAtom.atom';
import { FormikAtom } from '../../atoms/formik/FormikAtom.atom';
import { CopyrightMolecules } from '../../molecules/copyright/Copyright.mocules';

const defaultTheme = createTheme();

const LoginPage = ({ dispatch }: any) => {
  const navigate = useNavigate();

  const initialValues: any = {
    email: 'admin@gmail.com',
    password: '1234',
  };

  const LoginValidationSchema = Yup.object().shape({
    email: Yup.string().email('Please enter the right email format').required('Email Required'),
    password: Yup.string().required('Password required'),
  });

  const USER = {
    name: 'admin@gmail.com',
    password: '1234',
  };

  const handleLogin = (value: any) => {
    let { email, password } = value;
    if (email === USER.name && password === USER.password) {
      dispatch(login(true));
      navigate('/admin');
    } else {
      alert('Wrong email or password');
    }
  };

  const loginForm = (
    <Box sx={{ mt: 1 }}>
      <Field
        as={TextField}
        required
        fullWidth
        id='email'
        label='Email Address'
        name='email'
        autoComplete='email'
        autoFocus
        helperText={<ErrorMessageAtom name='email' />}
      />

      <Field
        as={TextField}
        margin='normal'
        variant='outlined'
        fullWidth
        name='password'
        label='Password'
        type='password'
        id='password'
        required={true}
        autoComplete='current-password'
        helperText={<ErrorMessageAtom name='password' />}
      />
      <FormControlLabel control={<Checkbox value='remember' color='primary' />} label='Remember me' />
      <Button type='submit' fullWidth variant='contained' sx={{ mt: 3, mb: 2 }}>
        Sign In
      </Button>
      <Grid container>
        <Grid item xs>
          <Link href='#' variant='body2'>
            Forgot password?
          </Link>
        </Grid>
        <Grid item>
          <Link href='/register' variant='body2'>
            {"Don't have an account? Sign Up"}
          </Link>
        </Grid>
      </Grid>
    </Box>
  );

  return (
    <ThemeProvider theme={defaultTheme}>
      <Container component='main' maxWidth='xs'>
        <CssBaseline />
        <Box
          sx={{
            marginTop: 8,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
          }}
        >
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component='h1' variant='h5'>
            Sign in
          </Typography>
          <FormikAtom
            initialValues={initialValues}
            validationSchema={LoginValidationSchema}
            onSubmit={handleLogin}
            children={loginForm}
          ></FormikAtom>
        </Box>
        <CopyrightMolecules />
      </Container>
    </ThemeProvider>
  );
};

export default connect()(LoginPage);
