import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CHeader } from './header.component';

describe('CHeader', () => {
  let component: CHeader;
  let fixture: ComponentFixture<CHeader>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CHeader],
    }).compileComponents();

    fixture = TestBed.createComponent(CHeader);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
