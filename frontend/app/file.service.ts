import { Injectable, WritableSignal } from '@angular/core';
import { signal } from '@angular/core';
import { toObservable } from '@angular/core/rxjs-interop';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';

import { Observable, Subject, switchMap, map, tap } from 'rxjs';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';

import { LRUCache } from 'lru-cache';
import { keyBy as ldKeyBy, get as ldGet } from 'lodash-es';

import { EnvService } from '@/env.service';
import { IFileListItem, ISocketMessage } from '../types/storage';
import { StorageService } from './storage.service';

export type TFileList = { [fileKey: string]: IFileListItem };

@Injectable({
  providedIn: 'root',
})
export class FileService {
  sockets = new LRUCache<string, WebSocketSubject<ISocketMessage>, unknown>({
    max: 500,
    maxSize: 5000,
    allowStale: false,
    updateAgeOnGet: false,
    updateAgeOnHas: false,
    ttl: 1000 * 60 * 1, // 1min

    sizeCalculation: (_) => 1,
    dispose: (ws: WebSocketSubject<ISocketMessage>, fileKey: string) => {
      if (ws === undefined || (ws && ws.closed)) {
        return;
      }
      // console.log(fileKey)
      ws.next({ event: 8 });
      ws.complete();
      ws.unsubscribe();
    },
    fetchMethod: async (
      fileKey,
      staleValue,
      // { options, signal, context }
    ) => {
      staleValue?.complete();
    },
  });

  files: WritableSignal<TFileList> = signal<TFileList>({});
  files$ = toObservable(this.files);

  changed$: Subject<IFileListItem> = new Subject<IFileListItem>();

  constructor(
    private env: EnvService,
    private storageService: StorageService,
    private http: HttpClient,
  ) {
    // this.changed$.subscribe({
    //   next: (x) => console.log(x),
    //   error: (err) => console.log(err),
    //   complete: () => console.log(237846),
    // })
  }

  getNativeSocket(f: IFileListItem): WebSocketSubject<ISocketMessage> {
    let ws: WebSocketSubject<ISocketMessage> | undefined;
    if (!this.sockets.has(f.fileKey)) {
      ws = this.newSocket(f);
    }
    if (!ws) {
      ws = this.sockets.get(f.fileKey);
    }
    if (ws === undefined || (ws && ws.closed)) {
      ws = this.newSocket(f);
    }
    // ws.error
    this.sockets.set(f.fileKey, ws);
    return ws;
  }

  getSocket(f: IFileListItem): Observable<ISocketMessage> {
    return this.getNativeSocket(f).asObservable();
  }

  disconnect(fileKey: string): boolean {
    const ws = this.sockets.get(fileKey);
    if (!ws) {
      return false;
    }
    this.sockets.delete(fileKey);
    return this.sockets.purgeStale();
  }

  loadFiles(): Observable<TFileList> {
    return this.storageService.loadSession().pipe(
      switchMap(({ root, active, paths, nodes }) => {
        return this.storageService.createSession(active);
      }),
      map((files) => {
        const list = ldKeyBy(files, 'fileKey');

        this.files.set(list);
        return list;
      }),
      // tap((v) => { console.log(v) }),
    );
  }

  activateFile(f: IFileListItem): Observable<IFileListItem> {
    return this.getSocket(f).pipe(
      // tap((v) => console.log(v)),
      map((v) => {
        const _parse = (stop1: boolean, stop2: boolean): ISocketMessage => {
          switch (typeof v) {
            case 'string':
              v = JSON.parse(v);
              break;
            default:
              stop1 = true;
              break;
          }
          switch (typeof v.payload) {
            case 'string':
              v.payload = JSON.parse(v.payload);
              break;
            default: {
              if (v.payload?.stage == 5) {
                switch (typeof v.payload?.message) {
                  case 'string':
                    v.payload.message = JSON.parse(v.payload?.message);
                    break;
                  default:
                    stop2 = true;
                    break;
                }
              } else {
                stop2 = true;
              }
              break;
            }
          }
          return stop1 && stop2 ? v : _parse(stop1, stop2);
        };
        return _parse(false, false);
      }),
      // tap((v) => console.log(v)),
      map((v: ISocketMessage) => {
        const stage: number = ldGet(v.payload, 'stage', 0);
        let result: IFileListItem = {
          ...f,
          isCompleted: !(!stage || stage < 5),
          isProgressing: stage > 0 && stage < 5,
          isDownloading: stage == 2,
        };
        if (stage == 5) {
          result = { ...result, ...(<{ message: object }>v.payload).message };
          result.sourceReady = !/(false|no|0|reject)/gi.test(
            <string>ldGet(result, 'sourceReady', '0'),
          );
        } else {
          result = { ...result, ...(<{}>v.payload) };
        }
        return result;
      }),
      // tap((v) => console.log(v)),
      tap((v) => this.changed$.next(v)),
    );
  }

  newSocket(f: IFileListItem) {
    return webSocket<ISocketMessage>(
      `ws://${this.env.get('endpoint.file')}/ws/${f.sessionId}/${f.fileKey}`,
    );
  }

  fetchFile(srcUrl: string, byteStart?: number, byteEnd?: number) {
    const headers = new HttpHeaders({
      range: `bytes=${byteStart || ''}-${byteEnd || ''}`,
      responseType: 'arraybuffer',
      endpoint: <string>this.env.get('endpoint.file'),
    });

    return this.http
      .get<unknown>(`${srcUrl}`, {
        headers,
        // observe: "body",
      })
      .pipe(tap((v) => console.log(v)));
  }
}
