import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideClientHydration } from '@angular/platform-browser';
import { withEventReplay, withHttpTransferCacheOptions } from '@angular/platform-browser';
import { HttpRequest, provideHttpClient, withFetch, withInterceptors } from '@angular/common/http';

import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { providePrimeNG } from 'primeng/config';

import { apiInterceptor } from '@/api.interceptor';
import DefaultThemePreset from '@/theme';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideClientHydration(
      withEventReplay(),
      withHttpTransferCacheOptions({
          filter: (req: HttpRequest<unknown>) => true, // to filter
          includeHeaders: [], // to include headers
          includePostRequests: true, // to include POST
          includeRequestsWithAuthHeaders: false, // to include with auth
      }),
    ),
    provideHttpClient(
      withFetch(),
      withInterceptors([apiInterceptor]),
    ),
    provideAnimationsAsync(),
    providePrimeNG({
      theme: {
        preset: DefaultThemePreset,
        options: {
          darkModeSelector: '.app-dark',
        }
      }
    })
  ]
};
