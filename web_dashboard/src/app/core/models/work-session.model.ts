export interface WorkSession {
  id: string;
  employeeID: string;
  start_time: string;
  end_time: string | null;
  total_time: string | null;
}

export interface WorkSessionsResponse {
  workSessions: WorkSession[];
}
