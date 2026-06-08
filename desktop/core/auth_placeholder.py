import threading
import numpy as np
from core import states
from core.cv_client import CVStorageClient
from core.camera_capture import CameraCapture

cv_client = CVStorageClient(api_url="https://watchhrs.gehrman.me/api/")

# Переменная для хранения текущего ID пользователя
current_employee_id = None

def check_photo(callback, camera: CameraCapture):
    """
    Реальная авторизация через CVImageStorage
    """
    def auth_thread():
        global current_employee_id
        try:
            best_frame = None
            frames_captured = 0
            
            # Пытаемся захватить четкий кадр
            while frames_captured < 30: 
                if hasattr(camera, '_last_frame'):
                    best_frame = camera._last_frame
                    if best_frame is not None:
                        break
                threading.Event().wait(0.12)
                frames_captured += 1
            
            if best_frame is None:
                callback(states.ERROR, None)
                return
            
            # Вызываем API
            employee_id = cv_client.authenticate_by_photo(best_frame)
            
            if employee_id:
                current_employee_id = employee_id
                print(f"✅ Авторизован, ID: {employee_id}")
                callback(states.SUCCESS, employee_id)
            else:
                print("❌ Авторизация не удалась")
                callback(states.ERROR, None)
                
        except Exception as e:
            print(f"❌ Auth error: {e}")
            callback(states.ERROR, None)
    
    thread = threading.Thread(target=auth_thread, daemon=True)
    thread.start()
