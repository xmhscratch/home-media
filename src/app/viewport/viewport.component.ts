import { Component, WritableSignal, signal, inject } from '@angular/core';
import { OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';

// import { tap, switchMap } from 'rxjs/operators';
// import { ElementRef, afterRender } from '@angular/core';
// import { $dt } from '@primeng/themes';

import { ButtonModule } from 'primeng/button';

import { zip, combineLatest, switchMap, Subscription } from 'rxjs';
import { map, isEmpty } from 'rxjs/operators'
import { toObservable } from '@angular/core/rxjs-interop';

import { IINode } from '../../types/storage'
import { StorageService } from '@/storage.service'

import { CGridview } from './gridview/gridview.component';

@Component({
  selector: 'app-viewport',
  standalone: true,
  imports: [ButtonModule, CGridview],
  templateUrl: './viewport.component.html',
  styleUrl: './viewport.component.scss',
})

export class CViewport implements OnInit, OnDestroy {

  private readonly route = inject(ActivatedRoute);

  rootId: WritableSignal<string> = signal('');
  nodeId: WritableSignal<string> = signal('');

  rootId$ = toObservable(this.rootId);
  nodeId$ = toObservable(this.nodeId);

  root: WritableSignal<IINode | null> = signal(null);
  paths: WritableSignal<IINode[]> = signal([]);
  nodes: WritableSignal<IINode[]> = signal([]);

  root$ = toObservable(this.root);
  paths$ = toObservable(this.paths);
  nodes$ = toObservable(this.nodes);

  destroy$: Subscription = new Subscription();

  constructor(
    private storage: StorageService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.destroy$.add(
      this.route.paramMap
        .pipe(
          map((params) => {
            const rootId = (params.get('rootId') || '');
            const nodeId = (params.get('nodeId') || '');
            return { rootId, nodeId }
          }),
        )
        .subscribe(({ rootId, nodeId }) => {
          this.rootId.set(rootId);
          this.nodeId.set(nodeId);
        })
    )
  }

  ngOnDestroy() {
    if (this.destroy$) { this.destroy$.unsubscribe() }
  }
}
