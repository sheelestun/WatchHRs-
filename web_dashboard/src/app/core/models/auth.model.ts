export interface LoginRequest {
  userID: string;
}

export interface AuthResponse {
  userID: string;
  role: 'manager' | 'employee';
  accessToken: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
}

export interface RegisterResponse {
  managerId: string;
}
