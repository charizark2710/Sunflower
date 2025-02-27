export interface DeviceData {
  device_id: string;
  device_name: string;
  firmware_ver: number;
  app_ver: number;
  type: string;
  status: string;
  lifetime: string;
  region: string;
}

export interface DeviceLogHistoryData {
  datetime: string;
  status: StatusEnum;
  message: string;
}

export interface DeviceChangeHistoryData {
  datetime: string;
  type: string;
  description: string;
}

export interface UserData {
  user_id: string;
  user_name: string;
  firstName?: string;
  lastName?: string;
  address: string;
  phone_num: string;
  email: string;
  type: TypeUserEnum;
  emailVerified: boolean;
  enabled: boolean;
}

export interface UserResponse {
  id: string,
  username: string,
  address: string,
  phone_num: string,
  email: string,
  type: TypeUserEnum,
  emailVerified: boolean,
  enabled: boolean,
  firstName?: string;
  lastName?: string;
}

export interface AdminData {
  admin_id: string;
  admin_name: string;
  status: string;
  role: string;
  auhentication: string;
}

export enum TypeUserEnum {
  Regular = 'regular',
  Industrial = 'industrial',
}

export interface ReceiptData {
  receive_id: string;
  datetime: string;
  component: string;
  service: string;
  total: number;
}

export enum StatusEnum {
  Default = 'default',
  Error = 'error',
  Warning = 'warning',
  Fatal = 'fatal',
}

export interface HeadCell {
  id?: any;
  label: string;
  numeric: boolean | any;
}


export interface RequestCreateDevices {
  name: string;
  type?: string;
  status?: string;
  id?: string;
  region?: string;
}

export interface RequestCreateUsers {
  firstName: string;
  lastName: string;
  email: string;
  username?: string;
  id?: string;
}
