# Luengo Engine

**Luengo Engine** is a lightweight, modular game engine written in Go with built-in scripting support using Lua. It is designed to be simple, flexible, and powerful — ideal for game developers who want performance and modding capabilities without the complexity.

## 🚀 Features

* 🔄 Written in Go for performance and concurrency
* 🧠 Embedded Lua scripting with [Gopher-Lua](https://github.com/yuin/gopher-lua)
* 🧩 Modular architecture for flexibility and expansion
* 🎮 Custom mod loading system
* 🧱 ECS-friendly foundation (coming soon)

## 📂 Project Structure

```
luengo/
├── engine/         # Core engine modules (graphics, input, audio, etc.)
├── core/           # Game loop, scene management, ECS (planned)
├── mod/            # Lua scripts and mods
│   └── main.lua    # Entry Lua script
├── assets/         # Game assets (images, audio, etc.)
├── main.go         # Entry point of the engine
└── README.md       # You're here 😄
```

## 🧪 Running the Engine

Make sure you have [Go](https://golang.org/dl/) installed.

```bash
git clone https://github.com/yourname/luengo.git
cd luengo
go run main.go
```

## 📝 Example Lua Script (`mod/main.lua`)

```lua
log("Hello from Lua!")
```

## 🛠 Roadmap

* [ ] ECS system for entities and components
* [ ] Input and audio modules
* [ ] Sprite rendering with OpenGL or Ebiten
* [ ] Lua sandboxing for secure modding
* [ ] Debug console (in-Lua or Go-based)

### 🧠 **Luengo Engine – Dev Summary**

**Nombre del proyecto:** Luengo Engine

**Core Language:** Go

**Scripting:** Lua (Gopher-Lua)

**Rendering:** Ebiten

**Audio:** Beep

**Input Handling:** Ebiten

**Arquitectura:** Modular, event-driven, ECS-ready

**Licencia:** MIT (Open Source)

---

### ✅ **Fases completadas (1 a 8)**

1. **Entity System** – Con soporte para entidades con ID, nombre, posición, sprite.
2. **Audio Module** – Soporte para reproducir sonidos `.wav` desde Lua (`play_sound`).
3. **Input Module** – Teclado integrado (WASD, Arrows), uso desde Lua (`is_key_pressed`).
4. **Rendering** – Motor 2D con Ebiten, renderiza sprites por coordenadas.
5. **Mod Loader** – Carga automática de todos los `.lua` en `mod/` y subdirectorios.
6. **Game Loop** – Soporte para `on_start()` y `on_update()` desde Lua.
7. **Debug Tools** – Consola visual (F3), comando `debug()`, FPS & posición.
8. **Technical Documentation** – Documento `TECHNICAL.md` con todo lo necesario pa’ devs.

---

### 📜 **Lua API Actual**

<pre class="!overflow-visible" data-start="1279" data-end="1752"><div class="contain-inline-size rounded-md border-[0.5px] border-token-border-medium relative bg-token-sidebar-surface-primary"><div class="flex items-center text-token-text-secondary px-4 py-2 text-xs font-sans justify-between h-9 bg-token-sidebar-surface-primary dark:bg-token-main-surface-secondary select-none rounded-t-[5px]">lua</div><div class="sticky top-9"><div class="absolute bottom-0 right-0 flex h-9 items-center pr-2"><div class="flex items-center rounded bg-token-sidebar-surface-primary px-2 font-sans text-xs text-token-text-secondary dark:bg-token-main-surface-secondary"><span class="" data-state="closed"><button class="flex gap-1 items-center select-none px-4 py-1" aria-label="Copy"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="icon-xs"><path fill-rule="evenodd" clip-rule="evenodd" d="M7 5C7 3.34315 8.34315 2 10 2H19C20.6569 2 22 3.34315 22 5V14C22 15.6569 20.6569 17 19 17H17V19C17 20.6569 15.6569 22 14 22H5C3.34315 22 2 20.6569 2 19V10C2 8.34315 3.34315 7 5 7H7V5ZM9 7H14C15.6569 7 17 8.34315 17 10V15H19C19.5523 15 20 14.5523 20 14V5C20 4.44772 19.5523 4 19 4H10C9.44772 4 9 4.44772 9 5V7ZM5 9C4.44772 9 4 9.44772 4 10V19C4 19.5523 4.44772 20 5 20H14C14.5523 20 15 19.5523 15 19V10C15 9.44772 14.5523 9 14 9H5Z" fill="currentColor"></path></svg>Copy</button></span><span class="" data-state="closed"><button class="flex select-none items-center gap-1 px-4 py-1"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="icon-xs"><path d="M2.5 5.5C4.3 5.2 5.2 4 5.5 2.5C5.8 4 6.7 5.2 8.5 5.5C6.7 5.8 5.8 7 5.5 8.5C5.2 7 4.3 5.8 2.5 5.5Z" fill="currentColor" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"></path><path d="M5.66282 16.5231L5.18413 19.3952C5.12203 19.7678 5.09098 19.9541 5.14876 20.0888C5.19933 20.2067 5.29328 20.3007 5.41118 20.3512C5.54589 20.409 5.73218 20.378 6.10476 20.3159L8.97693 19.8372C9.72813 19.712 10.1037 19.6494 10.4542 19.521C10.7652 19.407 11.0608 19.2549 11.3343 19.068C11.6425 18.8575 11.9118 18.5882 12.4503 18.0497L20 10.5C21.3807 9.11929 21.3807 6.88071 20 5.5C18.6193 4.11929 16.3807 4.11929 15 5.5L7.45026 13.0497C6.91175 13.5882 6.6425 13.8575 6.43197 14.1657C6.24513 14.4392 6.09299 14.7348 5.97903 15.0458C5.85062 15.3963 5.78802 15.7719 5.66282 16.5231Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path><path d="M14.5 7L18.5 11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg>Edit</button></span></div></div></div><div class="overflow-y-auto p-4" dir="ltr"><code class="!whitespace-pre language-lua"><span><span>log</span><span>(</span><span>"msg"</span><span>)                   </span><span>-- Mensaje normal</span><span>
</span><span>debug</span><span>(</span><span>"msg"</span><span>)                </span><span>-- Mensaje de depuración</span><span>
emit(</span><span>"event"</span><span>, </span><span>"data"</span><span>)       </span><span>-- Emitir eventos personalizados</span><span>
play_sound(</span><span>"ruta.wav"</span><span>)      </span><span>-- Reproducir audio</span><span>
is_key_pressed(</span><span>"W"</span><span>)         </span><span>-- Verifica si tecla está presionada</span><span>
move_player(dx, dy)         </span><span>-- Mueve al jugador</span><span>

</span><span>-- Hooks del juego:</span><span>
</span><span>function</span><span></span><span>on_start</span><span>()</span><span></span><span>end</span><span></span><span>-- Se llama una vez al iniciar</span><span>
</span><span>function</span><span></span><span>on_update</span><span>()</span><span></span><span>end</span><span></span><span>-- Se llama cada frame (~60 FPS)</span><span>
</span></span></code></div></div></pre>

---

### 🗂️ **Estructura sugerida**

<pre class="!overflow-visible" data-start="1792" data-end="1932"><div class="contain-inline-size rounded-md border-[0.5px] border-token-border-medium relative bg-token-sidebar-surface-primary"><div class="flex items-center text-token-text-secondary px-4 py-2 text-xs font-sans justify-between h-9 bg-token-sidebar-surface-primary dark:bg-token-main-surface-secondary select-none rounded-t-[5px]">bash</div><div class="sticky top-9"><div class="absolute bottom-0 right-0 flex h-9 items-center pr-2"><div class="flex items-center rounded bg-token-sidebar-surface-primary px-2 font-sans text-xs text-token-text-secondary dark:bg-token-main-surface-secondary"><span class="" data-state="closed"><button class="flex gap-1 items-center select-none px-4 py-1" aria-label="Copy"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="icon-xs"><path fill-rule="evenodd" clip-rule="evenodd" d="M7 5C7 3.34315 8.34315 2 10 2H19C20.6569 2 22 3.34315 22 5V14C22 15.6569 20.6569 17 19 17H17V19C17 20.6569 15.6569 22 14 22H5C3.34315 22 2 20.6569 2 19V10C2 8.34315 3.34315 7 5 7H7V5ZM9 7H14C15.6569 7 17 8.34315 17 10V15H19C19.5523 15 20 14.5523 20 14V5C20 4.44772 19.5523 4 19 4H10C9.44772 4 9 4.44772 9 5V7ZM5 9C4.44772 9 4 9.44772 4 10V19C4 19.5523 4.44772 20 5 20H14C14.5523 20 15 19.5523 15 19V10C15 9.44772 14.5523 9 14 9H5Z" fill="currentColor"></path></svg>Copy</button></span><span class="" data-state="closed"><button class="flex select-none items-center gap-1 px-4 py-1"><svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" class="icon-xs"><path d="M2.5 5.5C4.3 5.2 5.2 4 5.5 2.5C5.8 4 6.7 5.2 8.5 5.5C6.7 5.8 5.8 7 5.5 8.5C5.2 7 4.3 5.8 2.5 5.5Z" fill="currentColor" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"></path><path d="M5.66282 16.5231L5.18413 19.3952C5.12203 19.7678 5.09098 19.9541 5.14876 20.0888C5.19933 20.2067 5.29328 20.3007 5.41118 20.3512C5.54589 20.409 5.73218 20.378 6.10476 20.3159L8.97693 19.8372C9.72813 19.712 10.1037 19.6494 10.4542 19.521C10.7652 19.407 11.0608 19.2549 11.3343 19.068C11.6425 18.8575 11.9118 18.5882 12.4503 18.0497L20 10.5C21.3807 9.11929 21.3807 6.88071 20 5.5C18.6193 4.11929 16.3807 4.11929 15 5.5L7.45026 13.0497C6.91175 13.5882 6.6425 13.8575 6.43197 14.1657C6.24513 14.4392 6.09299 14.7348 5.97903 15.0458C5.85062 15.3963 5.78802 15.7719 5.66282 16.5231Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path><path d="M14.5 7L18.5 11" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path></svg>Edit</button></span></div></div></div><div class="overflow-y-auto p-4" dir="ltr"><code class="!whitespace-pre"><span><span>mod/
├── player/init.lua
├── enemy/slime.lua
├── items/potions.lua
├── world.lua
└── game/init.lua   ← aquí van on_start y on_update
</span></span></code></div></div></pre>


## ❤️ Made with amor by...

A crazy paisa who codes with café in hand ☕🎮

---

¡Gracias por usar  **Luengo Engine** , mi amor!
