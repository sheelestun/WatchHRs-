import { Component, inject } from '@angular/core';
import { MaintenanceService } from '../../core/services/maintenance.service';

@Component({
  selector: 'app-maintenance',
  standalone: true,
  styles: [`
    .overlay {
      position: fixed;
      inset: 0;
      z-index: 9999;
      background: #0f0f0f;
      color: #e0e0e0;
      font-family: system-ui, -apple-system, sans-serif;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
    }

    h1 {
      font-size: clamp(1.2rem, 4vw, 2rem);
      font-weight: 500;
      margin-bottom: 1.5rem;
      text-align: center;
      padding: 0 1rem;
    }

    img {
      width: 100vw;
      height: auto;
      display: block;
    }
  `],
  template: `
    @if (maintenance.offline()) {
      <div class="overlay">
        <h1>We are down for maintenance, my bad</h1>
        <img src="assets/kemps.webp" alt="maintenance gif">
      </div>
    }
  `,
})
export class MaintenanceComponent {
  readonly maintenance = inject(MaintenanceService);
}
