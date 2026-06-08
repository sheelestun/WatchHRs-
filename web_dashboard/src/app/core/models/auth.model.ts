export interface LoginRequest {
  userID: string;
}

export interface AuthResponse {
  userID: string;
  role: 'manager' | 'employee';
  accessToken: string;
}
