import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CBaseItem } from './base-item.component';

describe('CBaseItem', () => {
  let component: CBaseItem;
  let fixture: ComponentFixture<CBaseItem>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CBaseItem]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CBaseItem);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
