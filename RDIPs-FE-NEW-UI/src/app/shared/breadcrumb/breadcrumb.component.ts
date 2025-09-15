import { Component, Output, EventEmitter, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-breadcrumb',
  imports: [TranslateModule],
  templateUrl: './breadcrumb.component.html',
  styleUrl: './breadcrumb.component.scss'
})
export class BreadcrumbComponent {
  @Input() title = 'HOME PAGE';

  @Output() customize: EventEmitter<any> = new EventEmitter();
  @Output() export: EventEmitter<any> = new EventEmitter();

  customizeClicked($event: any): void{
    this.customize.emit($event);
  }

  exportClicked($event: any): void{
    this.export.emit($event);
  }
}
