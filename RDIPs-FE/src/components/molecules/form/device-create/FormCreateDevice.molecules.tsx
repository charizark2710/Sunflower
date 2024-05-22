import { Box, Button, TextField } from '@mui/material';
import { Field } from 'formik';
import React from 'react';
import { addDevice, updateDevice } from '../../../../axios/api';
import { initValue } from '../../../../utils/function';
import { DeviceData } from '../../../../utils/interface';
import ErrorMessageAtom from '../../../atoms/error-message/ErrorMessageAtom.atom';
import { FormikAtom } from '../../../atoms/formik/FormikAtom.atom';

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
    name: initValue(data?.device_name),
    type: initValue(data?.type),
    status: initValue(data?.status),
    region: initValue(data?.region),
    id: initValue(data?.device_id)
  };

  const handleSubmit = (formValue: RequestCreateDevices) => {
    if (state === 'create') {
      addDevice(formValue)
        .then((_) => onClosePopUp())
        .catch((_) => alert('Oop! There is some error happened!'));
    } else {
      updateDevice(formValue)
        .then((_) => onClosePopUp())
        .catch((_) => alert('Oop! There is some error happened!'));
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
