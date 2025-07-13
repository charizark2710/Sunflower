import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-text-group',
  imports: [TranslateModule],
  standalone: true,
  templateUrl: './text-group.component.html',
  styleUrl: './text-group.component.scss'
})
export class TextGroupComponent {
  @Input() title: string = '';
  @Input() description: string = '';
}
