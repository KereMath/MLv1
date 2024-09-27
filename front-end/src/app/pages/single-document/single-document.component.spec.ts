import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SingleDocumentComponent } from './single-document.component';

describe('SingleDocumentComponent', () => {
  let component: SingleDocumentComponent;
  let fixture: ComponentFixture<SingleDocumentComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SingleDocumentComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SingleDocumentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
