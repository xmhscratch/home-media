import {
  Component,
  Input,
  WritableSignal,
  signal,
  inject,
} from '@angular/core';
import { OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

import { ButtonModule } from 'primeng/button';
import { PanelModule } from 'primeng/panel';
import { AvatarModule } from 'primeng/avatar';
import { MenuModule } from 'primeng/menu';
import { ProgressSpinner } from 'primeng/progressspinner';
import { Skeleton } from 'primeng/skeleton';

// import { combineLatest, Observable, Subject, zip } from 'rxjs';
import { tap, isEmpty, map } from 'rxjs/operators';
import { toObservable } from '@angular/core/rxjs-interop';

// import { ElementRef, afterRender } from '@angular/core';
// import { $dt } from '@primeng/themes';

import { orderBy as ldOrderBy } from 'lodash-es';
import { IINode } from '../../../types/storage';
import { StorageService } from '@/storage.service';
import { CItem } from './item/item.component';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-gridview',
  standalone: true,
  imports: [
    ButtonModule,
    PanelModule,
    AvatarModule,
    MenuModule,
    ProgressSpinner,
    Skeleton,
    CItem,
  ],
  templateUrl: './gridview.component.html',
  styleUrl: './gridview.component.scss',
})
export class CGridview implements OnInit, OnDestroy {
  private readonly route = inject(ActivatedRoute);

  @Input() rootId: WritableSignal<string> = signal<string>('');
  @Input() nodeId: WritableSignal<string> = signal<string>('');

  rootId$ = toObservable(this.rootId);
  nodeId$ = toObservable(this.nodeId);

  paths: WritableSignal<IINode[]> = signal<IINode[]>([]);
  nodes: WritableSignal<IINode[]> = signal<IINode[]>([]);

  loaded: WritableSignal<boolean> = signal<boolean>(false);

  destroy$: Subscription = new Subscription();

  constructor(
    private storage: StorageService,
    private router: Router,
  ) {}

  ngOnInit() {
    this.destroy$.add(
      this.route.paramMap
        .pipe(
          tap(() => {
            this.destroy$.add(this.rootId$.subscribe());
            this.destroy$.add(this.nodeId$.subscribe());
          }),
        )
        .subscribe((a: ParamMap) => {
          this.destroy$.add(
            this.storage
              .switchNode(this.rootId(), this.nodeId())
              .pipe(
                tap(() => this.loaded.set(true)),
                map(({ nodes }) => ({
                  nodes: ldOrderBy(
                    nodes,
                    (i: IINode) => new Date(i.created_date as string).valueOf(),
                    'desc',
                  ),
                })),
                // tap((v) => console.log(v)),
              )
              .subscribe(({ nodes }) => {
                this.nodes.set(nodes);
              }),
          );
        }),
    );
  }

  ngOnDestroy(): void {
    if (this.destroy$) {
      this.destroy$.unsubscribe();
    }
  }

  selectItemHandler = ($event: Event, item: IINode) => {
    // this.router.navigate(['storage', `${item.root}`, `${item.id}`]);
  };
}

// constructor(elementRef: ElementRef) {
//   afterRender({
//     write: () => {},
//     read: () => {
//       console.log($dt('chip.icon'))
//     }
//   });
// }
