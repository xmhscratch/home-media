import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CAddressbar } from './addressbar.component';

describe('CAddressbar', () => {
  let component: CAddressbar;
  let fixture: ComponentFixture<CAddressbar>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CAddressbar],
    }).compileComponents();

    fixture = TestBed.createComponent(CAddressbar);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
