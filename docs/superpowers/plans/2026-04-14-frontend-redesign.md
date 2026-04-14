# Frontend Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Рефакторинг существующих компонентов — заменить хардкодные цвета на CSS-переменные, добавить Tailwind + shadcn-svelte, реализовать переключение светлой/тёмной темы (Zinc + Indigo).

**Architecture:** Глобальные CSS-переменные в `app.css` задают все токены обеих тем. Тема управляется через `stores/theme.ts` — сохраняется в localStorage, класс `.dark` ставится на `<html>`. Компоненты используют только переменные, хардкодных цветов не остаётся.

**Tech Stack:** SvelteKit 4, shadcn-svelte@0 (Svelte 4 compatible), Tailwind CSS v3, Lucide Svelte

---

### Task 1: Установить shadcn-svelte (Tailwind + компоненты)

> shadcn-svelte@0 совместима с Svelte 4. Она автоматически устанавливает и настраивает Tailwind v3.

**Files:**
- Modify: `frontend/vite.config.ts`
- Modify: `frontend/svelte.config.js`
- Modify: `frontend/package.json`
- Create: `frontend/tailwind.config.js`
- Create: `frontend/postcss.config.js`
- Modify: `frontend/src/app.html` (добавит shadcn init, если нужно)

- [ ] **Step 1: Запустить shadcn-svelte init**

```bash
cd frontend
npx --registry https://registry.npmjs.org/ shadcn-svelte@0 init
```

На вопросы отвечать:
- Style: `Default`
- Base color: `Zinc`
- CSS variables: `Yes`
- Global CSS file: `src/app.css`
- Tailwind config: `tailwind.config.js`
- Import alias: `$lib` (или оставить предложенное)

- [ ] **Step 2: Добавить нужные компоненты**

```bash
cd frontend
npx --registry https://registry.npmjs.org/ shadcn-svelte@0 add button input dialog
```

- [ ] **Step 3: Проверить, что TypeScript компилируется**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 ошибок (или только pre-existing warnings — не блокируют).

- [ ] **Step 4: Проверить, что проект собирается**

```bash
cd frontend
npm run build
```

Ожидаемый результат: `build` завершился без ошибок.

- [ ] **Step 5: Установить Lucide Svelte**

```bash
cd frontend
npm install --registry https://registry.npmjs.org/ lucide-svelte
```

- [ ] **Step 6: Коммит**

```bash
git add frontend/
git commit -m "chore(frontend): add tailwind, shadcn-svelte, lucide-svelte"
```

---

### Task 2: Создать тему — `app.css` + `stores/theme.ts`

**Files:**
- Modify: `frontend/src/app.css`
- Create: `frontend/src/stores/theme.ts`

- [ ] **Step 1: Написать тест для theme store**

Создать файл `frontend/tests/theme.test.ts`:

```ts
import { describe, it, expect, beforeEach } from 'vitest';

// Smoke-тест: theme store экспортирует нужное API
// Полноценное E2E тестирование — в Task 11
describe('theme store exports', () => {
  it('exports setTheme and getTheme functions', async () => {
    const mod = await import('../src/stores/theme');
    expect(typeof mod.setTheme).toBe('function');
    expect(typeof mod.getTheme).toBe('function');
    expect(typeof mod.theme).toBe('object'); // Svelte store
  });
});
```

> Примечание: если vitest не настроен — пропустить unit-тест и перейти сразу к реализации. Тема будет проверена через E2E в Task 11.

- [ ] **Step 2: Реализовать `stores/theme.ts`**

```ts
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export type ThemeMode = 'light' | 'dark' | 'system';

const STORAGE_KEY = 'pulse_theme';

function resolveTheme(mode: ThemeMode): 'light' | 'dark' {
  if (mode === 'system') {
    return browser && window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light';
  }
  return mode;
}

function applyTheme(mode: ThemeMode) {
  if (!browser) return;
  const resolved = resolveTheme(mode);
  document.documentElement.classList.toggle('dark', resolved === 'dark');
}

export function getTheme(): ThemeMode {
  if (!browser) return 'system';
  return (localStorage.getItem(STORAGE_KEY) as ThemeMode) ?? 'system';
}

export function setTheme(mode: ThemeMode) {
  if (!browser) return;
  localStorage.setItem(STORAGE_KEY, mode);
  applyTheme(mode);
  theme.set(mode);
}

export const theme = writable<ThemeMode>(getTheme());

export function initTheme() {
  const mode = getTheme();
  applyTheme(mode);
  if (browser) {
    window
      .matchMedia('(prefers-color-scheme: dark)')
      .addEventListener('change', () => {
        if (getTheme() === 'system') applyTheme('system');
      });
  }
}
```

