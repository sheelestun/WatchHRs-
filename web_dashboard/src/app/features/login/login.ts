import { Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressSpinnerModule,
  ],
  templateUrl: './login.html',
})
export class LoginComponent {
  userId = '';
  loading = signal(false);
  error = signal('');

  constructor(
    private auth: AuthService,
    private router: Router,
  ) {}

  login(): void {
    if (!this.userId.trim()) return;
    this.loading.set(true);
    this.error.set('');

    this.auth.login(this.userId.trim()).subscribe({
      next: (response) => {
        this.loading.set(false);
        if (response.role !== 'manager') {
          this.error.set('Access denied: only managers can use this dashboard.');
          this.auth.logout();
          return;
        }
        this.router.navigate(['/dashboard']);
      },
      error: () => {
        this.loading.set(false);
        this.error.set('Login failed. Please check your Manager ID and try again.');
      },
    });
  }
}
