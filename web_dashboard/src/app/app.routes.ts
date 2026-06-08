import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
  {
    path: 'login',
    loadComponent: () => import('./features/login/login').then((m) => m.LoginComponent),
  },
  {
    path: 'register',
    loadComponent: () => import('./features/register/register').then((m) => m.RegisterComponent),
  },
  {
    path: 'dashboard',
    loadComponent: () =>
      import('./features/dashboard/dashboard').then((m) => m.DashboardComponent),
    canActivate: [authGuard],
  },
  {
    path: 'employees/new',
    loadComponent: () =>
      import('./features/employee-form/employee-form').then((m) => m.EmployeeFormComponent),
    canActivate: [authGuard],
  },
  {
    path: 'employees/:id',
    loadComponent: () =>
      import('./features/employee-detail/employee-detail').then(
        (m) => m.EmployeeDetailComponent,
      ),
    canActivate: [authGuard],
  },
  { path: '**', redirectTo: 'dashboard' },
];
