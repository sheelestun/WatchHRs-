# How to start application easily?

## 1. Clone repository

```commandline
git clone https://github.com/sheelestun/WatchHRs-.git
```

## 2. Change main branch
```commandline
git checkout remotes/origin/develop
```

## 3. Rename .env
### There is an **.env.example** file in repo - we need to change it for docker-compose:
```commandline
mv .env.example .env
```

## 4. Start docker engine + start docker-compose
```commandline
docker compose up --build
```

## How to stop container?
```commandline
docker compose down
```
