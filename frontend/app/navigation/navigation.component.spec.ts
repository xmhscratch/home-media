import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CNavigation } from './navigation.component';

describe('CNavigation', () => {
  let component: CNavigation;
  let fixture: ComponentFixture<CNavigation>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CNavigation],
    }).compileComponents();

    fixture = TestBed.createComponent(CNavigation);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
