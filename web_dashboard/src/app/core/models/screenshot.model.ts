export interface Screenshot {
  id: string;
  filename: string;
  path: string;
  cnt_mouse_clicks: number;
  cnt_keyboard_clicks: number;
  timestamp: string;
}

export interface ScreenshotsResponse {
  employee_id: string;
  count: number;
  screenshots: Screenshot[];
}
