import { Component } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { BarChartComponent } from '../bar-chart/bar-chart.component';

@Component({
  selector: 'app-device-activity-statistic',
  standalone: true,
  imports: [BarChartComponent, TranslateModule],
  templateUrl: './device-activity-statistic.component.html',
  styleUrl: './device-activity-statistic.component.scss',
})
export class DeviceActivityStatisticComponent {}
