# Git Workflow (Lightweight GitFlow) — правила для репозитория

Цель: быстрые итерации без хаоса — предсказуемые ветки, единый стиль коммитов, PR как единая точка контроля качества.

Основа:
- Ветвление: GitFlow (упрощённый) (см. “A successful Git branching model”, Vincent Driessen)
- Коммиты: Conventional Commits
- Версии: Semantic Versioning (SemVer)

---

## 0) Термины и роли

- **main** — стабильная ветка релизов (то, что можно “выкатывать”).
- **develop** — интеграционная ветка (то, что готовится к следующему релизу).
- **Автор** — тот, кто делает изменения в своей ветке и открывает PR.
- **Ревьюер** — второй разработчик, который смотрит PR (минимум 1 апрув).

> Если у репозитория исторически `master`, то:
> - либо переименуйте в `main` (рекомендуется),
> - либо воспринимайте `master` как `main` и следуйте правилам дальше без изменений.

---

## 1) Политика веток (Branch Policy)

### Защищённые ветки (Protected branches)
**main** и **develop** должны быть защищены:
- запрет прямых push (только через PR),
- обязательный CI (минимум: сборка + тесты),
- минимум 1 approval перед merge,
- (опционально) линейная история (rebase/squash) — упрощает откаты.

### Типы веток
Ветки создаём только от правильной “базы”:

1) `feature/<ticket>-<kebab>`  
   - Новая функциональность.
   - База: `develop`
2) `bugfix/<ticket>-<kebab>`  
   - Багфикс НЕ в проде (найден на develop/в процессе).
   - База: `develop`
3) `release/<major>.<minor>.<patch>`  
   - Подготовка релиза (стабилизация, правки версий/доков).
   - База: `develop`
4) `hotfix/<ticket>-<kebab>`  
   - Срочный фикс в проде.
   - База: `main` (и потом обязательно возвращаем в `develop`).

### Примеры имён
- `feature/123-upload-endpoint`
- `bugfix/141-null-pointer-zip`
- `release/1.4.0`
- `hotfix/155-crash-on-startup`

---

## 2) Ежедневный рабочий цикл (самый частый сценарий)

### 2.1 Старт задачи
1. Обнови локальные ветки:
   ```bash
   git checkout develop
   git pull --rebase origin develop
   ```
2. Создай ветку:
   ```bash
   git checkout -b feature/123-upload-endpoint
   ```
3. Работай небольшими порциями (small batches):
   - коммиты логически цельные,
   - PR небольшой (условно до ~200–400 строк диффа; если больше — режем на части).

### 2.2 Коммиты (обязательно)
Используем **Conventional Commits**.

Формат:
```
<type>(<scope>): <subject>
```

Типы:
- `feat` — новая фича
- `fix` — багфикс
- `docs` — документация
- `refactor` — рефактор без изменения поведения
- `test` — тесты
- `chore` — инфраструктура/зависимости/рутина
- `ci` — пайплайны
- `perf` — оптимизация

Правила:
- `subject` в повелительном наклонении, без точки, кратко.
- При необходимости — тело коммита с мотивацией и рисками.
- Если есть “ломающее” изменение — `!` или `BREAKING CHANGE:`.

Примеры:
- `feat(api): add init upload endpoint`
- `fix(worker): handle empty pdf pages`
- `refactor(storage): extract minio client factory`
- `docs: update local dev instructions`

### 2.3 Push (порядок и частота)
- **Пушим только свою ветку**, не `develop/main`.
- Пушим часто (чтобы не потерять работу и включить CI):
  ```bash
  git push -u origin feature/123-upload-endpoint
  ```

### 2.4 Pull Request
PR открываем **в develop**:
- `feature/*` → `develop`
- `bugfix/*` → `develop`

Правила PR:
- PR = одна цель (одна фича/один фикс).
- Описание должно отвечать на:
  1) Что сделано?
  2) Почему?
  3) Как проверить?
  4) Риски/миграции/флаги?
