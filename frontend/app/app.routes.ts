import { Routes } from '@angular/router';

import { AppComponent } from '@/app.component';
import { CViewport } from '@/viewport/viewport.component';

export const routes: Routes = [
  { path: '', component: AppComponent },
  {
    path: 'storage',
    redirectTo: 'storage/678bb472f4420b064e4ab471/678bb472f4420b064e4ab471',
    // redirectTo: 'storage/67a61a61f2f1b1cbb8164e08/67a61a61f2f1b1cbb8164e08',
    pathMatch: 'full',
  },
  {
    path: 'storage',
    children: [
      {
        path: ':rootId/:nodeId',
        component: CViewport,
      },
    ],
    outlet: 'primary',
  },
];
