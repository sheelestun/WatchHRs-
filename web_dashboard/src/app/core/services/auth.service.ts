import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { AuthResponse, LoginRequest } from '../models/auth.model';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private accessToken: string | null = null;
  private userId: string | null = null;

  readonly isLoggedIn = signal(false);

  constructor(private http: HttpClient) {}

  login(userID: string): Observable<AuthResponse> {
    const body: LoginRequest = { userID };
    return this.http.post<AuthResponse>('/api/auth', body).pipe(
      tap((response) => {
        this.accessToken = response.accessToken;
        this.userId = response.userID;
        this.isLoggedIn.set(true);
      }),
    );
  }

  refresh(): Observable<AuthResponse> {
    return this.http
      .post<AuthResponse>('/api/refresh', {}, { withCredentials: true })
      .pipe(
        tap((response) => {
          this.accessToken = response.accessToken;
          this.userId = response.userID;
          this.isLoggedIn.set(true);
        }),
      );
  }

  logout(): void {
    this.accessToken = null;
    this.userId = null;
    this.isLoggedIn.set(false);
  }

  getToken(): string | null {
    return this.accessToken;
  }

  /** The userID is the managerId when role === 'manager' */
  getManagerId(): string | null {
    return this.userId;
  }

  isAuthenticated(): boolean {
    return !!this.accessToken;
  }
}
