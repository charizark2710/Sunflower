import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { BreadcrumbComponent } from '../shared/breadcrumb/breadcrumb.component';
import { DataTableComponent, TableColumn, TableData } from '../shared/data-table/data-table.component';

interface AdminData {
  admin_id: string;
  admin_name: string;
  status: string;
  role: string;
  authentication: string;
  history_active?: string;
}

@Component({
  selector: 'app-admins',
  imports: [RouterModule, NavigationComponent, BreadcrumbComponent, DataTableComponent],
  templateUrl: './admins.component.html',
  styleUrl: './admins.component.scss'
})
export class AdminsComponent {
  tableColumns: TableColumn[] = [
    { key: 'stt', label: 'labels.table.columns.stt', type: 'text', sortable: false },
    { key: 'admin_name', label: 'labels.table.columns.adminName', type: 'text', sortable: true },
    { key: 'status', label: 'labels.table.columns.status', type: 'status', sortable: true },
    { key: 'role', label: 'labels.table.columns.role', type: 'text', sortable: true },
    { key: 'authentication', label: 'labels.table.columns.authentication', type: 'text', sortable: true },
    { key: 'history_active', label: 'labels.table.columns.historyActive', type: 'text', sortable: true },
    { key: 'actions', label: 'labels.table.columns.actions', type: 'actions', sortable: false }
  ];

  adminListData: TableData[] = [
    {
      stt: 1,
      admin_id: 'AD001',
      admin_name: 'Esther Howard',
      status: 'Active',
      role: 'Admin',
      authentication: 'Admin',
      history_active: '15 minutes ago'
    },
    {
      stt: 2,
      admin_id: 'AD002',
      admin_name: 'Marvin McKinney',
      status: 'Active',
      role: 'Admin',
      authentication: 'Admin',
      history_active: '1 month ago'
    },
    {
      stt: 3,
      admin_id: 'AD003',
      admin_name: 'Jerome Bell',
      status: 'Active',
      role: 'Admin',
      authentication: 'Admin',
      history_active: '5 months ago'
    },
    {
      stt: 4,
      admin_id: 'AD004',
      admin_name: 'Albert Flores',
      status: 'Active',
      role: 'Admin',
      authentication: 'Admin',
      history_active: '2 week ago'
    },
    {
      stt: 5,
      admin_id: 'AD005',
      admin_name: 'Darlene Robertson',
      status: 'Active',
      role: 'Admin',
      authentication: 'Admin',
      history_active: '2 months ago'
    },
    {
      stt: 6,
      admin_id: 'AD006',
      admin_name: 'Kathryn Murphy',
      status: 'Active',
      role: 'Admin',
      authentication: 'Admin',
      history_active: '5 hours ago'
    },
    {
      stt: 7,
      admin_id: 'AD007',
      admin_name: 'Cody Fisher',
      status: 'Active',
      role: 'Admin',
      authentication: 'Admin',
      history_active: '15 minutes ago'
    }
  ];

  onEdit(admin: any): void {
    console.log('Edit admin:', admin);
    // Implement edit functionality
  }

  onDelete(admin: any): void {
    console.log('Delete admin:', admin);
    // Implement delete functionality
  }
}
