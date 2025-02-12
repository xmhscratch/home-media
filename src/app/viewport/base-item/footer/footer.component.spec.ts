import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CFooter } from './footer.component';

describe('CFooter', () => {
  let component: CFooter;
  let fixture: ComponentFixture<CFooter>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CFooter]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CFooter);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
