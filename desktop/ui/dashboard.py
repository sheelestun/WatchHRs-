import tkinter as tk
from tkinter import messagebox
from core import states

class DashboardView:
    def __init__(self, parent, employee_id, on_toggle_work, on_logout, on_work_finished):
        self.frame = tk.Frame(parent)
        self.frame.pack(fill="both", expand=True, padx=20, pady=20)

        self.employee_id = employee_id
        self.work_status = states.NOT_WORKING
        self.on_toggle_work = on_toggle_work
        self.on_logout = on_logout
        self.on_work_finished = on_work_finished

        self._build_ui()

    def _build_ui(self):
        # Приветствие
        tk.Label(self.frame, text=f"Сотрудник ID: {self.employee_id}", font=("Arial", 12)).pack(pady=5)
        tk.Label(self.frame, text=f"Статус: Рабочая сессия", font=("Arial", 14, "bold")).pack(pady=10)

        # Кнопка начала/завершения дня
        self.toggle_btn = tk.Button(
            self.frame, text="Начать рабочий день", font=("Arial", 12), 
            width=20, bg="#ccffcc", command=self._toggle_work
        )
        self.toggle_btn.pack(pady=10)

        # Кнопка выхода
        tk.Button(self.frame, text="Выйти", font=("Arial", 10), command=self.on_logout).pack(pady=5)

        # Метка для статистики
        self.stats_label = tk.Label(self.frame, text="", font=("Arial", 10), justify="left")
        self.stats_label.pack(pady=10)

    def _toggle_work(self):
        if self.work_status == states.WORKING:
            if not messagebox.askyesno("Подтверждение", "Завершить работу?"):
                return
            
            stats = self.on_work_finished()
            self._show_stats(stats)
            
            self.work_status = states.NOT_WORKING
            self.toggle_btn.config(text="Начать рабочий день", bg="#ccffcc")
        else:
            self.work_status = states.WORKING
            self.toggle_btn.config(text="Закончить работу", bg="#ffcccc")
            self.stats_label.config(text="Мониторинг запущен...")
        
        self.on_toggle_work(self.work_status)

    def _show_stats(self, stats):
        text = (
            f"✅ Сессия завершена\n\n"
            f"Мышь: {stats['mouse']}\n"
            f"Клавиатура: {stats['keyboard']}\n"
            f"Скриншотов сделано: {stats['screenshots']}"
        )
        self.stats_label.config(text=text, fg="green")
