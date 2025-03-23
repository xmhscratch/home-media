import { Component, ViewChild, WritableSignal } from '@angular/core'
import { signal, inject } from '@angular/core'
import { ChangeDetectorRef, OnInit } from '@angular/core'
import { toObservable } from '@angular/core/rxjs-interop';
import { Router, ActivatedRoute } from '@angular/router';

import { DrawerModule } from 'primeng/drawer'
import { ButtonModule } from 'primeng/button'
import { AvatarModule } from 'primeng/avatar'
// import { StyleClass } from 'primeng/styleclass'
import { Drawer } from 'primeng/drawer'
import { TreeNode } from 'primeng/api'
import { Tree } from 'primeng/tree'

import { tap, map, defer } from 'rxjs'
import { mergeMap } from 'rxjs/operators'

import { IINode, ITreeRootNode } from '../../types/storage'
import { StorageService } from '@/storage.service'

@Component({
  selector: 'app-navigation',
  standalone: true,
  imports: [DrawerModule, ButtonModule, AvatarModule, Tree],
  templateUrl: './navigation.component.html',
  styleUrl: './navigation.component.scss',
  providers: [StorageService],
})
export class CNavigation implements OnInit {
  @ViewChild('drawerRef') drawerRef!: Drawer

  private readonly route = inject(ActivatedRoute);

  visible: boolean = false

  roots: WritableSignal<ITreeRootNode[]> = signal<ITreeRootNode[]>([])
  roots$ = toObservable(this.roots)

  loaded: WritableSignal<boolean> = signal<boolean>(false)

  treeNodes: Array<TreeNode> = Array<TreeNode>()

  constructor(
    private storage: StorageService,
    private router: Router,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnInit() {
    // this.loaded.set(false)

    this.storage.getRoots()
      .pipe(
      // tap((v) => console.log(v)),
    )
      .subscribe((roots) => {
        this.roots.set(roots)

        this.roots$.pipe(
          mergeMap((v: Array<ITreeRootNode>) => v),
          map((v: ITreeRootNode, k) => {
            return <TreeNode>{
              key: `${k}`,
              label: v.label,
              data: v.rootId,
              icon: 'pi pi-fw pi-inbox',
            }
          }),
        )
          .subscribe((v) => this.treeNodes.push(v))

        this.loaded.set(true)
        this.changeDetectorRef.markForCheck()
      })
  }

  closeCallback(e: any): void {
    this.drawerRef.close(e)
  }

  handleNodeSelect(event: any) {
    if (!event.node) { return }
    const rootId = event.node.data

    this.router.navigate(['storage', rootId, rootId]);
  }
}
