import { Injectable, Inject, InjectionToken } from '@angular/core';
import { isDevMode } from '@angular/core';
import { get as ldGet } from 'lodash-es';

import { IEnvConfig } from '../types/env.d';
import { environment } from '../environment';

const envName = isDevMode() ? 'development' : 'production';

const GLOBAL_ENV_CONFIG = new InjectionToken<IEnvConfig>('GLOBAL_ENV_CONFIG', {
  providedIn: 'root',
  factory: () => {
    return <IEnvConfig>{
      env: envName,
      endpoint: {
        backend: environment.endpoint.backend,
        api: environment.endpoint.backend,
        file: environment.endpoint.file,
        frontend: environment.endpoint.frontend,
      },
    };
  },
});

@Injectable({
  providedIn: 'root',
})
export class EnvService {
  constructor(@Inject(GLOBAL_ENV_CONFIG) private env: IEnvConfig) {}

  get(path: string, defaultValue?: unknown): unknown {
    return ldGet(this.env || {}, path, defaultValue);
  }
}
