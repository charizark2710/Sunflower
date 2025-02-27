import { Box, Button, TextField } from '@mui/material';
import { Field } from 'formik';
import React from 'react';
import { addUser, updateUser } from '../../../../axios/api';
import { initValue } from '../../../../utils/function';
import { UserData } from '../../../../utils/interface';
import ErrorMessageAtom from '../../../atoms/error-message/ErrorMessageAtom.atom';
import { FormikAtom } from '../../../atoms/formik/FormikAtom.atom';

interface FormCreateUserMoleculesProps {
  onClosePopUp: (value?: any) => void;
  state?: string;
  data?: UserData;
}
export interface RequestCreateUsers {
  firstName: string;
  lastName: string;
  email: string;
  username: string;
  id: string;
}

export const FormCreateUserMolecules: React.FC<
  FormCreateUserMoleculesProps
> = ({ onClosePopUp, state = 'create', data }) => {
  const initialValues: any = {
    firstName: initValue(data?.firstName),
    lastName: initValue(data?.lastName),
    email: initValue(data?.email),
    username: initValue(data?.user_name),
    id: initValue(data?.user_id),
  };

  const handleSubmit = (formValue: RequestCreateUsers) => {
    if (state === 'create') {
      addUser(formValue)
        .then((_) => onClosePopUp())
        .catch((_) => alert('Oop! Can not create new user.'));
    } else {
      updateUser(formValue)
        .then((_) => onClosePopUp())
        .catch((_) => alert('Oop! Can not update new user.'));
    }
  };

  const createUserForm = (
    <Box sx={{ mt: 1 }}>
      {state === 'create' ? (
        <>
          <Field
            as={TextField}
            required
            fullWidth
            autoFocus
            id='username'
            label='User Name'
            name='username'
            autoComplete='username'
            helperText={<ErrorMessageAtom name='username' />}
          />
        </>
      ) : (
        ''
      )}

      <Field
        as={TextField}
        margin='normal'
        variant='outlined'
        fullWidth
        name='firstName'
        label='First name'
        type='text'
        required={true}
        helperText={<ErrorMessageAtom name='firstName' />}
      />

      <Field
        as={TextField}
        margin='normal'
        variant='outlined'
        fullWidth
        name='lastName'
        label='Last name'
        type='text'
        required={true}
        helperText={<ErrorMessageAtom name='lastName' />}
      />

      <Field
        as={TextField}
        required
        fullWidth
        id='email'
        label='Email'
        name='email'
        autoComplete='email'
        helperText={<ErrorMessageAtom name='email' />}
      />

      <Button
        className='add-button'
        type='submit'
        fullWidth
        variant='contained'
        sx={{ mt: 3, mb: 2 }}
      >
        {state === 'create' ? 'Add' : 'Update'}
      </Button>
    </Box>
  );

  return (
    <FormikAtom
      initialValues={initialValues}
      onSubmit={handleSubmit}
      children={createUserForm}
    ></FormikAtom>
  );
};
