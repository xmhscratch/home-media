import { Component, Input, Injector } from '@angular/core';
import { OnInit, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';

import { MessageService } from 'primeng/api';
import { ToastModule } from 'primeng/toast';
import { DialogService, DynamicDialog, DynamicDialogRef } from 'primeng/dynamicdialog';

import { IINode } from '../../../types/storage'
import { StorageService } from '@/storage.service'
import { CContent } from './content/content.component';
import { CFooter } from './footer/footer.component';
import { CHeader } from './header/header.component';
import { Subscription } from 'rxjs';
import { ConsoleLogger } from '@nestjs/common';

@Component({
  selector: 'app-base-item',
  standalone: true,
  imports: [DynamicDialog, ToastModule],
  // templateUrl: './base-item.component.html',
  styleUrl: './base-item.component.scss',
  providers: [DialogService, MessageService],
  template: '',
})
export class CBaseItem implements OnInit, OnDestroy {
  @Input() item: IINode = <IINode>{};

  cardStyles = {
    shadow: '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1);'
  }

  chipStyles = {
    borderRadius: '50px',
    paddingX: `0.35rem`,
    paddingY: `0.35rem`,
    iconSize: `0.5rem`,
    iconFontSize: '0.5rem',
    iconColor: '{rose.500}',
  }

  ref: DynamicDialogRef | undefined;

  router: Router;
  dialogService: DialogService;
  messageService: MessageService;
  storage: StorageService;

  destroy$: Subscription = new Subscription();

  constructor(
    protected injectors: Injector,
  ) {
    this.router = this.injectors.get(Router);
    this.storage = this.injectors.get(StorageService);
    this.dialogService = this.injectors.get(DialogService);
    this.messageService = this.injectors.get(MessageService);
  }

  ngOnInit(): void { }

  ngOnDestroy(): void {
    if (this.destroy$) { this.destroy$.unsubscribe() }
    if (this.ref) { this.ref.close(); }
  }

  handleEventClick($event: Event) {
    const isFolder = (this.item.depth || 0) > 0;
    const rootId: string = this.item.root;
    const nodeId: string = this.item.id;

    if (isFolder) {
      this.router.navigate(['storage', rootId, nodeId]);
    }
    else {
      this.destroy$.add(
        this.storage
          .switchNode(rootId, nodeId)
          .subscribe(() => this.show())
      )
    }
  }

  show() {
    this.ref = this.dialogService.open(CContent, {
      width: '65vw',
      modal: true,
      contentStyle: { overflow: 'hidden' },
      breakpoints: {
        '960px': '75vw',
        '640px': '90vw'
      },
      // data: this.item,
      templates: {
        footer: CFooter,
        header: CHeader,
      },
    });

    this.destroy$.add(
      this.ref.onClose.subscribe((data: any) => {
        this.ref = undefined
      })
    )
    this.destroy$.add(
      this.ref.onMaximize.subscribe((value) => { })
    )
  }
}
