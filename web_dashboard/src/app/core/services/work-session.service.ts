import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { WorkSessionsResponse } from '../models/work-session.model';

@Injectable({ providedIn: 'root' })
export class WorkSessionService {
  constructor(private http: HttpClient) {}

  /** date must be YYYY-MM-DD */
  getByDate(employeeId: string, date: string): Observable<WorkSessionsResponse> {
    return this.http.get<WorkSessionsResponse>(`/api/work_session/${employeeId}/${date}`);
  }
}
