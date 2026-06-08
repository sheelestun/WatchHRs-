import { Component, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { EmployeeService } from '../../core/services/employee.service';
import { AuthService } from '../../core/services/auth.service';
import { Employee } from '../../core/models/employee.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    MatTableModule,
    MatButtonModule,
    MatToolbarModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTooltipModule,
  ],
  templateUrl: './dashboard.html',
})
export class DashboardComponent implements OnInit {
  employees = signal<Employee[]>([]);
  loading = signal(true);
  displayedColumns = ['name', 'email', 'actions'];

  constructor(
    private employeeService: EmployeeService,
    private authService: AuthService,
    private router: Router,
    private snackBar: MatSnackBar,
  ) {}

  ngOnInit(): void {
    this.loadEmployees();
  }

  loadEmployees(): void {
    this.loading.set(true);
    this.employeeService.getAll().subscribe({
      next: (list) => {
        this.employees.set(list ?? []);
        this.loading.set(false);
      },
      error: () => {
        this.snackBar.open('Failed to load employees', 'Dismiss', { duration: 3000 });
        this.loading.set(false);
      },
    });
  }

  viewEmployee(id: string): void {
    this.router.navigate(['/employees', id]);
  }

  addEmployee(): void {
    this.router.navigate(['/employees/new']);
  }

  deleteEmployee(employee: Employee, event: Event): void {
    event.stopPropagation();
    if (!confirm(`Delete ${employee.name}? This cannot be undone.`)) return;

    this.employeeService.delete(employee.id).subscribe({
      next: () => {
        this.employees.update((list) => list.filter((e) => e.id !== employee.id));
        this.snackBar.open('Employee deleted', 'OK', { duration: 2000 });
      },
      error: () => {
        this.snackBar.open('Failed to delete employee', 'Dismiss', { duration: 3000 });
      },
    });
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}
