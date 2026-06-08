import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuthService } from './auth.service';
import { CreateEmployeeRequest, CreateEmployeeResponse, Employee } from '../models/employee.model';

@Injectable({ providedIn: 'root' })
export class EmployeeService {
  constructor(
    private http: HttpClient,
    private auth: AuthService,
  ) {}

  getAll(): Observable<Employee[]> {
    const managerId = this.auth.getManagerId();
    return this.http.get<Employee[]>(`/api/manager/${managerId}/employee/all`);
  }

  create(request: CreateEmployeeRequest): Observable<CreateEmployeeResponse> {
    return this.http.post<CreateEmployeeResponse>('/api/employee', request);
  }

  delete(employeeId: string): Observable<void> {
    return this.http.delete<void>(`/api/employee/${employeeId}`);
  }

  uploadPhoto(userId: string, file: File): Observable<unknown> {
    const formData = new FormData();
    formData.append('userId', userId);
    formData.append('screenshot', file);
    return this.http.post('/api/photo', formData);
  }
}
