import { Component, WritableSignal, Signal } from '@angular/core';
import { signal, computed, inject } from '@angular/core';
import { OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute, RouterModule } from '@angular/router';
import { NgIf, NgClass } from '@angular/common';

import { MenuItem } from 'primeng/api';
import { Breadcrumb } from 'primeng/breadcrumb';
import { MessageModule } from 'primeng/message';
import { DividerModule } from 'primeng/divider';
import { Tooltip } from 'primeng/tooltip';
import { CardModule } from 'primeng/card';

import { tap, map } from 'rxjs/operators';

import { CItem } from './item/item.component';
import { IINode } from '../../types/storage';
import { StorageService } from '@/storage.service';

@Component({
  selector: 'app-addressbar',
  standalone: true,
  imports: [
    NgIf,
    // NgClass,
    RouterModule,
    Breadcrumb,
    MessageModule,
    DividerModule,
    Tooltip,
    CardModule,
  ],
  templateUrl: './addressbar.component.html',
  styleUrl: './addressbar.component.scss',
})
export class CAddressbar implements OnInit, OnDestroy {
  private readonly route = inject(ActivatedRoute);

  breadcrumbStyles = {
    padding: '0rem',
    // background: 'transparent',
  };

  cardStyles = {
    bodyPadding: '0.25rem 1.25rem',
  };

  paths: WritableSignal<CItem[]> = signal<CItem[]>([]);
  items: Signal<CItem[]> = computed(() => <CItem[]>this.paths() || []);

  home: MenuItem | undefined;

  constructor(private storage: StorageService) {}

  ngOnInit(): void {
    this.route.paramMap
      .pipe(
        tap(() => {
          // this.rootId$.subscribe()
          // this.nodeId$.subscribe()
        }),
      )
      .subscribe((a) => {
        // const { paths } = this.storage.paths();
        // console.log(a)
      });

    this.home = { icon: 'pi pi-home', routerLink: '/storage/' };
    // this.paths.set(<CItem[]>paths);
  }

  ngOnDestroy(): void {}
}
