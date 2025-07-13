import { Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { BreadcrumbComponent } from '../shared/breadcrumb/breadcrumb.component';
import { DeviceActivityStatisticComponent } from '../shared/device-activity-statistic/device-activity-statistic.component';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { DataTableComponent, TableColumn } from '../shared/data-table/data-table.component';
import { TranslateModule } from '@ngx-translate/core';
import { TextGroupComponent } from '../shared/text-group/text-group.component';

interface DeviceResponse {
  id: string;
  name: string;
  firmware_ver: number;
  app_ver: number;
  type: string;
  status: string;
  life_time: string;
  region: string;
}

interface DeviceData {
  device_id: string;
  device_name: string;
  firmware_ver: number;
  app_ver: number;
  type: string;
  status: string;
  lifetime: string;
  region: string;
  user_name?: string;
}

@Component({
  selector: 'app-devices',
  imports: [
    BreadcrumbComponent,
    NavigationComponent,
    DeviceActivityStatisticComponent,
    DataTableComponent,
    RouterModule,
    TextGroupComponent,
    TranslateModule
  ],
  templateUrl: './devices.component.html',
  styleUrl: './devices.component.scss',
})
export class DevicesComponent implements OnInit {
  deviceListData: DeviceData[] = [];
  
  // Generic table configuration
  tableColumns: TableColumn[] = [
    { key: 'name', label: 'labels.table.columns.device', type: 'device-info' },
    { key: 'progress', label: 'labels.table.columns.progress', type: 'progress' },
    { key: 'user_name', label: 'labels.table.columns.username', type: 'text' },
    { key: 'region', label: 'labels.table.columns.region', type: 'text' },
    { key: 'lifetime', label: 'labels.table.columns.lifetime', type: 'text' },
    { key: 'status', label: 'labels.table.columns.status', type: 'status' },
    { key: 'actions', label: 'labels.table.columns.actions', type: 'actions' }
  ];

  ngOnInit() {
    this.getListDevice();
  }

  createDeviceData(data: DeviceResponse): DeviceData {
    const { id, name, firmware_ver, app_ver, type, status, life_time, region } = data;
    return {
      device_id: id,
      device_name: name,
      firmware_ver,
      app_ver,
      type,
      status,
      lifetime: life_time,
      region,
      user_name: `${name.toLowerCase()}.device` // Generate username for display
    };
  }

  // Mock API call - replace with actual HTTP service
  getListDevice() {
    const mockApiResponse: DeviceResponse[] = [
      {
        id: '1',
        name: 'Ephemeral',
        firmware_ver: 2.1,
        app_ver: 1.0,
        type: 'IoT Sensor',
        status: 'Active',
        life_time: '36s',
        region: 'Europe'
      },
      {
        id: '2', 
        name: 'Stack3d Lab',
        firmware_ver: 2.0,
        app_ver: 1.2,
        type: 'Gateway',
        status: 'Sleep',
        life_time: '1w',
        region: 'South America'
      },
      {
        id: '3',
        name: 'Warpspeed',
        firmware_ver: 1.9,
        app_ver: 1.1,
        type: 'Controller',
        status: 'Warning',
        life_time: '6d',
        region: 'Africa'
      },
      {
        id: '4',
        name: 'CloudWatch',
        firmware_ver: 2.2,
        app_ver: 1.3,
        type: 'Monitor',
        status: 'Error',
        life_time: '1d',
        region: 'Oceania'
      },
      {
        id: '5',
        name: 'ContrastAI',
        firmware_ver: 1.8,
        app_ver: 0.9,
        type: 'AI Device',
        status: 'Critical',
        life_time: '38s',
        region: 'North America'
      }
    ];

    this.deviceListData = mockApiResponse
      .reverse()
      .map((device: DeviceResponse) => this.createDeviceData(device));
  }

  navigateToDetailPage(deviceDetail: DeviceData) {
    console.log('Navigate to device detail:', deviceDetail);
    // TODO: Implement Angular router navigation
    // this.router.navigate(['/detail-device'], { state: deviceDetail });
  }
}
