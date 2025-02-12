import { Component, OnInit } from '@angular/core';

import { FormsModule } from '@angular/forms';
import { SelectButton, SelectButtonModule } from 'primeng/selectbutton';

import { Toolbar } from 'primeng/toolbar';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';

import { CAddressbar } from '../addressbar/addressbar.component';

// import { ElementRef, afterRender } from '@angular/core';
// import { $dt } from '@primeng/themes';

@Component({
  selector: 'app-toolbar',
  standalone: true,
  imports: [
    CAddressbar,
    FormsModule, SelectButton, SelectButtonModule, Toolbar, ButtonModule, InputTextModule,
  ],
  templateUrl: './toolbar.component.html',
  styleUrl: './toolbar.component.scss'
})
export class CToolbar implements OnInit {

  toolbarStyles = {
    borderColor: 'none',
  }

  viewModeOpts: any[] = [
    { label: 'Grid', value: 'grid', },
    { label: 'List', value: 'list' },
  ];

  themeModeOpts: any[] = [
    { label: 'Dark Mode', icon: 'pi pi-sun', value: 'dark' },
    { label: 'Light Mode', icon: 'pi pi-moon', value: 'light', },
  ];

  viewMode: string = 'grid';
  themeMode: string = 'light';

  ngOnInit() {
    // console.log(this.value)
  }
  // constructor(elementRef: ElementRef) {
  //   afterRender({
  //     write: () => {},
  //     read: () => {
  //       console.log($dt('toolbar.border.color'))
  //     }
  //   });
  // }
  toggleThemeMode() {
    this.setThemeMode(this.themeMode==='dark'?'light':'dark');
  }

  setThemeMode(themeMode?: string) {
    if (!themeMode) { return }
    const element = document.querySelector('html');
    element?.classList[themeMode==='dark'?'add':'remove']('app-dark');
  }
}
