import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DeviceActivityStatisticComponent } from './device-activity-statistic.component';

describe('DeviceActivityStatisticComponent', () => {
  let component: DeviceActivityStatisticComponent;
  let fixture: ComponentFixture<DeviceActivityStatisticComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DeviceActivityStatisticComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DeviceActivityStatisticComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
