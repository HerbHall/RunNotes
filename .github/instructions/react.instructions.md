---
applyTo: "ui/**/*.{ts,tsx}"
---

# React/TypeScript Coding Instructions

## Component Patterns

Define props interfaces above the component, use functional components:

```tsx
interface ItemCardProps {
    id: string
    label: string
    onSelect: (id: string) => void
}

export function ItemCard({ id, label, onSelect }: ItemCardProps) {
    return <button onClick={() => onSelect(id)}>{label}</button>
}
```

## State Management

- **useState** for client-only UI state
- **@docker/extension-api-client** for Docker Desktop SDK calls
- **No `useEffect` for data syncing** -- use nullable local override instead:

```tsx
// GOOD: nullable override, no useEffect sync
const [localOverride, setLocalOverride] = useState<string | null>(null)
const displayValue = localOverride ?? serverValue ?? ''

// Reset on save to re-sync from server
const handleSave = () => { setLocalOverride(null); refetch() }
```

## Docker Desktop Extension SDK

Always use the typed `ddClient` from `@docker/extension-api-client`:

```tsx
import { createDockerDesktopClient } from '@docker/extension-api-client'

const ddClient = createDockerDesktopClient()

// VM service calls
const result = await ddClient.extension.vm?.service?.get('/api/notes')

// Docker CLI execution
const output = await ddClient.docker.cli.exec('ps', ['--format', 'json'])
```

## Union Return Types

When a function returns a union type, add a type guard and use it at EVERY call site:

```tsx
// BAD: TS2339 -- access_token not on union
const result = await loginApi(user, pass)
setToken(result.access_token)

// GOOD: narrow first
const result = await loginApi(user, pass)
if (isMFAChallenge(result)) throw new Error('unexpected MFA')
setToken(result.access_token)
```

## TypeScript

- Strict mode -- no `any`; use `unknown` with type guards
- JSX short-circuit: `{expanded && item.details != null && <div/>}` (use `!= null`, not bare `&&` on `unknown`)
- Unused imports: ESLint catches these even when `tsc` does not -- verify every named import is used

## React Compiler Lint

Do not mutate `ref.current` during render -- wrap in `useEffect`:

```tsx
// BAD: ref mutation during render
onMessageRef.current = onMessage

// GOOD: wrap in effect
useEffect(() => { onMessageRef.current = onMessage }, [onMessage])
```

For Popper/Popover anchor elements, use callback ref with `useState` instead of `useRef`:

```tsx
// GOOD: callback ref avoids reading ref.current during render
const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null)
<Button ref={setAnchorEl}>Menu</Button>
<Popper anchorEl={anchorEl} open={open}>...</Popper>
```

## MUI v5 Patterns

This project uses MUI v5 via `@docker/docker-mui-theme`. Do NOT use MUI v6 patterns:

```tsx
// GOOD (MUI v5): InputProps for TextField adornments
<TextField InputProps={{ startAdornment: <SearchIcon /> }} />

// BAD (MUI v6): slotProps -- NOT supported in this project
<TextField slotProps={{ input: { startAdornment: <SearchIcon /> } }} />
```

Wrap disabled buttons in `<span>` when used inside `<Tooltip>`:

```tsx
<Tooltip title="Refresh">
  <span>
    <IconButton disabled={loading}><RefreshIcon /></IconButton>
  </span>
</Tooltip>
```

## Testing

Vitest + Testing Library. Use `mergeConfig(viteConfig, defineConfig(...))` so Vite
plugins and defines are inherited:

```ts
// vitest.config.ts
import { mergeConfig, defineConfig } from 'vitest/config'
import viteConfig from './vite.config'

export default mergeConfig(viteConfig, defineConfig({
    test: { environment: 'jsdom', globals: true, setupFiles: './src/test-setup.ts' }
}))
```
