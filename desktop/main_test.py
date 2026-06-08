import sys
import os
from main import main

if __name__ == "__main__":
    # Принудительно добавляем флаг --test для тестового режима
    if "--test" not in sys.argv:
        sys.argv.append("--test")
    main()
