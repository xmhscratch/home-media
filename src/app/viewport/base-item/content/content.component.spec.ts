import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CContent } from './content.component';

describe('CContent', () => {
  let component: CContent;
  let fixture: ComponentFixture<CContent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CContent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CContent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
