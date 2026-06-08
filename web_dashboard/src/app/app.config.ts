import { APP_INITIALIZER, ApplicationConfig, inject, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideRouter } from '@angular/router';
import { HttpClient, provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { catchError, of } from 'rxjs';
import { routes } from './app.routes';
import { authInterceptor } from './core/interceptors/auth.interceptor';
import { errorInterceptor } from './core/interceptors/error.interceptor';

// Ping /api/health on startup so the maintenance overlay shows immediately
// if the backend is already down when the page loads.
// The errorInterceptor handles setting MaintenanceService.offline — we just
// need to fire the request and swallow the re-thrown error so the app still boots.
function healthCheckInitializer() {
  const http = inject(HttpClient);
  return () => http.get('/api/health').pipe(catchError(() => of(null))).toPromise();
}

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes),
    provideHttpClient(withInterceptors([authInterceptor, errorInterceptor])),
    provideAnimationsAsync(),
    { provide: APP_INITIALIZER, useFactory: healthCheckInitializer, multi: true },
  ],
};
