# Frontend Redesign — Design Spec

**Date:** 2026-04-14
**Scope:** Рефакторинг существующих компонентов. Новые экраны не затрагиваются.

---

## Цели

1. **Визуал (А):** Устранить разрыв между тёмным сайдбаром и светлым контентом. Единая визуальная система для обеих тем.
2. **Код (В):** Заменить хардкодные цвета на CSS-переменные. Подключить shadcn-svelte + Tailwind как основу для всех будущих компонентов.

---

## Решения

| Вопрос | Решение |
|--------|---------|
| Тема | Светлая по умолчанию + тёмная через переключатель |
| CSS-подход | Shadcn-svelte + Tailwind + CSS custom properties |
| Сайдбар в светлой теме | Следует теме (светлый в light, тёмный в dark) |
| Палитра | Warm Zinc + Indigo |
| Скоуп | Только рефакторинг существующих компонентов |

---

## Дизайн-система

### CSS-переменные

```css
/* Светлая тема (по умолчанию) */
:root {
  --background:    #ebebea;   /* zinc-200 — фон приложения */
  --surface:       #fafafa;   /* zinc-50  — карточки, сайдбар */
  --border:        #d4d4d8;   /* zinc-300 */
  --text-primary:  #18181b;   /* zinc-900 */
  --text-muted:    #71717a;   /* zinc-500 */
  --text-subtle:   #a1a1aa;   /* zinc-400 */
  --primary:       #6366f1;   /* indigo-500 */
  --primary-hover: #4f46e5;   /* indigo-600 */
  --primary-bg:    #eef2ff;   /* indigo-50 */
  --primary-text:  #4338ca;   /* indigo-700 */
  --destructive:   #ef4444;
}

/* Тёмная тема */
.dark {
  --background:    #18181b;                  /* zinc-900 */
  --surface:       #27272a;                  /* zinc-800 */
  --border:        rgba(255, 255, 255, 0.08);
  --text-primary:  #fafafa;                  /* zinc-50 */
  --text-muted:    #a1a1aa;                  /* zinc-400 */
  --text-subtle:   #52525b;                  /* zinc-600 */
  --primary:       #6366f1;
  --primary-hover: #818cf8;                  /* indigo-400 */
  --primary-bg:    rgba(99, 102, 241, 0.15);
  --primary-text:  #a5b4fc;                  /* indigo-300 */
  --destructive:   #f87171;
}
```

### Типографика

- **Шрифт:** Inter (системный стек — нет загрузки)
  `font-family: Inter, ui-sans-serif, system-ui, -apple-system, sans-serif`
- **Body:** 13px / line-height 1.5
- **Labels / секции:** 11px / uppercase / font-weight 600 / letter-spacing 0.06em
- **Заголовки экранов:** 14–16px / font-weight 600–700

### Иконки

Заменить все текстовые псевдоиконки (`◆`, `›`, `▸`) на **Lucide Svelte**.
Конкретные иконки: `ChevronRight`, `ChevronDown`, `Search`, `Plus`, `LogOut`, `Sun`, `Moon`, `Monitor`.

### Переключение темы

- Класс `.dark` на теге `<html>`
- Три режима: `light` / `dark` / `system` (system следует `prefers-color-scheme`)
- Состояние сохраняется в `localStorage` под ключом `pulse_theme`
- Реализуется в `ProfilePopup` (переключатель уже есть в UI-спеке, но не реализован)

---

## Компоненты — что меняется

### 1. `app.css` (новый файл)

Создать `frontend/src/app.css`:
- CSS-переменные (все токены выше)
- Базовый reset: `box-sizing: border-box`, `margin: 0`
- `body { font-family: ...; background: var(--background); color: var(--text-primary); }`
- Импортировать в `+layout.svelte`

### 2. `+layout.svelte`

- Подключить `app.css` и Tailwind
- `background: #f9fafb` → `background: var(--background)`
- Добавить логику темы: читать из localStorage, слушать `prefers-color-scheme`, ставить/снимать класс `.dark` на `<html>`

### 3. `NavigationSidebar.svelte`

Наибольший объём изменений:
- Все хардкодные цвета → CSS-переменные
- `background: #1e1e2e` → `var(--surface)`
- `color: #cdd6f4` → `var(--text-primary)`
- Секции `rgba(205,214,244,0.5)` → `var(--text-subtle)`
- Активный элемент: `rgba(99,102,241,0.15)` → `var(--primary-bg)`, цвет → `var(--primary-text)`
- `◆` → `<Zap />` или `<Activity />` из Lucide (подобрать подходящую)
- `›` в chevron → `<ChevronRight />` / `<ChevronDown />`
- Добавить строку поиска (визуальная заглушка — кнопка без функционала, согласно UI-спеке поиск открывает Search Modal)
- `border: rgba(255,255,255,0.06)` → `var(--border)`

### 4. `ProfilePopup.svelte`

- Все цвета → CSS-переменные
- Реализовать переключатель темы: три кнопки (☀ / ☾ / ⬛), меняют класс `.dark` на `<html>` и сохраняют в localStorage

### 5. `ErrorToast.svelte`

- Цвета → CSS-переменные
- Фон ошибки: `var(--destructive)` + белый текст

### 6. `ConfirmationPopup.svelte`

- Цвета → CSS-переменные
- Деструктивная кнопка: `background: var(--destructive)`
- Использовать shadcn-svelte `Button` (variant `destructive`) и `Dialog`

### 7. `UnsavedChangesModal.svelte`

- Цвета → CSS-переменные
- Использовать shadcn-svelte `Dialog`

### 8. `ZeroState.svelte`

- Цвета → CSS-переменные
- `opacity`-based цвета → явные `var(--text-subtle)`

### 9. `auth/+page.svelte`

- Цвета → CSS-переменные
- Карточка формы: `background: var(--surface)`, `border: 1px solid var(--border)`
- Инпут и кнопка → shadcn-svelte `Input` и `Button`

### 10. `+error.svelte`

- Цвета → CSS-переменные

---

## Установка зависимостей

```bash
# Tailwind CSS
npm install -D tailwindcss @tailwindcss/vite --registry https://registry.npmjs.org/

# Shadcn-svelte (CLI-установщик)
npx shadcn-svelte@latest init --registry https://registry.npmjs.org/

# Lucide иконки
npm install lucide-svelte --registry https://registry.npmjs.org/

# Нужные shadcn-компоненты
npx shadcn-svelte@latest add button input dialog --registry https://registry.npmjs.org/
```

---

## Что не трогаем

- `about/+page.svelte` — минимальный контент, не является источником визуальных проблем
- `docs/+page.svelte` — таб с Swagger UI, стилизация не нужна
- Бизнес-логика во всех компонентах — только стили

---

## Чеклист перед сдачей (ui-ux-pro-max)

- [ ] Нет emoji как иконок — использовать SVG (Lucide)
- [ ] `cursor-pointer` на всех кликабельных элементах
- [ ] Hover-состояния с плавными переходами 150–200ms
- [ ] Светлая тема: контраст текста минимум 4.5:1
- [ ] Видимые focus-состояния для клавиатурной навигации
- [ ] `prefers-reduced-motion` учитывается в анимациях
- [ ] Обе темы проверены вручную
- [ ] Нет хардкодных цветов — только CSS-переменные
