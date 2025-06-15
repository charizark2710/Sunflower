import { NgFor } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { Dashboard_Data } from '../constant';
import { BreadcrumbComponent } from '../shared/breadcrumb/breadcrumb.component';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { DataCardComponent } from './data-card/data-card.component';
import { BarChartComponent } from '../shared/bar-chart/bar-chart.component';
import { TranslateModule } from '@ngx-translate/core';
import { DeviceActivityStatisticComponent } from '../shared/device-activity-statistic/device-activity-statistic.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    RouterModule,
    NavigationComponent,
    DataCardComponent,
    NgFor,
    BreadcrumbComponent,
    DeviceActivityStatisticComponent,
    TranslateModule
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent {
  CARD_DATA = Dashboard_Data;

  trackByFn(index: number, item: any): any {
    return item.id; // Assuming each item has a unique 'id' property
  }
}
