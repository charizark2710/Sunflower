import { HighchartsReact } from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import React, { useEffect, useMemo } from 'react';
import { DateFilter, TypeChart } from '../../utils/enum';

export interface HighChartsProps {
  timeType: DateFilter;
  typeChart: TypeChart;
  chartData: DataDetail[];
  titleChart: string | null;
}

export interface DataDetail {
  AO: number;
  AC: number;
  AL: number;
  dateTime: string;
  unit: string;
}

export const HighChartCustom: React.FC<HighChartsProps> = (props : HighChartsProps) => {
  let amountDisplay :number = 0;// distance between them
  useMemo(()=>{

    switch(props.timeType) {
      case DateFilter.Hour: {
        console.log("in 1", props.chartData);
        amountDisplay = 24
        break;
      }
      case DateFilter.Date: {
        console.log("in 2");

        break;
      }
      case DateFilter.Week: {
        console.log("in 3");

        break;
      }
      case DateFilter.Month: {
        console.log("in 4");

        break;
      }
      case DateFilter.Year: {
        console.log("in 5");

        break;
      }
    }
  }, [props, amountDisplay])

  function getOptions(typeChart: TypeChart) {
    switch(typeChart) {
      case 1: {
        return optionsSplineChart;
      }
      case 2: {
        return optionsBarChart;
      }
    }
  }
  const monthData = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  const hourData = [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24]
  
  function filterDataBy(type : DateFilter, data: number) {

    return null;
  }

  const optionsSplineChart = {
    chart: {
      type: 'spline'
    },
    title: {
      text: props.titleChart
    },
    xAxis: {
      categories: hourData,
      title: {
        text: `Time (${DateFilter[props.timeType]})`
      }
    },
    yAxis: {
      title: {
        text: `A (${props.chartData[0].unit})`
      }
    },
    plotOptions: {
      line: {
        dataLabels: {
          enabled: true
        },
        enableMouseTracking: false
      }
    },
    series: [
      {
        name: 'Power consumption',
        data: props.chartData.map(data => data.AC)
      },
      {
        name: 'Power obtained',
        data: props.chartData.map(data => data.AO)
      },
      {
        name: 'Power loss',
        data: props.chartData.map(data => data.AL)
      }
    ]
  };

  const optionsBarChart = {
    chart: {
      type: 'column'
    },
    title: {
      text: props.titleChart
    },
    subtitle: {
      text: 'Source: RDIPs.com'
    },
    xAxis: {
      categories: hourData,
      title: {
        text: `Time (${DateFilter[props.timeType]})`
      },
      crosshair: true
    },
    yAxis: {
      min: 0,
      title: {
        text: `A (${props.chartData[0].unit})`
      }
    },
    tooltip: {
      headerFormat: '<span style="font-size:10px">{point.key}</span><table>',
      pointFormat: '<tr><td style="color:{series.color};padding:0">{series.name}: </td>' +
        '<td style="padding:0"><b>{point.y:.1f} mm</b></td></tr>',
      footerFormat: '</table>',
      shared: true,
      useHTML: true
    },
    plotOptions: {
      column: {
        pointPadding: 0.2,
        borderWidth: 0
      }
    },
    series: [
      {
        name: 'Power consumption',
        data: props.chartData.map(data => data.AC)
      },
      {
        name: 'Power obtained',
        data: props.chartData.map(data => data.AO)
      },
      {
        name: 'Power loss',
        data: props.chartData.map(data => data.AL)
      }
    ]
  };


  function highChartsRender(ops: any) {
    Highcharts.chart({
      chart: {
        type: 'pie',
        renderTo: 'abc',
      },
      title: {
        verticalAlign: 'middle',
        floating: true,
        text: "Earth's Atmospheric Composition",
        style: {
          fontSize: '10px',
        },
      },
      plotOptions: {
        pie: {
          dataLabels: {
            format: '{point.name}: {point.percentage:.1f} %',
          },
          innerSize: '70%',
        },
      },
      series: ops,
    });
  }

  useEffect(()=>{

  },[])

  return (
    <div className='highchart-container'>
      <HighchartsReact
        highcharts={Highcharts}
        options={getOptions(props.typeChart)}
      />
    </div>
  );
};