- До ревью: актуализируй ветку:
  ```bash
  git fetch origin
  git rebase origin/develop
  # реши конфликты, если есть
  git push --force-with-lease
  ```
  `--force-with-lease` допускается **только** в своих feature/bugfix ветках.

Ревью:
- минимум 1 approval,
- CI зелёный,
- автор исправляет замечания и сам разрешает конфликты.

### 2.5 Merge стратегия
Рекомендуемая стратегия для скорости и читаемой истории:
- **Squash merge** для `feature/*` и `bugfix/*` (1 PR = 1 коммит в develop/main).
- Заголовок squash-коммита тоже в Conventional Commits стиле:
  - например: `feat(api): add init upload endpoint (#123)`

После merge:
- удалить ветку на сервере (GitHub “Delete branch”).

---

## 3) Релизный цикл (release/*)

### 3.1 Когда создаём release/*
Когда develop накопил изменения, которые хотим выпустить.

Шаги:
1. Создать ветку релиза от develop:
   ```bash
   git checkout develop
   git pull --rebase origin develop
   git checkout -b release/1.4.0
   git push -u origin release/1.4.0
   ```
2. В release-ветке:
   - обновить версию (где принято: tag/файл/переменные),
   - обновить changelog (если ведёте),
   - только стабилизация и правки релиза (без новых фич).

3. Открыть PR:
   - `release/1.4.0` → `main`
   - после мержа в `main` **обязательно** вернуть изменения в `develop` (PR `main` → `develop` или merge `release/*` → `develop`).

### 3.2 Теги релиза
Тегируем релиз в `main` по SemVer: `vMAJOR.MINOR.PATCH`.

Пример:
```bash
git checkout main
git pull origin main
git tag -a v1.4.0 -m "Release v1.4.0"
git push origin v1.4.0
```

---

## 4) Срочный фикс в проде (hotfix/*)

Сценарий: в `main` критичный баг.

Шаги:
1. Создать hotfix от main:
   ```bash
   git checkout main
   git pull --rebase origin main
   git checkout -b hotfix/155-crash-on-startup
   git push -u origin hotfix/155-crash-on-startup
   ```
2. Коммиты по правилам, PR в `main`.
3. После мержа в `main`:
   - создать тег `vX.Y.(Z+1)` по SemVer,
   - вернуть фикс в `develop` (PR `main` → `develop` или cherry-pick).

---

## 5) Минимальные правила качества (Definition of Done)

PR считается готовым, если:
- CI зелёный (build/test).
- Нет “закомментированного мусора”, временных логов, дебажных файлов.
- Для критичных изменений: добавлены/обновлены тесты (где применимо).
- Документация обновлена, если менялись интерфейсы/энвы/эндпоинты.
- Нет секретов в репозитории (ключи, токены).

---

## 6) Быстрые правила “как в стартапе” (скорость без бардака)

- Маленькие PR, частые мерджи.
- Draft PR как можно раньше (видно прогресс, можно обсуждать заранее).
- Не копим конфликтов: ребейзимся на develop минимум раз в день.
- Если задача > 1–2 дня — режем на подзадачи и PR-инкременты.
- Никаких “giant commit: final version”.

---

## 7) Рекомендованные настройки Git (локально)

Один раз:
```bash
git config --global pull.rebase true
git config --global fetch.prune true
git config --global rebase.autoStash true
```

---

## 8) Шпаргалка по командам

Создать фичу:
```bash
git checkout develop
git pull --rebase origin develop
git checkout -b feature/123-something
git push -u origin feature/123-something
```

Обновить ветку перед PR:
```bash
git fetch origin
git rebase origin/develop
git push --force-with-lease
```

---

## 9) Что запрещено

- Прямой push в `main` и `develop`.
- Merge без ревью (кроме явно согласованных экстренных случаев).
- Коммиты без понятного смысла (“fix”, “update”, “wip” без контекста).
- Долгоживущие feature-ветки неделями.

---

## 10) Если правила нужно менять

Меняем только через PR в этот README (чтобы было зафиксировано, почему и когда поменяли процесс).