- [ ] **Step 3: Дополнить `app.css` CSS-переменными**

Открыть `frontend/src/app.css`. Найти секцию `:root` (shadcn-svelte init её создаёт) и заменить целиком на:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    /* Zinc + Indigo — светлая тема */
    --background:    240 4% 92%;    /* #ebebea */
    --foreground:    240 10% 10%;   /* #18181b */
    --card:          0 0% 98%;      /* #fafafa */
    --card-foreground: 240 10% 10%;
    --popover:       0 0% 100%;
    --popover-foreground: 240 10% 10%;
    --primary:       239 84% 67%;   /* #6366f1 */
    --primary-foreground: 0 0% 100%;
    --secondary:     240 5% 88%;    /* zinc-300 */
    --secondary-foreground: 240 10% 10%;
    --muted:         240 5% 88%;
    --muted-foreground: 240 4% 46%; /* #71717a */
    --accent:        239 84% 95%;   /* indigo-50 */
    --accent-foreground: 239 84% 40%; /* indigo-700 */
    --destructive:   0 84% 60%;     /* #ef4444 */
    --destructive-foreground: 0 0% 100%;
    --border:        240 6% 83%;    /* #d4d4d8 */
    --input:         240 6% 83%;
    --ring:          239 84% 67%;
    --radius:        0.5rem;
  }

  .dark {
    /* Zinc — тёмная тема */
    --background:    240 10% 11%;   /* #18181b */
    --foreground:    0 0% 98%;      /* #fafafa */
    --card:          240 6% 16%;    /* #27272a */
    --card-foreground: 0 0% 98%;
    --popover:       240 6% 16%;
    --popover-foreground: 0 0% 98%;
    --primary:       239 84% 67%;   /* #6366f1 */
    --primary-foreground: 0 0% 100%;
    --secondary:     240 5% 20%;
    --secondary-foreground: 0 0% 98%;
    --muted:         240 6% 16%;
    --muted-foreground: 240 5% 65%; /* #a1a1aa */
    --accent:        239 84% 20%;
    --accent-foreground: 239 84% 80%;
    --destructive:   0 91% 71%;     /* #f87171 */
    --destructive-foreground: 0 0% 100%;
    --border:        240 6% 22%;
    --input:         240 6% 22%;
    --ring:          239 84% 67%;
  }
}

@layer base {
  * {
    @apply border-border;
  }

  body {
    @apply bg-background text-foreground;
    font-family: Inter, ui-sans-serif, system-ui, -apple-system, sans-serif;
    font-size: 13px;
    line-height: 1.5;
  }
}
```

- [ ] **Step 4: Проверить TypeScript**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 новых ошибок.

- [ ] **Step 5: Коммит**

```bash
git add frontend/src/app.css frontend/src/stores/theme.ts
git commit -m "feat(frontend): add theme store and CSS design tokens"
```

---

### Task 3: Обновить `+layout.svelte`

**Files:**
- Modify: `frontend/src/routes/+layout.svelte`

- [ ] **Step 1: Заменить файл**

```svelte
<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { getActiveAccount } from '../stores/session';
  import { initTheme } from '../stores/theme';
  import ErrorToast from '../components/ui/ErrorToast.svelte';
  import NavigationSidebar from '../components/ui/NavigationSidebar.svelte';

  const PUBLIC_ROUTES = ['/auth', '/auth/verify'];

  $: isPublic = PUBLIC_ROUTES.some((route) => $page.url.pathname.startsWith(route));

  onMount(() => {
    initTheme();
    if (!isPublic && !getActiveAccount()) {
      goto('/auth');
    }
  });
</script>

<ErrorToast />

