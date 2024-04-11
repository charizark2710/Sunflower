import { Visibility, VisibilityOff } from '@mui/icons-material';
import { Field, Form, Formik } from 'formik';
import { useState } from 'react';
import { connect } from 'react-redux';
import ErrorMessageAtom from '../../atoms/error-message/ErrorMessageAtom.atom';
import SunFlowerLogo from '../../atoms/image/SunFlowerLogo.atom';
import './LoginPage.page.scss';
import { login } from '../../../redux/slice/authSlice';
import { Box, Button, Container, CssBaseline, Link, TextField, ThemeProvider, Typography, createTheme } from '@mui/material';
import { AccountInfo } from '../../../model/page';
import { useNavigate } from 'react-router-dom';
import * as Yup from 'yup';

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
    email: Yup.string(),
    password: Yup.string().required('Password required'),
  });

  const handleLogin = (value: any) => {
    let { email, password } = value;
    dispatch(login({ username: email, password: password }))
      .unwrap()
      .then(() => {
        navigate('/admin');
      })
      .catch(() => {
        alert('Wrong email or password');
      });
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
          <Box className='administrator-title'>Administrator</Box>

          <Formik initialValues={initialValues} validationSchema={LoginValidationSchema} onSubmit={handleLogin}>
            {({ errors, isValid }: { errors: AccountInfo; isValid: boolean }) => (
              <Form>
                <Box sx={{ mt: 1 }}>
                  <Box className='flex-justify-start flex-align-center'>
                    <Typography component={'span'} className='input-title'>
                      Username
                    </Typography>
                    <Typography component={'span'}>
                      <Field
                        as={TextField}
                        required
                        fullWidth
                        id='email'
                        name='email'
                        placeholder='Your email'
                        autoComplete='email'
                        className={errors.email ? 'error' : ''}
                        helperText={<ErrorMessageAtom name='email' />}
                      />
                    </Typography>
                  </Box>
                  <Box className='flex-justify-start flex-align-center'>
                    <Typography component={'span'} className='input-title'>
                      Password
                    </Typography>
                    <Typography component={'span'}>
                      <Field
                        as={TextField}
                        margin='normal'
                        variant='outlined'
                        fullWidth
                        name='password'
                        type={showPassword ? 'text' : 'password'}
                        placeholder='Your password'
                        id='password'
                        className={errors.password ? 'error' : ''}
                        required={true}
                        autoComplete='current-password'
                        helperText={<ErrorMessageAtom name='password' />}
                      />
                    </Typography>
                    {!showPassword ? (
                      <Visibility onClick={handleClickShowPassword} className='ml-12' />
                    ) : (
                      <VisibilityOff onClick={handleClickShowPassword} className='ml-12' />
                    )}
                  </Box>
                  <Box className='group-button flex-justify-center'>
                    <Button type='submit' className={isValid ? 'enabled login-button' : 'login-button'}>
                      Login
                    </Button>
                    <Link href='#' variant='body2' className='forgot-pass'>
                      <Typography component={'span'}>Forgot your password?</Typography>
                    </Link>
                  </Box>
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
