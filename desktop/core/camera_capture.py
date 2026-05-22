import cv2
import threading

class CameraCapture:
    def __init__(self, video_source=0):
        self.video_source = video_source
        self.cap = cv2.VideoCapture(self.video_source)
        self.running = False
        self.frame_callback = None
        self._thread = None

    def set_frame_callback(self, callback):
        self.frame_callback = callback

    def _update_frame(self):
        while self.running:
            ret, frame = self.cap.read()
            if ret:
                self._last_frame = frame  # ✅ Сохраняем последний кадр для анализа
                if self.frame_callback:
                    frame_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
                    self.frame_callback(frame_rgb)

    def start(self):
        if self.running:
            return True
        
        # ✅ Если камера была закрыта, открываем заново
        if not self.cap.isOpened():
            self.cap = cv2.VideoCapture(self.video_source)
            
        if not self.cap.isOpened():
            print("Ошибка: камера не найдена")
            return False
            
        self.running = True
        self._thread = threading.Thread(target=self._update_frame, daemon=True)
        self._thread.start()
        return True

    def stop(self):
        self.running = False
        if self._thread:
            self._thread.join(timeout=1)
        if self.cap.isOpened():
            self.cap.release()
