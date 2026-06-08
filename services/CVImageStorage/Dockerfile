FROM python:3.12-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1 \
    libglib2.0-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY storage+recognition ./storage+recognition

WORKDIR /app/storage+recognition/api

ENV STORAGE_ROOT=/data/storage
ENV HOST=0.0.0.0
ENV PORT=8081

EXPOSE 8081

CMD ["python", "main.py"]
