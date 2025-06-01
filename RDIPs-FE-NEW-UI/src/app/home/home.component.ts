import { NgFor } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { Dashboard_Data } from '../constant';
import { BreadcrumbComponent } from '../shared/breadcrumb/breadcrumb.component';
import { NavigationComponent } from '../shared/navigation/navigation.component';
import { DataCardComponent } from './data-card/data-card.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    RouterModule,
    NavigationComponent,
    DataCardComponent,
    NgFor,
    BreadcrumbComponent,
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss',
})
export class HomeComponent {
  CARD_DATA = Dashboard_Data;

  trackByFn(index: number, item: any): any {
    return item.id; // Assuming each item has a unique 'id' property
  }
}
