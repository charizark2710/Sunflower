import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { DevicesComponent } from './devices/devices.component';
import { UsersComponent } from './users/users.component';
import { AdminsComponent } from './admins/admins.component';
import { CampaignComponent } from './campaign/campaign.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'devices', component: DevicesComponent },
  { path: 'users', component: UsersComponent },
  { path: 'admins', component: AdminsComponent },
  { path: 'campaign', component: CampaignComponent }
];
