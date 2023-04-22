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
  datetime: string,
  status: StatusEnum,
  message: string
}

export interface DeviceChangeHistoryData {
  datetime: string,
  type: string,
  description: string
}

export enum StatusEnum {
  Default = 'default',
  Error = 'error',
  Warning = 'warning',
  Fatal = 'fatal'
}

export interface HeadCell {
  disablePadding: boolean;
  id?: any;
  label: string;
  numeric: boolean;
}