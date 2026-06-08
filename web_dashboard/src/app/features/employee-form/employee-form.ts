import { Component, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatStepperModule } from '@angular/material/stepper';
import { EmployeeService } from '../../core/services/employee.service';
import { AuthService } from '../../core/services/auth.service';

type Step = 'form' | 'photo' | 'done';

@Component({
  selector: 'app-employee-form',
  standalone: true,
  imports: [
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatToolbarModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTooltipModule,
    MatStepperModule,
  ],
  templateUrl: './employee-form.html',
})
export class EmployeeFormComponent {
  // Step 1: employee info
  name = '';
  email = '';
  loading = signal(false);
  error = signal('');

  // Step 2: photo upload
  createdEmployeeId = signal('');
  photoFile: File | null = null;
  photoLoading = signal(false);

  step = signal<Step>('form');

  constructor(
    private employeeService: EmployeeService,
    private authService: AuthService,
    private router: Router,
    private snackBar: MatSnackBar,
  ) {}

  submitForm(): void {
    if (!this.name.trim() || !this.email.trim()) return;
    this.loading.set(true);
    this.error.set('');

    const managerID = this.authService.getManagerId() ?? '';
    this.employeeService
      .create({ name: this.name.trim(), email: this.email.trim(), managerID })
      .subscribe({
        next: (res) => {
          this.loading.set(false);
          this.createdEmployeeId.set(res.employeeId);
          this.step.set('photo');
        },
        error: () => {
          this.loading.set(false);
          this.error.set('Failed to create employee. Please try again.');
        },
      });
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.photoFile = input.files?.[0] ?? null;
  }

  uploadPhoto(): void {
    if (!this.photoFile) return;
    this.photoLoading.set(true);

    this.employeeService.uploadPhoto(this.createdEmployeeId(), this.photoFile).subscribe({
      next: () => {
        this.photoLoading.set(false);
        this.step.set('done');
      },
      error: () => {
        this.photoLoading.set(false);
        this.snackBar.open('Photo upload failed. You can upload it later.', 'OK', {
          duration: 4000,
        });
        this.step.set('done');
      },
    });
  }

  skipPhoto(): void {
    this.step.set('done');
  }

  goToDashboard(): void {
    this.router.navigate(['/dashboard']);
  }

  goBack(): void {
    this.router.navigate(['/dashboard']);
  }
}
