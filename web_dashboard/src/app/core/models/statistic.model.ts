export interface Statistic {
  id: string;
  employeeID: string;
  count_mouse_clicks: number;
  count_keyboard_clicks: number;
  timestamp: string;
}

export interface StatisticsResponse {
  screenshots: Statistic[];
}
