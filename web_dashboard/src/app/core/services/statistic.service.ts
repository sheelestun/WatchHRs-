import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { StatisticsResponse } from '../models/statistic.model';

@Injectable({ providedIn: 'root' })
export class StatisticService {
  constructor(private http: HttpClient) {}

  /** date must be YYYY-MM-DD */
  getByDate(employeeId: string, date: string): Observable<StatisticsResponse> {
    return this.http.get<StatisticsResponse>(`/api/statistic/${employeeId}/${date}`);
  }
}
