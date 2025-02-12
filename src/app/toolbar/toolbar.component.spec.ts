import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CToolbar } from './toolbar.component';

describe('CToolbar', () => {
  let component: CToolbar;
  let fixture: ComponentFixture<CToolbar>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CToolbar]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CToolbar);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
