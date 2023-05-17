export interface DeviceData {
  device_id: string;
  device_name: string;
  firmware_ver: number;
  app_ver: number;
  type: string;
  status: string;
  lifetime: string;
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
  address: string;
  phone_num: string;
  email: string;
  type: TypeUserEnum;
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
