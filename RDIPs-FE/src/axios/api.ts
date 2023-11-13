import { RequestCreateDevices } from '../components/molecules/form/device-create/FormCreateDevice.molecules';
import { axiosClient } from './axiosClient';

export const getAllDevices = (id ?: string) => {
  if(id) {
    return axiosClient.get(`devices/${id}/?detail=true`);
  }
  return axiosClient.get('devices');
};

export const addDevice = (request: RequestCreateDevices) => {
  return axiosClient.post('devices', {
    name: request.name,
    type: request.type,
    region: request.region
  });
};

export const updateDevice = (request: RequestCreateDevices) => {
  return axiosClient.put(`devices/${request.id}`, {
    name: request.name,
    status: request.status
  });
};
