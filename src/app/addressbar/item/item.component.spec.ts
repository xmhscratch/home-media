import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CItem } from './item.component';

describe('CItem', () => {
  let component: CItem;
  let fixture: ComponentFixture<CItem>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CItem]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CItem);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
