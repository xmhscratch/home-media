import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CGridview } from './gridview.component';

describe('GridviewComponent', () => {
  let component: CGridview;
  let fixture: ComponentFixture<CGridview>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CGridview],
    }).compileComponents();

    fixture = TestBed.createComponent(CGridview);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
