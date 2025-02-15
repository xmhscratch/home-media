import { HttpInterceptorFn, HttpRequest, HttpEvent, HttpHandlerFn } from '@angular/common/http';

import { Subject, Observable, Subscription } from 'rxjs';

import { ENDPOINT_STORAGE_API } from './environment';

export const apiInterceptor: HttpInterceptorFn = (req: HttpRequest<any>, next: HttpHandlerFn): Observable<HttpEvent<unknown>> => {
  const endpoint = req.headers.get('endpoint') || ENDPOINT_STORAGE_API
  const apiReq = req.clone({ url: new URL(req.url, endpoint).toString() });

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
