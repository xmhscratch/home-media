import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CPlayer } from './player.component';

describe('CPlayer', () => {
  let component: CPlayer;
  let fixture: ComponentFixture<CPlayer>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CPlayer]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CPlayer);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
