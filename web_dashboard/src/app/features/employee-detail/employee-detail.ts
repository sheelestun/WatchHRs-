import { Component, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { DomSanitizer, SafeUrl } from '@angular/platform-browser';
import { map } from 'rxjs';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatTooltipModule } from '@angular/material/tooltip';
import { WorkSessionService } from '../../core/services/work-session.service';
import { StatisticService } from '../../core/services/statistic.service';
import { ScreenshotService } from '../../core/services/screenshot.service';
import { WorkSession } from '../../core/models/work-session.model';
import { Statistic } from '../../core/models/statistic.model';
import { Screenshot } from '../../core/models/screenshot.model';

@Component({
  selector: 'app-employee-detail',
  standalone: true,
  imports: [
    FormsModule,
    MatTabsModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatToolbarModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatFormFieldModule,
    MatInputModule,
    MatTooltipModule,
  ],
  templateUrl: './employee-detail.html',
})
export class EmployeeDetailComponent implements OnInit {
  employeeId = '';
  selectedDate = new Date().toISOString().split('T')[0]; // YYYY-MM-DD

  workSessions = signal<WorkSession[]>([]);
  statistics = signal<Statistic[]>([]);
  screenshots = signal<Screenshot[]>([]);
  screenshotUrls = signal<Record<string, SafeUrl>>({});

  loadingWorkSessions = signal(false);
  loadingStats = signal(false);
  loadingScreenshots = signal(false);

  workSessionColumns = ['start_time', 'end_time', 'total_time'];
  statColumns = ['timestamp', 'count_mouse_clicks', 'count_keyboard_clicks'];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private workSessionService: WorkSessionService,
    private statisticService: StatisticService,
    private screenshotService: ScreenshotService,
    private snackBar: MatSnackBar,
    private sanitizer: DomSanitizer,
  ) {}

  ngOnInit(): void {
    this.employeeId = this.route.snapshot.paramMap.get('id') ?? '';
    this.loadData();
  }

  loadData(): void {
    this.screenshotUrls.set({});
    this.loadWorkSessions();
    this.loadStatistics();
    this.loadScreenshots();
  }

  loadWorkSessions(): void {
    this.loadingWorkSessions.set(true);
    this.workSessionService.getByDate(this.employeeId, this.selectedDate).subscribe({
      next: (res) => {
        this.workSessions.set(res.workSessions ?? []);
        this.loadingWorkSessions.set(false);
      },
      error: () => {
        this.snackBar.open('Failed to load work sessions', 'Dismiss', { duration: 3000 });
        this.loadingWorkSessions.set(false);
      },
    });
  }

  loadStatistics(): void {
    this.loadingStats.set(true);
    this.statisticService.getByDate(this.employeeId, this.selectedDate).subscribe({
      next: (res) => {
        this.statistics.set(res.screenshots ?? []);
        this.loadingStats.set(false);
      },
      error: () => {
        this.snackBar.open('Failed to load activity stats', 'Dismiss', { duration: 3000 });
        this.loadingStats.set(false);
      },
    });
  }

  loadScreenshots(): void {
    this.loadingScreenshots.set(true);
    this.screenshotService.getByDate(this.employeeId, this.selectedDate).subscribe({
      next: (res) => {
        this.screenshots.set(res.screenshots ?? []);
        this.loadingScreenshots.set(false);
        // Pre-load each image as an object URL so auth headers are sent
        (res.screenshots ?? []).forEach((s) => this.loadScreenshotImage(s));
      },
      error: () => {
        this.snackBar.open('Failed to load screenshots', 'Dismiss', { duration: 3000 });
        this.loadingScreenshots.set(false);
      },
    });
  }

  private loadScreenshotImage(screenshot: Screenshot): void {
    this.screenshotService
      .getFile(this.employeeId, screenshot.filename)
      .pipe(
        map((blob) => {
          const url = URL.createObjectURL(blob);
          return this.sanitizer.bypassSecurityTrustUrl(url);
        }),
      )
      .subscribe({
        next: (safeUrl) => {
          this.screenshotUrls.update((urls) => ({ ...urls, [screenshot.filename]: safeUrl }));
        },
      });
  }

  downloadArchive(): void {
    this.screenshotService.downloadArchive(this.employeeId, this.selectedDate).subscribe({
      next: (blob) => {
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `screenshots-${this.employeeId}-${this.selectedDate}.zip`;
        a.click();
        URL.revokeObjectURL(url);
      },
      error: () => {
        this.snackBar.open('Failed to download archive', 'Dismiss', { duration: 3000 });
      },
    });
  }

  formatTime(iso: string | null): string {
    if (!iso) return '—';
    return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  goBack(): void {
    this.router.navigate(['/dashboard']);
  }
}
