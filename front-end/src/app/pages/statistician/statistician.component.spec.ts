import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StatisticianComponent } from './statistician.component';

describe('StatisticianComponent', () => {
  let component: StatisticianComponent;
  let fixture: ComponentFixture<StatisticianComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StatisticianComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(StatisticianComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
