import { Component, ViewChild, OnInit } from '@angular/core';
// import { NgTemplateOutlet } from '@angular/common';
import { RouterOutlet } from '@angular/router';

import { SpeedDial } from 'primeng/speeddial';
import { CardModule } from 'primeng/card';

import { IINode } from '../types/storage'
import { StorageService } from '@/storage.service'

import { CNavigation } from './navigation/navigation.component';
import { CViewport } from './viewport/viewport.component';
import { CToolbar } from './toolbar/toolbar.component'

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    SpeedDial, CardModule,
    CToolbar, CNavigation, CViewport,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent implements OnInit {
  @ViewChild('navigation') navElRef!: CNavigation;

  title = 'home-media';

  constructor(
    private storage: StorageService,
  ) { }

  ngOnInit(): void { }

  ngAfterViewInit() { }

  toggleNavigationPanel() {
    this.navElRef.visible = true
  }
}
