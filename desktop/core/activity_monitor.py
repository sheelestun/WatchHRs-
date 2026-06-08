from pynput import mouse, keyboard
from PIL import ImageGrab
import time
import threading
import os

class ActivityMonitor:
    def __init__(self, on_report_callback=None):
        self.count_mouse = 0
        self.count_keyboard = 0
        self.count_screenshot = 0
        self.running = False
        
        self.on_report_callback = on_report_callback
        self.m_listener = None
        self.k_listener = None
        self.screenshot_thread = None
        self._stop_event = threading.Event()

    def _click_mouse(self, x, y, button, pressed):
        if pressed:
            self.count_mouse += 1

    def _click_keyboard(self, key):
        self.count_keyboard += 1

    def _take_screenshot(self):
        if not self.running:
            return
        
        self.count_screenshot += 1
        os.makedirs("screenshots", exist_ok=True)
        filename = os.path.join("screenshots", f"temp_screen_{int(time.time())}.png")
        
        try:
            screenshot = ImageGrab.grab()
            screenshot.save(filename)
            
            # Если есть коллбэк - сообщаем о новом скриншоте и текущей статистике
            if self.on_report_callback:
                # Передаем копии значений, чтобы сбросить их для следующего интервала
                m = self.count_mouse
                k = self.count_keyboard
                
                # Сбрасываем счетчики для нового интервала (согласно функциональным требованиям)
                self.count_mouse = 0
                self.count_keyboard = 0
                
                # Вызываем в отдельном потоке, чтобы не блокировать таймер
                threading.Thread(
                    target=self.on_report_callback, 
                    args=(filename, m, k),
                    daemon=True
                ).start()
                
        except Exception as e:
            print(f"Screenshot error: {e}")

    def _screenshot_timer(self, interval):
        while self.running and not self._stop_event.is_set():
            # Ждем интервал (например 600 сек = 10 мин)
            if self._stop_event.wait(interval):
                break
            if self.running:
                self._take_screenshot()

    def start(self, interval=600):
        """Запустить мониторинг. interval в секундах (дефолт 10 мин)"""
        self.running = True
        self.count_mouse = 0
        self.count_keyboard = 0
        self.count_screenshot = 0
        self._stop_event.clear()
        
        self.m_listener = mouse.Listener(on_click=self._click_mouse)
        self.k_listener = keyboard.Listener(on_press=self._click_keyboard)
        
        self.m_listener.start()
        self.k_listener.start()
        
        self.screenshot_thread = threading.Thread(
            target=self._screenshot_timer, 
            args=(interval,),
            daemon=True
        )
        self.screenshot_thread.start()

    def stop(self):
        """Остановить мониторинг и вернуть финальную статистику"""
        self.running = False
        self._stop_event.set()
        
        if self.m_listener:
            self.m_listener.stop()
        if self.k_listener:
            self.k_listener.stop()
        
        return {
            "mouse": self.count_mouse,
            "keyboard": self.count_keyboard,
            "screenshots": self.count_screenshot
        }
