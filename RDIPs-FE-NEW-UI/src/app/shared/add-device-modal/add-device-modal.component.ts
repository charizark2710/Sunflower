import { Component, EventEmitter, Input, Output, OnInit, OnChanges } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TranslateModule } from '@ngx-translate/core';

export interface NewDevice {
  name: string;
  type: string;
  region: string;
  status?: string;
}

export interface EditDevice extends NewDevice {
  device_id: string;
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
export class AddDeviceModalComponent implements OnInit, OnChanges {
  @Input() isVisible = false;
  @Input() isEditMode = false;
  @Input() deviceToEdit: EditDevice | null = null;
  @Output() close = new EventEmitter<void>();
  @Output() addDevice = new EventEmitter<NewDevice>();
  @Output() updateDevice = new EventEmitter<EditDevice>();

  newDevice: NewDevice = {
    name: '',
    type: '',
    region: '',
    status: ''
  };

  ngOnInit() {
    // Initialize form data when deviceToEdit changes
    if (this.isEditMode && this.deviceToEdit) {
      this.newDevice = {
        name: this.deviceToEdit.name,
        type: this.deviceToEdit.type,
        region: this.deviceToEdit.region,
        status: this.deviceToEdit.status || 'Active'
      };
    }
  }

  ngOnChanges() {
    // Update form when deviceToEdit input changes
    if (this.isEditMode && this.deviceToEdit) {
      this.newDevice = {
        name: this.deviceToEdit.name,
        type: this.deviceToEdit.type,
        region: this.deviceToEdit.region,
        status: this.deviceToEdit.status || 'Active'
      };
    } else if (!this.isEditMode) {
      this.resetForm();
    }
  }

  closeModal() {
    this.resetForm();
    this.close.emit();
  }

  onSubmit() {
    if (this.isFormValid()) {
      if (this.isEditMode && this.deviceToEdit) {
        // Emit update event with device ID
        const updatedDevice: EditDevice = {
          device_id: this.deviceToEdit.device_id,
          name: this.newDevice.name,
          type: this.newDevice.type,
          region: this.newDevice.region,
          status: this.newDevice.status || 'Active'
        };
        this.updateDevice.emit(updatedDevice);
      } else {
        // Emit add event
        this.addDevice.emit({ ...this.newDevice });
      }
      this.resetForm();
    }
  }

  private resetForm() {
    this.newDevice = {
      name: '',
      type: '',
      region: '',
      status: 'Active'
    };
  }

  private isFormValid(): boolean {
    if (this.isEditMode) {
      // For edit mode, only name and status are required (region/type are disabled)
      return !!(this.newDevice.name && this.newDevice.status);
    }
    // For add mode, all fields are required
    return !!(this.newDevice.name && this.newDevice.type && this.newDevice.region);
  }

  onOverlayClick() {
    this.closeModal();
  }

  onModalClick(event: Event) {
    event.stopPropagation();
  }
}
