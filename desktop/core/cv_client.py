import requests
import cv2
import numpy as np
import os
from typing import Optional, Tuple

class CVStorageClient:
    def __init__(self, api_url: str = "https://watchhrs.gehrman.me/api"):
        self.api_url = api_url
        self._access_token: Optional[str] = None

    def _auth_headers(self) -> dict:
        if self._access_token:
            return {"Authorization": f"Bearer {self._access_token}"}
        return {}

    def authenticate_by_photo(self, frame: np.ndarray) -> Optional[str]:
        """
        Отправить фото на /auth

        Returns:
            employeeId (str) or None
        """
        try:
            success, buffer = cv2.imencode('.jpg', frame, [cv2.IMWRITE_JPEG_QUALITY, 85])
            if not success:
                return None

            files = {'photo': ('employee.jpg', buffer.tobytes(), 'image/jpeg')} #from desktop-add
            #from develop files = {'file': ('employee.jpg', buffer.tobytes(), 'image/jpeg')}
            
            response = requests.post(
                f"{self.api_url}/auth",
                files=files,
                timeout=10
            )

            # from desktop-add
            print(f"🔍 Auth response: {response.status_code} {response.text[:300]}")
            if response.status_code == 200:
                data = response.json()
                self._access_token = data.get('accessToken')
                return data.get('userID')

            ## from develop
            #print(f"🔍 Auth response: {response.status_code} {response.text[:300]}")
            #if response.status_code == 200:
            #    data = response.json()
            #    return data.get('employeeId')

            return None

        except Exception as e:
            print(f"❌ CV Auth Error: {e}")
            return None

    def send_statistics(self, employee_id: str, mouse_clicks: int, keyboard_clicks: int) -> Optional[str]:
        """
        POST {address}/statistic/{employeeId}
        Returns: screenshotId
        """
        try:
            payload = {
                "count_mouse_clicks": mouse_clicks,
                "count_keyboard_clicks": keyboard_clicks
            }
            response = requests.post(
                f"{self.api_url}/statistic/{employee_id}",
                json=payload,
                headers=self._auth_headers(),
                timeout=10
            )
            if response.status_code == 200:
                return response.json().get("screenshotId")
            print(f"❌ Stats error: {response.status_code} {response.text[:200]}")
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
                files = {'screenshot': (filename, f, 'image/png')}
                response = requests.post(
                    f"{self.api_url}/screenshot/{employee_id}",
                    files=files,
                    headers=self._auth_headers(),
                    timeout=20
                )
                return response.status_code == 200
        except Exception as e:
            print(f"❌ Screenshot upload error: {e}")
            return False

    def start_work_session(self, employee_id: str):
        """POST {address}/work_session/{employeeId}/start"""
        try:
            requests.post(
                f"{self.api_url}/work_session/{employee_id}/start",
                headers=self._auth_headers(),
                timeout=10,
            )
        except Exception as e:
            print(f"❌ Start session error: {e}")

    def stop_work_session(self, employee_id: str):
        """POST {address}/work_session/{employeeId}/stop"""
        try:
            requests.post(
                f"{self.api_url}/work_session/{employee_id}/stop",
                headers=self._auth_headers(),
                timeout=10,
            )
        except Exception as e:
            print(f"❌ Stop session error: {e}")
