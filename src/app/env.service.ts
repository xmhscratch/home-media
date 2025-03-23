import { Injectable, Inject, InjectionToken } from '@angular/core';
import { inject, isDevMode, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { get as ldGet, map as ldMap } from 'lodash-es';

import {
  IEnvConfig,
} from '../types/env.d';

const envName = isDevMode() ? 'development' : 'production';

const GLOBAL_ENV_CONFIG = new InjectionToken<IEnvConfig>('GLOBAL_ENV_CONFIG', {
  providedIn: 'root',
  factory: (
  ) => {
    const platformId: Object = inject(PLATFORM_ID);
    const envConfig = require(`../../config.${envName}${isPlatformBrowser(platformId) ? '.browser' : ''}.json`);

    return <IEnvConfig>{
      env: envName,
      endpoint: {
        backend: envConfig.endpoint.backend,
        api: envConfig.endpoint.backend,
        file: envConfig.endpoint.file,
        frontend: envConfig.endpoint.frontend,
      },
    };
  }
});

@Injectable({
  providedIn: 'root'
})
export class EnvService {

  constructor(
    @Inject(GLOBAL_ENV_CONFIG) private env: IEnvConfig,
  ) { }

  get(path: string, defaultValue?: any): any {
    return ldGet(this.env || {}, path, defaultValue);
  }
}
