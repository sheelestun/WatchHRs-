import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, tap, throwError } from 'rxjs';
import { MaintenanceService } from '../services/maintenance.service';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  if (!req.url.includes('/api/')) {
    return next(req);
  }

  const maintenance = inject(MaintenanceService);

  return next(req).pipe(
    tap(() => maintenance.setOffline(false)),
    catchError((err) => {
      if (err.status === 0 || err.status >= 500) {
        maintenance.setOffline(true);
      }
      return throwError(() => err);
    }),
  );
};
