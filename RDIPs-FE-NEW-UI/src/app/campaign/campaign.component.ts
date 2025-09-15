import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { BreadcrumbComponent } from '../shared/breadcrumb/breadcrumb.component';
import { DataTableComponent, TableColumn, TableData } from '../shared/data-table/data-table.component';

interface CampaignData {
  campaign_id: string;
  member_user: string;
  plan_fee: string;
  element: boolean;
  family: boolean;
  package: boolean;
}

@Component({
  selector: 'app-campaign',
  imports: [RouterModule, NavigationComponent, BreadcrumbComponent, DataTableComponent],
  templateUrl: './campaign.component.html',
  styleUrl: './campaign.component.scss'
})
export class CampaignComponent {
  tableColumns: TableColumn[] = [
    { key: 'stt', label: 'labels.table.columns.stt', type: 'text', sortable: false },
    { key: 'member_user', label: 'labels.table.columns.memberUser', type: 'text', sortable: true },
    { key: 'plan_fee', label: 'labels.table.columns.planFee', type: 'text', sortable: true },
    { key: 'element', label: 'labels.table.columns.element', type: 'text', sortable: true },
    { key: 'family', label: 'labels.table.columns.family', type: 'text', sortable: true },
    { key: 'package', label: 'labels.table.columns.package', type: 'text', sortable: true },
    { key: 'actions', label: 'labels.table.columns.actions', type: 'actions', sortable: false }
  ];

  campaignListData: TableData[] = [
    {
      stt: 1,
      campaign_id: 'C001',
      member_user: 'Sự Kiện Khánh Thành Tháp Thiết',
      plan_fee: 'Regular',
      element: true,
      family: false,
      package: false
    },
    {
      stt: 2,
      campaign_id: 'C002',
      member_user: 'Chiến Dịch "Mua 1 Tặng 1"',
      plan_fee: 'Industrial',
      element: true,
      family: true,
      package: false
    },
    {
      stt: 3,
      campaign_id: 'C003',
      member_user: 'Sale Bão Hè 2024',
      plan_fee: 'Common',
      element: false,
      family: false,
      package: false
    },
    {
      stt: 4,
      campaign_id: 'C004',
      member_user: 'Ưu Đãi Cuối Năm',
      plan_fee: 'Common',
      element: false,
      family: false,
      package: true
    },
    {
      stt: 5,
      campaign_id: 'C005',
      member_user: 'Tháng Sáu Lớn Nhất Năm',
      plan_fee: 'Common',
      element: false,
      family: false,
      package: false
    },
    {
      stt: 6,
      campaign_id: 'C006',
      member_user: 'Ưu Tiệc Giảm Giá Cuối Tuần',
      plan_fee: 'Common',
      element: true,
      family: true,
      package: true
    },
    {
      stt: 7,
      campaign_id: 'C007',
      member_user: 'Ngày Độc Thân Giảm Cho Bạn',
      plan_fee: 'Common',
      element: false,
      family: false,
      package: true
    }
  ];

  onEdit(campaign: any): void {
    console.log('Edit campaign:', campaign);
    // Implement edit functionality
  }

  onDelete(campaign: any): void {
    console.log('Delete campaign:', campaign);
    // Implement delete functionality
  }
}
