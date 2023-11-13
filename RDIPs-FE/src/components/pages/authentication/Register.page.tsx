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
import ErrorMessageAtom from '../../atoms/error-message/ErrorMessageAtom.atom';
import { FormikAtom } from '../../atoms/formik/FormikAtom.atom';
import { CopyrightMolecules } from '../../molecules/copyright/Copyright.mocules';
import { Field } from 'formik';
import * as Yup from 'yup';

const theme = createTheme();

export default function Register() {
  const initialValues: any = {
    firstName: '',
    lastName: '',
    email: 'admin@gmail.com',
    password: '1234',
  };

  const RegisterValidationSchema = Yup.object().shape({
    firstName: Yup.string().required('First name required'),
    lastName: Yup.string().required('Last name required'),
    email: Yup.string().email('Please enter the right email format').required('Email Required'),
    password: Yup.string().required('Password required'),
  });

  const handeRegister = (value: any) => {
    console.log('registered');
  };

  const registerForm = (
    <Box sx={{ mt: 3 }}>
      <Grid container spacing={2}>
        <Grid item xs={12} sm={6}>
          <Field
            as={TextField}
            autoComplete='given-name'
            name='firstName'
            required
            fullWidth
            id='firstName'
            label='First Name'
            helperText={<ErrorMessageAtom name='firstName' />}
            autoFocus
          />
        </Grid>
        <Grid item xs={12} sm={6}>
          <Field
            as={TextField}
            required
            fullWidth
            id='lastName'
            label='Last Name'
            name='lastName'
            helperText={<ErrorMessageAtom name='lastName' />}
            autoComplete='family-name'
          />
        </Grid>
        <Grid item xs={12}>
          <Field
            as={TextField}
            required
            fullWidth
            id='email'
            label='Email Address'
            name='email'
            helperText={<ErrorMessageAtom name='email' />}
            autoComplete='email'
          />
        </Grid>
        <Grid item xs={12}>
          <Field
            as={TextField}
            required
            fullWidth
            name='password'
            label='Password'
            type='password'
            id='password'
            helperText={<ErrorMessageAtom name='password' />}
            autoComplete='new-password'
          />
        </Grid>
        <Grid item xs={12}>
          <FormControlLabel
            control={<Checkbox value='allowExtraEmails' color='primary' />}
            label='I want to receive inspiration, marketing promotions and updates via email.'
          />
        </Grid>
      </Grid>
      <Button type='submit' fullWidth variant='contained' sx={{ mt: 3, mb: 2 }}>
        Sign Up
      </Button>
      <Grid container justifyContent='flex-end'>
        <Grid item>
          <Link href='/login' variant='body2'>
            Already have an account? Sign in
          </Link>
        </Grid>
      </Grid>
    </Box>
  );

  return (
    <ThemeProvider theme={theme}>
      <Container component='main' maxWidth='sm'>
        <CssBaseline />
        <Box
          sx={{
            marginTop: 8,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            boxShadow: 3,
            borderRadius: 2,
            px: 4,
            py: 6,
          }}
        >
          <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
            <LockOutlinedIcon />
          </Avatar>
          <Typography component='h1' variant='h5'>
            Sign up
          </Typography>
          <FormikAtom
            initialValues={initialValues}
            validationSchema={RegisterValidationSchema}
            onSubmit={handeRegister}
            children={registerForm}
          ></FormikAtom>
        </Box>
        <CopyrightMolecules />
      </Container>
    </ThemeProvider>
  );
}
