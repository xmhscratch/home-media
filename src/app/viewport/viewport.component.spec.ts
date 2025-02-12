import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CViewport } from './viewport.component';

describe('CViewport', () => {
  let component: CViewport;
  let fixture: ComponentFixture<CViewport>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CViewport]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CViewport);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
