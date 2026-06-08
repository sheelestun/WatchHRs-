import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MaintenanceComponent } from './features/maintenance/maintenance';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, MaintenanceComponent],
  template: `
    <router-outlet />
    <app-maintenance />
  `,
})
export class App {}
