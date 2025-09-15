import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { BreadcrumbComponent } from '../shared/breadcrumb/breadcrumb.component';
import { DataTableComponent, TableColumn, TableData } from '../shared/data-table/data-table.component';

export enum TypeUserEnum {
  Regular = 'Regular',
  Industrial = 'Industrial',
  Common = 'Common',
}

interface UserData {
  user_id: string;
  user_name: string;
  firstName?: string;
  lastName?: string;
  address: string;
  phone_num: string;
  email: string;
  type: TypeUserEnum;
  emailVerified: boolean;
  enabled: boolean;
  totalDevices?: number;
  status?: string;
}

@Component({
  selector: 'app-users',
  imports: [RouterModule, NavigationComponent, BreadcrumbComponent, DataTableComponent],
  templateUrl: './users.component.html',
  styleUrl: './users.component.scss'
})
export class UsersComponent {
  tableColumns: TableColumn[] = [
    { key: 'stt', label: 'labels.table.columns.stt', type: 'text', sortable: false },
    { key: 'user_name', label: 'labels.table.columns.userName', type: 'text', sortable: true },
    { key: 'address', label: 'labels.table.columns.address', type: 'text', sortable: true },
    { key: 'phone_num', label: 'labels.table.columns.phoneNumber', type: 'text', sortable: true },
    { key: 'email', label: 'labels.table.columns.email', type: 'text', sortable: true },
    { key: 'type', label: 'labels.table.columns.type', type: 'text', sortable: true },
    { key: 'totalDevices', label: 'labels.table.columns.totalDevices', type: 'text', sortable: true },
    { key: 'status', label: 'labels.table.columns.status', type: 'status', sortable: true },
    { key: 'actions', label: 'labels.table.columns.actions', type: 'actions', sortable: false }
  ];

  userListData: TableData[] = [
    {
      stt: 1,
      user_id: 'U001',
      user_name: 'Esther Howard',
      address: '8080 Belmont St',
      phone_num: '(208) 555-0112',
      email: 'martha...',
      type: 'Regular',
      totalDevices: 12,
      status: 'Active'
    },
    {
      stt: 2,
      user_id: 'U002',
      user_name: 'Marvin McKinney',
      address: '7529 E. Pecan St',
      phone_num: '(219) 555-0114',
      email: 'okctmt...',
      type: 'Industrial',
      totalDevices: 34,
      status: 'Active'
    },
    {
      stt: 3,
      user_id: 'U003',
      user_name: 'Jerome Bell',
      address: '775 Rolling Gr...',
      phone_num: '(208) 555-0111',
      email: 'trethtu...',
      type: 'Common',
      totalDevices: 24,
      status: 'Active'
    },
    {
      stt: 4,
      user_id: 'U004',
      user_name: 'Albert Flores',
      address: '8558 Green Rd.',
      phone_num: '(702) 555-0122',
      email: 'trungkle...',
      type: 'Common',
      totalDevices: 17,
      status: 'Active'
    },
    {
      stt: 5,
      user_id: 'U005',
      user_name: 'Darlene Robertson',
      address: '3890 Poplar Dr.',
      phone_num: '(208) 555-0100',
      email: 'thuhang...',
      type: 'Common',
      totalDevices: 45,
      status: 'Active'
    },
    {
      stt: 6,
      user_id: 'U006',
      user_name: 'Kathryn Murphy',
      address: '3890 Poplar Dr.',
      phone_num: '(406) 555-0120',
      email: 'vuhath...',
      type: 'Common',
      totalDevices: 120,
      status: 'Active'
    },
    {
      stt: 7,
      user_id: 'U007',
      user_name: 'Cody Fisher',
      address: '3605 Parker Rd.',
      phone_num: '(319) 555-0115',
      email: 'dangho...',
      type: 'Common',
      totalDevices: 100,
      status: 'Active'
    }
  ];

  onEdit(user: any): void {
    console.log('Edit user:', user);
    // Implement edit functionality
  }

  onDelete(user: any): void {
    console.log('Delete user:', user);
    // Implement delete functionality
  }
}
