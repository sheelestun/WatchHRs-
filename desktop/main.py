import tkinter as tk
import sys
import os

# Добавляем текущую директорию в путь, чтобы импорты core и ui работали
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from ui.main_window import MainWindow

import argparse

def main():
    parser = argparse.ArgumentParser(description="WatchHRs Desktop App")
    parser.add_argument("--test", action="store_true", help="Enable test mode with auto-auth")
    args = parser.parse_args()

    root = tk.Tk()
    app = MainWindow(root, test_mode=args.test)
    root.mainloop()

if __name__ == "__main__":
    main()
