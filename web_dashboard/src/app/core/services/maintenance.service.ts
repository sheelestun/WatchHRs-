import { Injectable, signal } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class MaintenanceService {
  readonly offline = signal(false);

  setOffline(value: boolean): void {
    this.offline.set(value);
  }
}