{#if isPublic}
  <slot />
{:else}
  <div class="app-shell">
    <NavigationSidebar currentPath={$page.url.pathname} />
    <main class="main-content">
      <slot />
    </main>
  </div>
{/if}

<style>
  .app-shell {
    display: flex;
    height: 100vh;
    overflow: hidden;
  }

  .main-content {
    flex: 1;
    overflow: auto;
    background: hsl(var(--background));
  }
</style>
```

- [ ] **Step 2: Проверить TypeScript**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 новых ошибок.

- [ ] **Step 3: Коммит**

```bash
git add frontend/src/routes/+layout.svelte
git commit -m "feat(frontend): init theme on mount, import app.css"
```

---

### Task 4: Обновить `NavigationSidebar.svelte`

**Files:**
- Modify: `frontend/src/components/ui/NavigationSidebar.svelte`

- [ ] **Step 1: Заменить файл**

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { getActiveAccount, session } from '../../stores/session';
  import { Search, ChevronRight, ChevronDown, Plus, Activity } from 'lucide-svelte';
  import ProfilePopup from './ProfilePopup.svelte';
  import ZeroState from './ZeroState.svelte';

  export let currentPath: string = '/';

  let profileOpen = false;
  let account = getActiveAccount();
  let filtersOpen = true;
  let teamsOpen = true;

  const STORAGE_KEY = 'pulse_sidebar_state';

  onMount(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const state = JSON.parse(raw);
        filtersOpen = state.filtersOpen ?? true;
        teamsOpen = state.teamsOpen ?? true;
      }
    } catch {
      /* ignore */
    }
  });

  session.subscribe(() => {
    account = getActiveAccount();
  });

  function saveSidebarState() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({ filtersOpen, teamsOpen }));
    } catch {
      /* ignore */
    }
  }

  function toggleFilters() {
    filtersOpen = !filtersOpen;
    saveSidebarState();
  }

  function toggleTeams() {
    teamsOpen = !teamsOpen;
    saveSidebarState();
  }
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <div class="logo-icon">
      <Activity size={16} />
    </div>
    <span class="logo-text">Pulse</span>
  </div>

  <div class="search-row">
    <button type="button" class="search-btn" aria-label="Поиск">
      <Search size={14} />
      <span>Поиск...</span>
    </button>
  </div>

  <nav class="sidebar-nav">
    <button
      type="button"
      class="nav-item"
      class:active={currentPath === '/'}
      on:click={() => goto('/')}
    >
      Personal Activity
    </button>

    <div class="section">
      <button type="button" class="section-header" on:click={toggleFilters}>
        {#if filtersOpen}
          <ChevronDown size={12} />
        {:else}
          <ChevronRight size={12} />
        {/if}
        Личные фильтры
      </button>
      {#if filtersOpen}
        <div class="section-body">
          <button type="button" class="action-btn">
            <Plus size={12} />
            Создать фильтр
          </button>
          <ZeroState message="Нет фильтров" />
          <button type="button" class="action-btn">Импортировать задачи</button>
        </div>
      {/if}
    </div>

    <div class="section">
      <button type="button" class="section-header" on:click={toggleTeams}>
        {#if teamsOpen}
          <ChevronDown size={12} />
        {:else}
          <ChevronRight size={12} />
        {/if}
        Команды
      </button>
      {#if teamsOpen}
        <div class="section-body">
          <button type="button" class="action-btn">
            <Plus size={12} />
            Создать команду
          </button>
          <ZeroState message="Нет команд" />
        </div>
      {/if}
    </div>

    <button type="button" class="nav-item">Автоматизации</button>
    <button type="button" class="nav-item">Метки</button>
    <button
      type="button"
      class="nav-item"
      class:active={currentPath.startsWith('/about')}
      on:click={() => goto('/about')}
    >
      О системе
    </button>
    <button
      type="button"
      class="nav-item"
      class:active={currentPath.startsWith('/docs')}
      on:click={() => goto('/docs')}
    >
      Документация
    </button>
  </nav>

  <div class="sidebar-footer">
    {#if account}
      <button
        type="button"
        class="profile-btn"
        on:click={() => (profileOpen = !profileOpen)}
        aria-label="Профиль"
      >
        <span class="avatar">{account.username[0]?.toUpperCase()}</span>
        <span class="username">{account.username}</span>
      </button>
    {/if}
    <ProfilePopup bind:open={profileOpen} />
  </div>
</aside>

<style>
  .sidebar {
    width: 240px;
    background: hsl(var(--card));
    color: hsl(var(--foreground));
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    position: relative;
    height: 100vh;
    border-right: 1px solid hsl(var(--border));
  }

  .sidebar-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 16px;
    border-bottom: 1px solid hsl(var(--border));
    flex-shrink: 0;
  }

  .logo-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    background: hsl(var(--primary));
    color: hsl(var(--primary-foreground));
    border-radius: 6px;
  }

  .logo-text {
    font-size: 15px;
    font-weight: 700;
    color: hsl(var(--foreground));
  }

  .search-row {
    padding: 8px 10px;
    flex-shrink: 0;
  }

  .search-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    width: 100%;
    padding: 6px 8px;
    background: hsl(var(--muted));
    border: 1px solid hsl(var(--border));
    border-radius: 6px;
    cursor: pointer;
    font-size: 12px;
    color: hsl(var(--muted-foreground));
    transition: background 0.15s;
    text-align: left;
  }

  .search-btn:hover {
    background: hsl(var(--secondary));
  }

  .sidebar-nav {
    flex: 1;
    overflow-y: auto;
    padding: 4px 0;
  }

  .nav-item {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 7px 12px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 13px;
    color: hsl(var(--foreground));
    text-align: left;
    transition: background 0.15s;
    border-radius: 0;
  }

  .nav-item:hover {
    background: hsl(var(--muted));
  }

  .nav-item.active {
    background: hsl(var(--accent));
    color: hsl(var(--accent-foreground));
    font-weight: 500;
  }

  .section {
    margin: 2px 0;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 4px;
    width: 100%;
    padding: 6px 12px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 11px;
    font-weight: 600;
    color: hsl(var(--muted-foreground));
    text-transform: uppercase;
    letter-spacing: 0.05em;
    text-align: left;
    transition: background 0.15s;
  }

  .section-header:hover {
    background: hsl(var(--muted));
  }

  .section-body {
    padding: 0 0 4px 16px;
  }

  .action-btn {
    display: flex;
    align-items: center;
    gap: 5px;
    width: 100%;
    padding: 5px 12px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 12px;
    color: hsl(var(--muted-foreground));
    text-align: left;
    transition: color 0.15s;
  }

  .action-btn:hover {
    color: hsl(var(--foreground));
  }

  .sidebar-footer {
    margin-top: auto;
    padding: 10px 8px;
    border-top: 1px solid hsl(var(--border));
    flex-shrink: 0;
  }

  .profile-btn {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
    border-radius: 6px;
    color: hsl(var(--foreground));
    transition: background 0.15s;
    text-align: left;
  }

  .profile-btn:hover {
    background: hsl(var(--muted));
  }

  .avatar {
    width: 28px;
    height: 28px;
    background: hsl(var(--primary));
    color: hsl(var(--primary-foreground));
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    font-weight: 600;
    flex-shrink: 0;
  }

  .username {
    font-size: 13px;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: hsl(var(--foreground));
  }
</style>
```

- [ ] **Step 2: Проверить TypeScript**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 новых ошибок.

- [ ] **Step 3: Коммит**

```bash
git add frontend/src/components/ui/NavigationSidebar.svelte
git commit -m "feat(frontend): refactor NavigationSidebar to CSS variables and Lucide icons"
```

---

### Task 5: Обновить `ProfilePopup.svelte` + переключатель темы

**Files:**
- Modify: `frontend/src/components/ui/ProfilePopup.svelte`

- [ ] **Step 1: Заменить файл**

```svelte
<script lang="ts">
  import { goto } from '$app/navigation';
  import { Sun, Moon, Monitor, LogOut, UserPlus } from 'lucide-svelte';
  import { session, getActiveAccount, removeAccount, switchAccount } from '../../stores/session';
  import { auth } from '../../api/auth';
  import { theme, setTheme } from '../../stores/theme';
  import type { ThemeMode } from '../../stores/theme';

  export let open = false;

  let loggingOut = false;

  async function handleLogout() {
    const account = getActiveAccount();
    if (!account) return;

    loggingOut = true;
    try {
      await auth.logout(account.refreshToken);
    } catch {
      // продолжаем logout даже при сетевой ошибке
    } finally {
      removeAccount(account.id);
      loggingOut = false;
      open = false;
      const remaining = getActiveAccount();
      goto(remaining ? '/' : '/auth');
    }
  }

  function handleSwitch(id: string) {
    switchAccount(id);
    open = false;
    goto('/');
  }

  function handleAddUser() {
    open = false;
    goto('/auth');
  }

  function handleBackdropClick() {
    open = false;
  }

  const THEME_OPTIONS: { mode: ThemeMode; icon: typeof Sun; label: string }[] = [
    { mode: 'light', icon: Sun, label: 'Светлая' },
    { mode: 'dark', icon: Moon, label: 'Тёмная' },
    { mode: 'system', icon: Monitor, label: 'Системная' },
  ];
</script>

{#if open}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="backdrop" on:click={handleBackdropClick} />

  <div class="popup" role="dialog" aria-label="Профиль">
    <ul class="accounts-list">
      {#each $session.accounts as account, index (account.id)}
        <li>
          <button
            type="button"
            class="account-item"
            class:active={index === $session.activeIndex}
            on:click={() => handleSwitch(account.id)}
          >
            <span class="avatar">{account.username[0]?.toUpperCase()}</span>
            <span class="account-info">
              <span class="username">{account.username}</span>
              <span class="email">{account.email}</span>
            </span>
          </button>
        </li>
      {/each}
    </ul>

    <div class="theme-section">
      <span class="theme-label">Тема</span>
      <div class="theme-buttons">
        {#each THEME_OPTIONS as opt}
          <button
            type="button"
            class="theme-btn"
            class:active={$theme === opt.mode}
            on:click={() => setTheme(opt.mode)}
            aria-label={opt.label}
            title={opt.label}
          >
            <svelte:component this={opt.icon} size={14} />
          </button>
        {/each}
      </div>
    </div>

    <div class="actions">
      <button type="button" class="action-btn" on:click={handleAddUser}>
        <UserPlus size={14} />
        Добавить пользователя
      </button>
      <button
        type="button"
        class="action-btn logout-btn"
        on:click={handleLogout}
        disabled={loggingOut}
      >
        {#if loggingOut}
          <span class="loader" />
        {:else}
          <LogOut size={14} />
          Выйти
        {/if}
      </button>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    z-index: 99;
  }

  .popup {
    position: fixed;
    bottom: 64px;
    left: 16px;
    z-index: 100;
    background: hsl(var(--popover));
    color: hsl(var(--popover-foreground));
    border: 1px solid hsl(var(--border));
    border-radius: 10px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
    min-width: 240px;
    overflow: hidden;
  }

  .accounts-list {
    list-style: none;
    margin: 0;
    padding: 8px 0;
    border-bottom: 1px solid hsl(var(--border));
  }

  .account-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 10px 16px;
    background: none;
    border: none;
    cursor: pointer;
    text-align: left;
    transition: background 0.15s;
  }

  .account-item:hover {
    background: hsl(var(--muted));
  }

  .account-item.active {
    background: hsl(var(--accent));
  }

  .avatar {
    width: 32px;
    height: 32px;
    background: hsl(var(--primary));
    color: hsl(var(--primary-foreground));
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 13px;
    font-weight: 600;
    flex-shrink: 0;
  }

  .account-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
    overflow: hidden;
  }

  .username {
    font-size: 13px;
    font-weight: 500;
    color: hsl(var(--foreground));
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .email {
    font-size: 12px;
    color: hsl(var(--muted-foreground));
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .theme-section {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 16px;
    border-bottom: 1px solid hsl(var(--border));
  }

  .theme-label {
    font-size: 12px;
    color: hsl(var(--muted-foreground));
  }

  .theme-buttons {
    display: flex;
    gap: 2px;
    background: hsl(var(--muted));
    border-radius: 6px;
    padding: 2px;
  }

  .theme-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: none;
    border: none;
    cursor: pointer;
    border-radius: 4px;
    color: hsl(var(--muted-foreground));
    transition: background 0.15s, color 0.15s;
  }

  .theme-btn:hover {
    background: hsl(var(--secondary));
    color: hsl(var(--foreground));
  }

  .theme-btn.active {
    background: hsl(var(--card));
    color: hsl(var(--foreground));
  }

  .actions {
    padding: 8px 0;
    display: flex;
    flex-direction: column;
  }

  .action-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 10px 16px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 13px;
    color: hsl(var(--foreground));
    text-align: left;
    transition: background 0.15s;
  }

  .action-btn:hover:not(:disabled) {
    background: hsl(var(--muted));
  }

  .logout-btn {
    color: hsl(var(--destructive));
  }

  .logout-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .loader {
    width: 14px;
    height: 14px;
    border: 2px solid hsl(var(--destructive) / 0.3);
    border-top-color: hsl(var(--destructive));
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
    display: inline-block;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
```

- [ ] **Step 2: Проверить TypeScript**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 новых ошибок.

- [ ] **Step 3: Коммит**

```bash
git add frontend/src/components/ui/ProfilePopup.svelte
git commit -m "feat(frontend): refactor ProfilePopup, add theme switcher"
```

---

### Task 6: Обновить малые компоненты (`ZeroState`, `ErrorToast`, `+error.svelte`)

**Files:**
- Modify: `frontend/src/components/ui/ZeroState.svelte`
- Modify: `frontend/src/components/ui/ErrorToast.svelte`
- Modify: `frontend/src/routes/+error.svelte`

- [ ] **Step 1: Заменить `ZeroState.svelte`**

```svelte
<script lang="ts">
  export let message: string = 'Нет данных';
</script>

<div class="zero-state">
  <p class="message">{message}</p>
</div>

<style>
  .zero-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 16px 0;
  }

  .message {
    font-size: 12px;
    color: hsl(var(--muted-foreground));
    margin: 0;
  }
</style>
```

- [ ] **Step 2: Заменить `ErrorToast.svelte`**

```svelte
<script lang="ts">
  import { X } from 'lucide-svelte';
  import { toastError } from '../../stores/toast';

  let timer: ReturnType<typeof setTimeout>;

  $: if ($toastError) {
    clearTimeout(timer);
    timer = setTimeout(() => toastError.set(null), 4000);
  }

  function dismiss() {
    clearTimeout(timer);
    toastError.set(null);
  }
</script>

{#if $toastError}
  <div class="error-toast" role="alert">
    <span>{$toastError}</span>
    <button type="button" aria-label="Закрыть" on:click={dismiss}>
      <X size={14} />
    </button>
  </div>
{/if}

<style>
  .error-toast {
    position: fixed;
    bottom: 24px;
    right: 24px;
    z-index: 9999;
    display: flex;
    align-items: center;
    gap: 12px;
    background: hsl(var(--destructive));
    color: hsl(var(--destructive-foreground));
    padding: 12px 16px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    cursor: pointer;
    max-width: 360px;
    font-size: 14px;
  }

  button {
    display: flex;
    align-items: center;
    background: none;
    border: none;
    color: hsl(var(--destructive-foreground));
    cursor: pointer;
    padding: 0;
    flex-shrink: 0;
  }
</style>
```

- [ ] **Step 3: Заменить `+error.svelte`**

```svelte
<script lang="ts">
  import { page } from '$app/stores';
</script>

<div class="error-page">
  {#if $page.status === 403}
    <h1>Нет доступа</h1>
    <p>У вас нет прав для просмотра этой страницы.</p>
  {:else if $page.status === 404}
    <h1>Не найдено</h1>
    <p>Запрашиваемая страница не существует или была удалена.</p>
  {:else}
    <h1>Произошла ошибка</h1>
    <p>{$page.error?.message ?? 'Что-то пошло не так. Попробуйте обновить страницу.'}</p>
  {/if}
  <p class="status-code">Код ошибки: {$page.status}</p>
  <a href="/" class="home-link">На главную</a>
</div>

<style>
  .error-page {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: hsl(var(--background));
    text-align: center;
    padding: 32px;
  }

  h1 {
    font-size: 28px;
    font-weight: 700;
    color: hsl(var(--foreground));
    margin: 0 0 12px;
  }

  p {
    font-size: 14px;
    color: hsl(var(--muted-foreground));
    margin: 0 0 8px;
  }

  .status-code {
    font-size: 12px;
    color: hsl(var(--muted-foreground));
    opacity: 0.6;
  }

  .home-link {
    margin-top: 24px;
    color: hsl(var(--primary));
    font-size: 14px;
    text-decoration: none;
    font-weight: 500;
  }

  .home-link:hover {
    text-decoration: underline;
  }
</style>
```

- [ ] **Step 4: Проверить TypeScript**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 новых ошибок.

- [ ] **Step 5: Коммит**

```bash
git add frontend/src/components/ui/ZeroState.svelte \
        frontend/src/components/ui/ErrorToast.svelte \
        frontend/src/routes/+error.svelte
git commit -m "feat(frontend): refactor ZeroState, ErrorToast, error screen to CSS variables"
```

---

### Task 7: Обновить модалки (`ConfirmationPopup`, `UnsavedChangesModal`)

**Files:**
- Modify: `frontend/src/components/ui/ConfirmationPopup.svelte`
- Modify: `frontend/src/components/ui/UnsavedChangesModal.svelte`

> Используем shadcn-svelte `Button`. Dialog-примитив сохраняем самописным — он уже корректно управляет z-index и backdrop.

- [ ] **Step 1: Заменить `ConfirmationPopup.svelte`**

```svelte
<script lang="ts">
  import { Button } from '$lib/components/ui/button';

  export let open: boolean = false;
  export let operation: string;
  export let onConfirm: () => void;
  export let onCancel: () => void;
</script>

{#if open}
  <div class="backdrop" on:click={onCancel} role="presentation" />
  <div class="popup" role="dialog" aria-modal="true">
    <p class="message">Вы подтверждаете {operation}?</p>
    <div class="actions">
      <Button variant="outline" on:click={onCancel}>Отмена</Button>
      <Button variant="destructive" on:click={onConfirm}>Да, подтверждаю</Button>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 99;
  }

  .popup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 100;
    background: hsl(var(--card));
    color: hsl(var(--card-foreground));
    border: 1px solid hsl(var(--border));
    border-radius: 10px;
    padding: 24px;
    min-width: 320px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.16);
  }

  .message {
    font-size: 14px;
    color: hsl(var(--foreground));
    margin: 0 0 20px;
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
</style>
```

- [ ] **Step 2: Заменить `UnsavedChangesModal.svelte`**

```svelte
<script lang="ts">
  import { Button } from '$lib/components/ui/button';

  export let open: boolean = false;
  export let onConfirm: () => void;
  export let onCancel: () => void;
</script>

{#if open}
  <div class="backdrop" on:click={onCancel} role="presentation" />
  <div class="popup" role="dialog" aria-modal="true">
    <p class="message">Вы уверены? Несохранённые данные будут потеряны.</p>
    <div class="actions">
      <Button variant="outline" on:click={onCancel}>Отмена</Button>
      <Button variant="destructive" on:click={onConfirm}>Да, закрыть</Button>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 99;
  }

  .popup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 100;
    background: hsl(var(--card));
    color: hsl(var(--card-foreground));
    border: 1px solid hsl(var(--border));
    border-radius: 10px;
    padding: 24px;
    min-width: 320px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.16);
  }

  .message {
    font-size: 14px;
    color: hsl(var(--foreground));
    margin: 0 0 20px;
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
</style>
```

- [ ] **Step 3: Проверить TypeScript**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 новых ошибок.

- [ ] **Step 4: Коммит**

```bash
git add frontend/src/components/ui/ConfirmationPopup.svelte \
        frontend/src/components/ui/UnsavedChangesModal.svelte
git commit -m "feat(frontend): refactor modals to CSS variables and shadcn Button"
```

---

### Task 8: Обновить экран авторизации

**Files:**
- Modify: `frontend/src/routes/auth/+page.svelte`

- [ ] **Step 1: Заменить файл**

```svelte
<script lang="ts">
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { auth } from '../../api/auth';

  let email = '';
  let emailError = '';
  let loading = false;
  let sent = false;
  let touched = false;

  function validateEmail(value: string): string {
    if (!value.trim()) return 'Введите email';
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return 'Введите корректный email';
    return '';
  }

  function handleBlur() {
    touched = true;
    emailError = validateEmail(email);
  }

  async function handleSubmit() {
    touched = true;
    emailError = validateEmail(email);
    if (emailError) return;

    loading = true;
    try {
      await auth.sendMagicLink(email);
      sent = true;
    } catch (error) {
      emailError = error instanceof Error ? error.message : 'Ошибка отправки';
    } finally {
      loading = false;
    }
  }

  async function handleResend() {
    sent = false;
    loading = true;
    try {
      await auth.sendMagicLink(email);
      sent = true;
    } catch (error) {
      emailError = error instanceof Error ? error.message : 'Ошибка отправки';
      sent = false;
    } finally {
      loading = false;
    }
  }
</script>

<div class="auth-screen">
  <div class="auth-card">
    <h1 class="auth-title">Pulse</h1>

    {#if sent}
      <div class="sent-state">
        <p class="sent-message">
          Проверьте почту — ссылка отправлена на <strong>{email}</strong>
        </p>
        <Button variant="outline" class="w-full" on:click={handleResend} disabled={loading}>
          {#if loading}
            <span class="loader" />
          {:else}
            Отправить снова
          {/if}
        </Button>
      </div>
    {:else}
      <form on:submit|preventDefault={handleSubmit} novalidate>
        <div class="field" class:has-error={touched && emailError}>
          <label for="email">Email</label>
          <Input
            id="email"
            type="email"
            bind:value={email}
            on:blur={handleBlur}
            placeholder="you@example.com"
            autocomplete="email"
            disabled={loading}
            class={touched && emailError ? 'border-destructive' : ''}
          />
          {#if touched && emailError}
            <span class="field-error">{emailError}</span>
          {/if}
        </div>

        <Button type="submit" class="w-full" disabled={loading}>
          {#if loading}
            <span class="loader" />
          {:else}
            Войти
          {/if}
        </Button>
      </form>
    {/if}
  </div>
</div>

<style>
  .auth-screen {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    background: hsl(var(--background));
  }

  .auth-card {
    background: hsl(var(--card));
    color: hsl(var(--card-foreground));
    border: 1px solid hsl(var(--border));
    border-radius: 12px;
    padding: 40px;
    width: 100%;
    max-width: 400px;
    box-shadow: 0 2px 16px rgba(0, 0, 0, 0.08);
  }

  .auth-title {
    font-size: 28px;
    font-weight: 700;
    color: hsl(var(--foreground));
    margin: 0 0 32px;
    text-align: center;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 20px;
  }

  label {
    font-size: 14px;
    font-weight: 500;
    color: hsl(var(--foreground));
  }

  .field-error {
    font-size: 12px;
    color: hsl(var(--destructive));
  }

  .sent-state {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .sent-message {
    font-size: 14px;
    color: hsl(var(--foreground));
    line-height: 1.5;
    margin: 0;
  }

  .loader {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.4);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
    display: inline-block;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
```

- [ ] **Step 2: Проверить TypeScript**

```bash
cd frontend
npm run check
```

Ожидаемый результат: 0 новых ошибок.

- [ ] **Step 3: Коммит**

```bash
git add frontend/src/routes/auth/+page.svelte
git commit -m "feat(frontend): refactor auth screen to CSS variables and shadcn Input/Button"
```

---

### Task 9: E2E тест переключения темы + финальная проверка

**Files:**
- Create: `e2e/theme.spec.ts`

- [ ] **Step 1: Написать E2E тест**

Создать `e2e/theme.spec.ts`:

```ts
import { test, expect } from '@playwright/test';

test.describe('theme switching', () => {
  test('light theme is default', async ({ page }) => {
    await page.goto('/auth');
    const htmlClass = await page.locator('html').getAttribute('class');
    expect(htmlClass ?? '').not.toContain('dark');
  });

  test('dark theme persists after reload', async ({ page }) => {
    await page.goto('/auth');

    // Устанавливаем тёмную тему через localStorage напрямую
    await page.evaluate(() => {
      localStorage.setItem('pulse_theme', 'dark');
    });
    await page.reload();

    const htmlClass = await page.locator('html').getAttribute('class');
    expect(htmlClass).toContain('dark');
  });

  test('light theme persists after reload', async ({ page }) => {
    await page.goto('/auth');

    await page.evaluate(() => {
      localStorage.setItem('pulse_theme', 'light');
    });
    await page.reload();

    const htmlClass = await page.locator('html').getAttribute('class');
    expect(htmlClass ?? '').not.toContain('dark');
  });
});
```

- [ ] **Step 2: Запустить E2E тест**

```bash
cd /Users/gaevivan/projects/pulse
npx --registry https://registry.npmjs.org/ playwright test e2e/theme.spec.ts
```

Ожидаемый результат: все 3 теста проходят.

> Если приложение не запущено — сначала запустить `make up` или `make dev`, затем повторить.

- [ ] **Step 3: Финальная сборка**

```bash
cd frontend
npm run build
```

Ожидаемый результат: сборка без ошибок.

- [ ] **Step 4: Финальный коммит**

```bash
git add e2e/theme.spec.ts
git commit -m "test(e2e): add theme persistence tests"
```
