import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ScreenshotsResponse } from '../models/screenshot.model';

@Injectable({ providedIn: 'root' })
export class ScreenshotService {
  constructor(private http: HttpClient) {}

  /** date must be YYYY-MM-DD */
  getByDate(employeeId: string, date: string): Observable<ScreenshotsResponse> {
    return this.http.get<ScreenshotsResponse>(`/api/screenshot/${employeeId}/${date}`);
  }

  /** Returns the image as a Blob so it can be displayed with auth headers applied */
  getFile(employeeId: string, filename: string): Observable<Blob> {
    return this.http.get(`/api/screenshot/${employeeId}/file/${filename}`, {
      responseType: 'blob',
    });
  }

  /** Downloads the full-day archive as a Blob */
  downloadArchive(employeeId: string, date: string): Observable<Blob> {
    return this.http.get(`/api/screenshot/${employeeId}/${date}/archive`, {
      responseType: 'blob',
    });
  }
}
