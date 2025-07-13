import { Component, EventEmitter, Input, Output, ViewChild, ElementRef } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';

export interface NewDevice {
  name: string;
  type: string;
  region: string;
}

@Component({
  selector: 'app-add-device-modal',
  imports: [
    FormsModule,
    CommonModule,
    TranslateModule
  ],
  templateUrl: './add-device-modal.component.html',
  styleUrl: './add-device-modal.component.scss',
})
export class AddDeviceModalComponent {
  @ViewChild('modal', { static: false }) modalRef!: ElementRef;
  @Output() addDevice = new EventEmitter<NewDevice>();

  isVisible = false;
  newDevice: NewDevice = {
    name: '',
    type: '',
    region: ''
  };

  // Programmatic methods that parent can call via ViewChild
  open() {
    this.isVisible = true;
    this.resetForm();
  }

  close() {
    this.isVisible = false;
    this.resetForm();
  }

  onSubmit() {
    if (this.isFormValid()) {
      this.addDevice.emit({ ...this.newDevice });
      this.close();
    }
  }

  private resetForm() {
    this.newDevice = {
      name: '',
      type: '',
      region: ''
    };
  }

  private isFormValid(): boolean {
    return !!(this.newDevice.name && this.newDevice.type && this.newDevice.region);
  }

  onOverlayClick() {
    this.close();
  }

  onModalClick(event: Event) {
    event.stopPropagation();
  }
}