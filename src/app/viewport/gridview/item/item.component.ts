import { Component, Injector, OnInit } from '@angular/core';
import { NgClass } from '@angular/common';

import { Chip } from 'primeng/chip';
import { Image } from 'primeng/image';
import { Fluid } from 'primeng/fluid';
import { Tag } from 'primeng/tag';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { MessageService } from 'primeng/api';
import { PanelModule } from 'primeng/panel';
import { AvatarModule } from 'primeng/avatar';
import { MenuModule } from 'primeng/menu';
import { DialogService, DynamicDialog, DynamicDialogRef } from 'primeng/dynamicdialog';

import { CBaseItem } from '../../base-item/base-item.component';
import { ItemDirective } from './item.directive';

@Component({
  selector: 'app-item',
  standalone: true,
  imports: [
    NgClass,
    CardModule, ButtonModule, PanelModule, AvatarModule, MenuModule, Chip, Image, Fluid, Tag, DynamicDialog,
  ],
  templateUrl: './item.component.html',
  styleUrl: './item.component.scss',
  outputs: ['selectItem'],
  providers: [{ provide: ItemDirective }, DialogService, MessageService],
})
export class CItem extends CBaseItem implements OnInit {
  currentClasses: Record<string, boolean> = {}

  isSpecial: boolean = false

  constructor(injectors: Injector) {
    super(injectors);
  }

  override ngOnInit(): void {
    this.currentClasses = {};
  }
}
