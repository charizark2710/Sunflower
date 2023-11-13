import { Box, Button, TextField } from '@mui/material';
import { Field } from 'formik';
import React from 'react';
import { addDevice, updateDevice } from '../../../../axios/api';
import ErrorMessageAtom from '../../../atoms/error-message/ErrorMessageAtom.atom';
import { FormikAtom } from '../../../atoms/formik/FormikAtom.atom';
import { DeviceData } from '../../../../utils/interface';

interface FormCreateDeviceMoleculesProps {
  onClosePopUp: (value?: any) => void;
  state?: string;
  data?: DeviceData;
}
export interface RequestCreateDevices {
  name: string;
  type?: string;
  status?: string;
  id?: string;
  region?: string;
}

export const FormCreateDeviceMolecules: React.FC<FormCreateDeviceMoleculesProps> = ({
  onClosePopUp,
  state = 'create',
  data,
}) => {
  const initialValues: any = {
    name: data?.device_name || '',
    type: data?.type || '',
    status: data?.status || '',
    id: data?.device_id || '',
  };

  const handleSubmit = (formValue: RequestCreateDevices) => {
    if (state === 'create') {
      addDevice(formValue)
        .then((data: any) => onClosePopUp())
        .catch((err: any) => alert(err));
    } else {
      updateDevice(formValue)
        .then((data: any) => onClosePopUp())
        .catch((err: any) => alert(err));
    }
  };

  const createDeviceForm = (
    <Box sx={{ mt: 1 }}>
      <Field
        as={TextField}
        required
        fullWidth
        id='name'
        label='Device Name'
        name='name'
        autoComplete='name'
        autoFocus
        helperText={<ErrorMessageAtom name='name' />}
      />

      {state === 'create' ? (
        <>
          <Field
            as={TextField}
            margin='normal'
            variant='outlined'
            fullWidth
            name='type'
            label='Type'
            type='text'
            required={true}
            helperText={<ErrorMessageAtom name='type' />}
          />
          <Field
            as={TextField}
            margin='normal'
            variant='outlined'
            fullWidth
            name='region'
            label='Region'
            type='text'
            required={true}
            helperText={<ErrorMessageAtom name='region' />}
          />
        </>
      ) : (
        <Field
          as={TextField}
          margin='normal'
          variant='outlined'
          fullWidth
          name='status'
          label='Status'
          type='text'
          required={true}
          helperText={<ErrorMessageAtom name='status' />}
        />
      )}

      <Button className='add-button' type='submit' fullWidth variant='contained' sx={{ mt: 3, mb: 2 }}>
        {state === 'create' ? 'Add' : 'Update'}
      </Button>
    </Box>
  );

  return <FormikAtom initialValues={initialValues} onSubmit={handleSubmit} children={createDeviceForm}></FormikAtom>;
};
