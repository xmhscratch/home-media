import { Component } from '@angular/core';
import { OnInit } from '@angular/core';

import { DynamicDialogRef } from 'primeng/dynamicdialog';
import { ButtonModule } from 'primeng/button';
import { MenuItem } from 'primeng/api';
// import { Menubar } from 'primeng/menubar';
import { BadgeModule } from 'primeng/badge';
import { AvatarModule } from 'primeng/avatar';
import { InputTextModule } from 'primeng/inputtext';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    // Menubar,
    ButtonModule,
    BadgeModule,
    AvatarModule,
    InputTextModule,
    CommonModule,
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss',
})
export class CHeader implements OnInit {
  items: MenuItem[] | undefined;

  constructor(public ref: DynamicDialogRef) {}

  ngOnInit() {
    this.items = [
      {
        label: 'Home',
        icon: 'pi pi-home',
      },
    ];
  }

  closeDialog(data: any) {
    this.ref.close(data);
  }
}
