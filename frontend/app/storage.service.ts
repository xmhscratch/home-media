import { Injectable, signal, WritableSignal } from '@angular/core';
import { toObservable } from '@angular/core/rxjs-interop';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';

import {
  Observable,
  Subject,
  zip,
  combineLatest,
  of,
  Subscription,
} from 'rxjs';
import { switchMap, map, tap } from 'rxjs/operators';

import { map as ldMap } from 'lodash-es';
// import { DynamicDialogRef } from 'primeng/dynamicdialog';
// import { forEach } from 'lodash-es';

import {
  IINode,
  ITreeRootNode,
  FileMetaInfo,
  SessionInfo,
  IFileListItem,
} from '../types/storage';

import { EnvService } from './env.service';

@Injectable({
  providedIn: 'root',
})
export class StorageService {
  roots$: Observable<ITreeRootNode[]> = of([]).pipe(
    switchMap(() => this.fetchRoots()),
  );

  rootId: WritableSignal<string> = signal<string>('');
  // rootId: WritableSignal<string> = signal<string>('678bb472f4420b064e4ab471')
  // rootId: WritableSignal<string> = signal<string>('67a61a61f2f1b1cbb8164e08')
  rootId$ = toObservable(this.rootId);

  activeId: WritableSignal<string> = signal<string>(this.rootId());
  activeId$ = toObservable(this.activeId);

  root$: Observable<IINode> = this.rootId$.pipe(
    switchMap((rootId) => this.fetchOne(rootId, rootId, 'node')),
  );
  active$: Observable<IINode> = combineLatest([
    this.rootId$,
    this.activeId$,
  ]).pipe(
    switchMap(([rootId, activeId]) => this.fetchOne(rootId, activeId, 'node')),
  );
  paths$: Observable<IINode[]> = combineLatest([
    this.rootId$,
    this.activeId$,
  ]).pipe(
    switchMap(([rootId, activeId]) => this.fetchAll(rootId, activeId, 'paths')),
  );
  nodes$: Observable<IINode[]> = combineLatest([
    this.rootId$,
    this.activeId$,
  ]).pipe(
    switchMap(([rootId, activeId]) =>
      this.fetchAll(rootId, activeId, 'children'),
    ),
  );

  defaultHeaders: HttpHeaders = new HttpHeaders({
    enctype: 'multipart/form-data',
    endpoint: <string>this.env.get('endpoint.backend'),
  });

  onchange = new Subject<unknown>();
  switchNode$: Subscription = new Subscription();

  constructor(
    private env: EnvService,
    private http: HttpClient,
  ) {}

  getRoots() {
    return this.roots$
      .pipe
      // tap((v) => { console.log(v) }),
      ();
  }

  switchNode(rootId?: string, nodeId?: string) {
    this.switchNode$.unsubscribe();

    if (rootId) {
      if (this.rootId() == '') {
        this.rootId.set(rootId);
      } else {
        this.rootId.update((old: string) => (!rootId ? old : rootId));
      }
      if (this.activeId() == '') {
        this.activeId.set(nodeId || rootId);
      } else {
        this.activeId.update((old: string) => (!nodeId ? old : nodeId));
      }
    }

    return zip(this.root$, this.active$, this.paths$, this.nodes$).pipe(
      // tap((v) => { console.log(v) }),
      map(([root, active, paths, nodes]) => ({ root, active, paths, nodes })),
      tap((v) => {
        this.onchange.next(v);
      }),
      tap(() => {
        this.switchNode$.add(this.root$.subscribe());
        this.switchNode$.add(this.active$.subscribe());
        this.switchNode$.add(this.paths$.subscribe());
        this.switchNode$.add(this.nodes$.subscribe());
      }),
      // tap(({ root, active, paths, nodes }) => {
      //   console.log({ root, active, paths, nodes })
      //   // switchMap(([rootId, activeId]) => this.fetchAll(rootId, activeId, 'children')),
      // }),
      // tap((v) => { console.log(v) }),
      // share(),
    );
  }

  loadSession() {
    return zip(this.root$, this.active$, this.paths$, this.nodes$).pipe(
      // tap((v) => { console.log(v) }),
      map(([root, active, paths, nodes]) => ({ root, active, paths, nodes })),
    );
  }

  fetchAll(
    rootId: string,
    nodeId: string,
    fnName: string,
  ): Observable<IINode[]> {
    return this.http
      .get<IINode[]>(`/storage/${rootId}/${nodeId || rootId}/${fnName}`, {
        headers: this.defaultHeaders,
        params: new HttpParams({ fromString: `extend=1&branch=0` }),
        // observe: "body",
      })
      .pipe(map((data: IINode[]) => data));
  }

  fetchOne(rootId: string, nodeId: string, fnName: string): Observable<IINode> {
    return this.http
      .get<IINode>(`/storage/${rootId}/${nodeId || rootId}/${fnName}`, {
        headers: this.defaultHeaders,
        params: new HttpParams({ fromString: `extend=1&branch=0` }),
        // observe: "body",
      })
      .pipe(map((data: IINode) => data));
  }

  fetchRoots(): Observable<ITreeRootNode[]> {
    return this.http
      .get<ITreeRootNode[]>(`/storage/roots`, {
        headers: this.defaultHeaders,
        // params: new HttpParams({ fromString: `` }),
        // observe: "body",
      })
      .pipe(map((data: ITreeRootNode[]) => data));
  }

  createSession(node: IINode): Observable<IFileListItem[]> {
    const formData = new FormData();

    zip(Object.keys(node), Object.values(node))
      .pipe(
        // tap((v) => {
        //   console.log(v);
        // }),
        map(([key, val]) => {
          formData.append(key, <string | Blob>val);
        }),
        // tap((v) => {
        //   console.log(v);
        // }),
      )
      .subscribe()
      .unsubscribe();

    const headers = new HttpHeaders({
      enctype: 'multipart/form-data',
      endpoint: <string>this.env.get('endpoint.api'),
    });

    return this.http
      .put<SessionInfo>(`/`, formData, {
        headers,
        // observe: "body",
      })
      .pipe(
        // tap((v) => console.log(v)),
        map((ss: SessionInfo): IFileListItem[] => {
          return ldMap(
            ss.files,
            (v: FileMetaInfo, k: string) =>
              <IFileListItem>{
                path: v.path,
                size: v.size,
                fileKey: k,
                sessionId: ss.id,
                nodeId: ss.nodeId,
              },
          );
        }),
        // tap((v) => console.log(v)),
      );
  }

  fetchSource(sessionId: string, filePath: string): Observable<IINode> {
    const formData = new FormData();

    const headers = new HttpHeaders({
      enctype: 'multipart/form-data',
      endpoint: <string>this.env.get('endpoint.api'),
    });

    return this.http
      .post<IINode>(`/${sessionId}/${filePath}`, formData, {
        headers,
        // observe: "body",
      })
      .pipe(map((data: IINode) => data));
  }
}
