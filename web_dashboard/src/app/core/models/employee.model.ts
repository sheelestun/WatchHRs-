export interface Employee {
  id: string;
  name: string;
  email: string;
  managerID: string;
}

export interface CreateEmployeeRequest {
  name: string;
  email: string;
  managerID: string;
}

export interface CreateEmployeeResponse {
  employeeId: string;
}
