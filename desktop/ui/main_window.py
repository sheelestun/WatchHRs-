import tkinter as tk
import threading
import os
from PIL import Image, ImageTk
from core import auth_placeholder, states
from core.activity_monitor import ActivityMonitor
from core.camera_capture import CameraCapture
from core.cv_client import CVStorageClient
from ui.dashboard import DashboardView

class MainWindow:
    def __init__(self, root, test_mode=False):
        self.root = root
        self.root.title("WatchHRs Desktop")
        self.root.geometry("400x450")
        
        self.test_mode = test_mode
        self.client = CVStorageClient(api_url="https://watchhrs.gehrman.me/api")
        self.employee_id = None
        
        # Инициализируем монитор с коллбэком для отправки данных каждые 10 мин
        self.monitor = ActivityMonitor(on_report_callback=self._send_periodic_report)
        self.is_working = False
        self.camera = CameraCapture()

        self.container = tk.Frame(root)
        self.container.pack(fill="both", expand=True)

        self._show_login()
        self.root.protocol("WM_DELETE_WINDOW", self._on_closing)

    def _send_periodic_report(self, screenshot_path, mouse_clicks, keyboard_clicks):
        """Вызывается монитором каждые 10 минут"""
        if not self.employee_id:
            return
            
        print(f"📊 Отправка статистики: M:{mouse_clicks}, K:{keyboard_clicks}")
        # 1. Отправляем статистику и получаем screenshotId
        screenshot_id = self.client.send_statistics(self.employee_id, mouse_clicks, keyboard_clicks)
        
        if screenshot_id:
            print(f"📸 Загрузка скриншота {screenshot_id}...")
            # 2. Загружаем сам скриншот
            success = self.client.upload_screenshot(self.employee_id, screenshot_id, screenshot_path)
            if success:
                print("✅ Скриншот успешно загружен")
                # Удаляем временный файл
                try:
                    os.remove(screenshot_path)
                except:
                    pass
            else:
                print("❌ Ошибка загрузки скриншота")
        else:
            print("❌ Ошибка получения screenshotId")

    def _on_closing(self):
        if self.is_working:
            from tkinter import messagebox
            if messagebox.askyesno("Выход", "Завершить рабочий день и выйти?"):
                self._stop_work()
                self.root.destroy()
        else:
            self.camera.stop()
            self.root.destroy()

    def _show_login(self):
        self._clear_screen()
        tk.Label(self.container, text="Авторизация сотрудника", font=("Arial", 14)).pack(pady=20)
        
        self.video_container = tk.Frame(self.container, bg="black", width=320, height=240)
        self.video_container.pack(pady=10)
        self.video_container.pack_propagate(False)
        
        self.video_label = tk.Label(self.video_container, bg="black")
        self.video_label.pack(expand=True)
        
        self.status_lbl = tk.Label(self.container, text="Нажмите 'Войти' для сканирования")
        self.status_lbl.pack(pady=10)
        
        self.auth_btn = tk.Button(self.container, text="Войти", command=self._start_auth, width=15, height=2)
        self.auth_btn.pack(pady=10)

    def _start_auth(self):
        self.auth_btn.config(state="disabled")
        self.status_lbl.config(text="Поиск лица...", fg="blue")
        
        if self.test_mode:
            # Имитируем задержку для "типа работы"
            self.root.after(1000, lambda: self._on_auth_result(states.SUCCESS, "test-employee-uuid-123"))
            return

        if self.camera.start():
            self.camera.set_frame_callback(self._update_video_frame)
        auth_placeholder.check_photo(self._on_auth_result, self.camera, self.client)

    def _update_video_frame(self, frame):
        img = Image.fromarray(frame)
        img = img.resize((320, 240), Image.Resampling.LANCZOS)
        photo = ImageTk.PhotoImage(img)
        self.video_label.img = photo
        self.video_label.config(image=photo)

    def _on_auth_result(self, status, employee_id):
        self.root.after(0, lambda: self._handle_auth(status, employee_id))

    def _handle_auth(self, status, employee_id):
        if status == states.SUCCESS:
            self.employee_id = employee_id
            self.camera.stop()
            self._show_dashboard()
        else:
            self.status_lbl.config(text="Ошибка авторизации. Попробуйте снова.", fg="red")
            self.auth_btn.config(state="normal")

    def _show_dashboard(self):
        self._clear_screen()
        self.dashboard = DashboardView(
            self.container, 
            self.employee_id,
            on_toggle_work=self._on_toggle_work,
            on_logout=self._logout,
            on_work_finished=self._stop_work
        )

    def _on_toggle_work(self, status):
        if status == states.WORKING:
            self.is_working = True
            self.client.start_work_session(self.employee_id)
            # Передаем интервал: 10 секунд для теста, 600 для продакшна
            interval = 10 if self.test_mode else 600
            self.monitor.start(interval=interval)
        else:
            self._stop_work()

    def _stop_work(self):
        self.is_working = False
        self.client.stop_work_session(self.employee_id)
        return self.monitor.stop()

    def _logout(self):
        if self.is_working:
            self._stop_work()
        self.employee_id = None
        self._show_login()

    def _clear_screen(self):
        for widget in self.container.winfo_children():
            widget.destroy()
