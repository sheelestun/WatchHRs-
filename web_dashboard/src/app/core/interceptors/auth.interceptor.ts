import { HttpInterceptorFn, HttpRequest, HttpHandlerFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, switchMap, throwError } from 'rxjs';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

const AUTH_BYPASS_URLS = ['/api/auth', '/api/refresh'];

function addToken(req: HttpRequest<unknown>, token: string): HttpRequest<unknown> {
  return req.clone({ setHeaders: { Authorization: `Bearer ${token}` } });
}

export const authInterceptor: HttpInterceptorFn = (
  req: HttpRequest<unknown>,
  next: HttpHandlerFn,
) => {
  const auth = inject(AuthService);
  const router = inject(Router);

  // Skip auth endpoints
  if (AUTH_BYPASS_URLS.some((url) => req.url.includes(url))) {
    return next(req);
  }

  const token = auth.getToken();
  const authedReq = token ? addToken(req, token) : req;

  return next(authedReq).pipe(
    catchError((err) => {
      if (err.status !== 401) {
        return throwError(() => err);
      }

      // Try to refresh the access token (refresh token is in httpOnly cookie)
      return auth.refresh().pipe(
        switchMap((response) => next(addToken(req, response.accessToken))),
        catchError(() => {
          auth.logout();
          router.navigate(['/login']);
          return throwError(() => err);
        }),
      );
    }),
  );
};
