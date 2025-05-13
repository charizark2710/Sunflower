import { Component } from '@angular/core';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-admins',
  imports: [RouterModule, NavigationComponent],
  templateUrl: './admins.component.html',
  styleUrl: './admins.component.scss'
})
export class AdminsComponent {

}
