import { Visibility, VisibilityOff } from '@mui/icons-material';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import CssBaseline from '@mui/material/CssBaseline';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { Field, Form, Formik } from 'formik';
import { useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router';
import * as Yup from 'yup';
import { login } from '../../../redux/actions/authentication';
import ErrorMessageAtom from '../../atoms/error-message/ErrorMessageAtom.atom';
import SunFlowerLogo from '../../atoms/image/SunFlowerLogo.atom';
import './LoginPage.page.scss';

const defaultTheme = createTheme();

const LoginPage = ({ dispatch }: any) => {
  const [showPassword, setShowPassword] = useState(false);

  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };

  const navigate = useNavigate();

  const initialValues: any = {
    email: 'admin@gmail.com',
    password: '1234',
  };

  const LoginValidationSchema = Yup.object().shape({
    email: Yup.string().email('Wrong email format').required('Email Required'),
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
          <SunFlowerLogo />
          <div className='administrator-title'>Administrator</div>

          <Formik
            initialValues={initialValues}
            validationSchema={LoginValidationSchema}
            onSubmit={handleLogin}
          >
            {({}) => (
              <Form>
                <Box sx={{ mt: 1 }}>
                  <div className='flex-justify-start flex-align-center'>
                    <span className='input-title'>Email</span>
                    <span>
                      <Field
                        as={TextField}
                        required
                        fullWidth
                        id='email'
                        name='email'
                        placeholder='Your email'
                        autoComplete='email'
                        helperText={<ErrorMessageAtom name='email' />}
                      />
                    </span>
                  </div>
                  <div className='flex-justify-start flex-align-center'>
                    <span className='input-title'>Password</span>
                    <span>
                      <Field
                        as={TextField}
                        margin='normal'
                        variant='outlined'
                        fullWidth
                        name='password'
                        type={showPassword ? 'text' : 'password'}
                        placeholder='Your password'
                        id='password'
                        required={true}
                        autoComplete='current-password'
                        helperText={<ErrorMessageAtom name='password' />}
                      />
                    </span>
                    {!showPassword ? (
                      <Visibility
                        onClick={handleClickShowPassword}
                        className='ml-12'
                      />
                    ) : (
                      <VisibilityOff
                        onClick={handleClickShowPassword}
                        className='ml-12'
                      />
                    )}
                  </div>
                  <div className='group-button flex-justify-center'>
                    <button type='submit' className='login-button'>
                      Login
                    </button>
                    <Link href='#' variant='body2' className='forgot-pass'>
                      <span>Forgot your password?</span>
                    </Link>
                  </div>
                </Box>
              </Form>
            )}
          </Formik>
        </Box>
      </Container>
    </ThemeProvider>
  );
};

export default connect()(LoginPage);
