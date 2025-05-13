import { Component } from '@angular/core';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-devices',
  imports: [RouterModule, NavigationComponent],
  templateUrl: './devices.component.html',
  styleUrl: './devices.component.scss'
})
export class DevicesComponent {

}
