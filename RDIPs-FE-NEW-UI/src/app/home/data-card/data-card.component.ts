import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-data-card',
  standalone: true,
  imports: [TranslateModule],
  templateUrl: './data-card.component.html',
  styleUrl: './data-card.component.scss'
})
export class DataCardComponent {
  @Input() card!: {
    title: string;
    number: string;
    rate: string;
  };
}
