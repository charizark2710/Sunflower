import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { BreadcrumbComponent } from '../shared/breadcrumb/breadcrumb.component';
import { DeviceActivityStatisticComponent } from '../shared/device-activity-statistic/device-activity-statistic.component';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { DataTableComponent } from '../shared/data-table/data-table.component';

@Component({
  selector: 'app-devices',
  imports: [
    BreadcrumbComponent,
    NavigationComponent,
    DeviceActivityStatisticComponent,
    DataTableComponent,
    RouterModule,
  ],
  templateUrl: './devices.component.html',
  styleUrl: './devices.component.scss',
})
export class DevicesComponent {}
