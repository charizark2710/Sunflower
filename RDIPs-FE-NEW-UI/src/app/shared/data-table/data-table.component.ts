import { Component, Input, Output, EventEmitter } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { CommonModule } from '@angular/common';

export interface TableColumn {
  key: string;
  label: string;
  type?: 'text' | 'status' | 'progress' | 'actions' | 'device-info';
  sortable?: boolean;
}

export interface TableData {
  [key: string]: any;
}

const ELEMENT_DATA: TableData[] = [
  { 
    id: 1, 
    name: 'Ephemeral', 
    email: 'ephemeral.io', 
    region: 'Europe', 
    lifetime: '36s', 
    status: 'Active', 
    progress: 60,
    type: 'Regular'
  },
  { 
    id: 2, 
    name: 'Stack3d Lab', 
    email: 'stack3dlab.com', 
    region: 'South America', 
    lifetime: '1w', 
    status: 'Sleep', 
    progress: 72,
    type: 'Industrial'
  },
  { 
    id: 3, 
    name: 'Warpspeed', 
    email: 'getwarpspeed...', 
    region: 'Africa', 
    lifetime: '6d', 
    status: 'Warning', 
    progress: 78,
    type: 'Common'
  },
  { 
    id: 4, 
    name: 'CloudWatch', 
    email: 'cloudwatch.vip', 
    region: 'Oceania', 
    lifetime: '1d', 
    status: 'Error', 
    progress: 38,
    type: 'Common'
  },
  { 
    id: 5, 
    name: 'ContrastAI', 
    email: 'contrastai.com', 
    region: 'North America', 
    lifetime: '38s', 
    status: 'Critical', 
    progress: 42,
    type: 'Common'
  }
]
@Component({
  selector: 'app-data-table',
  standalone: true,
  imports: [TranslateModule, CommonModule],
  templateUrl: './data-table.component.html',
  styleUrl: './data-table.component.scss',
})
export class DataTableComponent {
  @Input() columns: TableColumn[] = [
    { key: 'name', label: 'labels.table.columns.device', type: 'device-info' },
    { key: 'progress', label: 'labels.table.columns.progress', type: 'progress' },
    { key: 'email', label: 'labels.table.columns.username', type: 'text' },
    { key: 'region', label: 'labels.table.columns.region', type: 'text' },
    { key: 'lifetime', label: 'labels.table.columns.lifetime', type: 'text' },
    { key: 'status', label: 'labels.table.columns.status', type: 'status' },
    { key: 'actions', label: 'labels.table.columns.actions', type: 'actions' }
  ];
  @Input() dataSource: TableData[] = ELEMENT_DATA;
  @Input() showDeviceInfo: boolean = true; // Controls device info display
  @Input() showProgress: boolean = true; // Controls progress bar display
  @Input() showActions: boolean = true; // Controls action buttons

  @Output() onEdit = new EventEmitter<TableData>();
  @Output() onDelete = new EventEmitter<TableData>();

  getStatusClass(status: number | string | undefined): string {
    if (!status) return 'active'; // Handle undefined/null
    
    if (typeof status === 'string') {
      switch (status.toLowerCase()) {
        case 'active': return 'active';
        case 'sleep': return 'sleep';
        case 'warning': return 'warning';
        case 'error': return 'error';
        case 'critical': return 'critical';
        default: return 'active';
      }
    }
    
    switch (status) {
      case 1: return 'active';
      case 2: return 'sleep';
      case 3: return 'warning';
      case 4: return 'error';
      case 5: return 'critical';
      default: return 'active';
    }
  }

  getStatusText(status: number | string | undefined): string {
    if (!status) return 'Active'; // Handle undefined/null
    
    if (typeof status === 'string') {
      return status.charAt(0).toUpperCase() + status.slice(1);
    }
    
    switch (status) {
      case 1: return 'Active';
      case 2: return 'Sleep';
      case 3: return 'Warning';
      case 4: return 'Error';
      case 5: return 'Critical';
      default: return 'Active';
    }
  }

  getProgressClass(progress: number | undefined): string {
    if (!progress) return 'medium'; // Handle undefined
    if (progress >= 70) return 'high';
    if (progress >= 40) return 'medium';
    return 'low';
  }

  // Action handlers
  editItem(item: TableData): void {
    this.onEdit.emit(item);
  }

  deleteItem(item: TableData): void {
    this.onDelete.emit(item);
  }
}
