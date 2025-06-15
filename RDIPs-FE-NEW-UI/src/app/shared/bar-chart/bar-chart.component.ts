import { Component, Input, OnInit } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import Chart from 'chart.js/auto';

@Component({
  selector: 'app-bar-chart',
  imports: [],
  templateUrl: './bar-chart.component.html',
  styleUrl: './bar-chart.component.scss',
})
export class BarChartComponent implements OnInit {
  chart: any;
  @Input() titleChart: string = 'labels.chart.title';
  @Input() legendChart: string = 'labels.chart.legend';

  constructor(private translateService: TranslateService) {}

  ngOnInit(): void {
    this.createChart();
  }

  createChart() {
    const titleChartConverted = this.translateService.instant(this.titleChart);
    const legendChartConverted = this.translateService.instant(
      this.legendChart
    );

    this.chart = new Chart('MyChart', {
      type: 'bar', //this denotes tha type of chart
      data: {
        // values on X-Axis
        labels: [
          '2022-05-10',
          '2022-05-11',
          '2022-05-12',
          '2022-05-13',
          '2022-05-14',
          '2022-05-15',
          '2022-05-16',
          '2022-05-17',
          '2022-05-14',
          '2022-05-15',
          '2022-05-16',
          '2022-05-17',
        ],
        datasets: [
          {
            label: 'data1',
            data: [
              '467',
              '576',
              '572',
              '79',
              '92',
              '574',
              '573',
              '576',
              '467',
              '576',
              '572',
              '79',
            ],
            backgroundColor: 'rgba(214, 187, 251, 1)',
            borderRadius: 5
          },
          {
            label: 'data2',
            data: [
              '542',
              '542',
              '536',
              '327',
              '17',
              '0.00',
              '538',
              '541',
              '17',
              '0.00',
              '538',
              '541',
            ],
            backgroundColor: ' rgba(158, 119, 237, 1)',
            borderRadius: 5
          },
          {
            label: 'data3',
            data: [
              '542',
              '542',
              '536',
              '327',
              '17',
              '0.00',
              '538',
              '541',
              '536',
              '327',
              '17',
              '0.00',
            ],
            backgroundColor: 'rgba(51, 55, 65, 1)',
            borderRadius: 5
          },
        ],
      },
      options: {
        plugins: {
          title: {
            display: true,
            text: titleChartConverted,
            position: 'left',
          },
          legend: {
            display: true,
            position: 'bottom',
            labels: {
              generateLabels: () => {
                return [
                  {
                    text: legendChartConverted,
                    fontColor: '#1F242F',
                    fillStyle: 'transparent',
                    strokeStyle: 'transparent',
                    lineWidth: 0,
                  },
                ];
              },
            },
          },
        },
        aspectRatio: 2.5,
        scales: {
          x: {
            stacked: true,
          },
          y: {
            stacked: true,
          },
        },
      },
    });

    Chart.defaults.backgroundColor = '#9BD0F5';
    Chart.defaults.borderColor = '#1F242F';
    Chart.defaults.color = '#1F242F';
  }
}
