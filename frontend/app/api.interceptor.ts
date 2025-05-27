import { inject } from '@angular/core';
import { DOCUMENT } from '@angular/common';
import {
  HttpInterceptorFn,
  HttpRequest,
  HttpEvent,
  HttpHandlerFn,
} from '@angular/common/http';

import { Observable } from 'rxjs';
// import { Subject, Subscription } from 'rxjs';

import { EnvService } from './env.service';

export const apiInterceptor: HttpInterceptorFn = (
  req: HttpRequest<unknown>,
  next: HttpHandlerFn,
): Observable<HttpEvent<unknown>> => {
  const env = inject(EnvService);
  const document: Document = inject(DOCUMENT);

  const endpoint: string =
    <string>req.headers.get('endpoint') || <string>env.get('endpoint.backend');

  const apiReq = req.clone({
    url: new URL(
      req.url,
      `${document.location.protocol}//` + endpoint,
    ).toString(),
  });

  return next(apiReq);

  // const completed: Subject<any> = new Subject<any>();

  // const obs = new Observable<HttpEvent<unknown>>((o) => {
  //   next(apiReq)
  //     .pipe(
  //       // tap(console.log),
  //       takeUntil(completed),
  //     )
  //     .subscribe((v: any) => {
  //       o.next(v);
  //       o.complete();
  //     });
  // })
  // obs.subscribe(() => {
  //   completed.next(void 0);
  // })
  // return obs
};
