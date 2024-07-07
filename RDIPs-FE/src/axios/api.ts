import { RequestCreateDevices, RequestCreateUsers } from '../utils/interface';
import { axiosClient } from './axiosClient';

export const getAllDevices = (id ?: string) => {
  if(id) {
    return axiosClient.get(`devices/${id}/?detail=true`);
  }
  return axiosClient.get('devices');
};

export const getAllUsers = (id ?: string) => {
  if(id) {
    return axiosClient.get(`users/${id}/?detail=true`);
  }
  return axiosClient.get('users');
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

export const addUser = (request: RequestCreateUsers) => {
  return axiosClient.post('users', {
    username: request.username,
    firstName: request.firstName,
    lastName: request.lastName,
    email: request.email,
    enabled: true
  });
};

export const updateUser = (request: RequestCreateUsers) => {
  return axiosClient.put(`users/${request.id}`, {
    firstName: request.firstName,
    lastName: request.lastName,
    email: request.email,
    enabled: true
  });
};
