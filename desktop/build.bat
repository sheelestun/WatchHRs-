@echo off
echo Building WatchHRs Standard...
python -m PyInstaller --noconfirm --onefile --windowed --name "WatchHRs" main.py

echo Building WatchHRs Test...
python -m PyInstaller --noconfirm --onefile --windowed --name "WatchHRs_Test" main_test.py

echo Done! Executables are in the 'dist' folder.
pause
