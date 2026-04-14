# Code Conventions

## Git Commits

Каждый коммит, завершающий реализацию задачи из GitHub Issues, должен содержать closing keyword с номером задачи:

```
feat: add magic link auth

closes #3
```

Поддерживаемые ключевые слова: `closes`, `fixes`, `resolves`. При мерже в main GitHub автоматически закроет указанную задачу.

Если коммит частично реализует задачу, использовать `ref #N` вместо `closes #N` — задача останется открытой.

Если один коммит закрывает несколько задач, перечислить каждую на отдельной строке:

```
feat: add task CRUD and labels API

closes #11
closes #10
```
