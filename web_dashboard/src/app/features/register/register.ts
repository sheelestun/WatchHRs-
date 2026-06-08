import { Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [
    FormsModule,
    RouterLink,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressSpinnerModule,
  ],
  templateUrl: './register.html',
})
export class RegisterComponent {
  name = '';
  email = '';
  loading = signal(false);
  error = signal('');
  managerId = signal('');

  constructor(
    private auth: AuthService,
    private router: Router,
  ) {}

  register(): void {
    if (!this.name.trim() || !this.email.trim()) return;
    this.loading.set(true);
    this.error.set('');

    this.auth.register(this.name.trim(), this.email.trim()).subscribe({
      next: (response) => {
        this.loading.set(false);
        this.managerId.set(response.managerId);
      },
      error: () => {
        this.loading.set(false);
        this.error.set('Registration failed. Please try again.');
      },
    });
  }
}
