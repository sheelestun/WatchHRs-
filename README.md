# WatchHRs-

```
WatchHRs/
|   readme.md
|   services/web/
|   services/desktop/
|   services/API Gateway/
```

# Workflow

![gitflow graph](/assets/gitflow.jpg)

Наименования веток:
- `main` — то что работает, то что распространяется
- `develop` — форк main; то над чем сейчас работаем
- `feature/*` — форк develop; реализация конкретной фичи

Доп названия веток:
- `release` — форк develop; этап "чистки" всего проекта, перед тем как пушить в main
- `hotfix` — форк main; починка багов

