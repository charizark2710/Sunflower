import { Component } from '@angular/core';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-campaign',
  imports: [RouterModule, NavigationComponent],
  templateUrl: './campaign.component.html',
  styleUrl: './campaign.component.scss'
})
export class CampaignComponent {

}
