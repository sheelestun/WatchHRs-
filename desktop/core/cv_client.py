import requests
import cv2
import numpy as np
import os
from typing import Optional, Tuple

class CVStorageClient:
    def __init__(self, api_url: str = "http://localhost:8080"): # Default to common port, or adjust as needed
        self.api_url = api_url
    
    def authenticate_by_photo(self, frame: np.ndarray) -> Optional[str]:
        """
        Отправить фото на /employee/auth
        
        Returns:
            employeeId (str) or None
        """
        try:
            success, buffer = cv2.imencode('.png', frame)
            if not success:
                return None
            
            files = {'file': ('employee.png', buffer.tobytes(), 'image/png')}
            response = requests.post(
                f"{self.api_url}/employee/auth",
                files=files,
                timeout=10
            )
            
            if response.status_code == 200:
                data = response.json()
                return data.get('employeeId')
            
            return None
            
        except Exception as e:
            print(f"❌ CV Auth Error: {e}")
            return None

    def send_statistics(self, employee_id: str, mouse_clicks: int, keyboard_clicks: int) -> Optional[str]:
        """
        POST {address}/statistic/{employeeId}/
        Returns: screenshotId
        """
        try:
            payload = {
                "count_mouse_clicks": mouse_clicks,
                "count_keyboard_clicks": keyboard_clicks
            }
            response = requests.post(
                f"{self.api_url}/statistic/{employee_id}/",
                json=payload,
                timeout=10
            )
            if response.status_code == 200:
                return response.json().get("screenshotId")
            return None
        except Exception as e:
            print(f"❌ Stats sending error: {e}")
            return None

    def upload_screenshot(self, employee_id: str, screenshot_id: str, file_path: str):
        """
        POST {address}/screenshot/{employeeId}
        Request: employeeId-screenshotId.png
        """
        try:
            filename = f"{employee_id}-{screenshot_id}.png"
            with open(file_path, 'rb') as f:
                files = {'file': (filename, f, 'image/png')}
                response = requests.post(
                    f"{self.api_url}/screenshot/{employee_id}",
                    files=files,
                    timeout=20
                )
                return response.status_code == 200
        except Exception as e:
            print(f"❌ Screenshot upload error: {e}")
            return False

    def start_work_session(self, employee_id: str):
        """POST {address}/work_session/{employeeId}/start"""
        try:
            requests.post(f"{self.api_url}/work_session/{employee_id}/start", timeout=10)
        except Exception as e:
            print(f"❌ Start session error: {e}")

    def stop_work_session(self, employee_id: str):
        """POST {address}/work_session/{employeeId}/stop"""
        try:
            requests.post(f"{self.api_url}/work_session/{employee_id}/stop", timeout=10)
        except Exception as e:
            print(f"❌ Stop session error: {e}")
