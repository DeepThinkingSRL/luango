# 📘 Luango Engine – Technical Documentation

## 🧠 Overview

Luango Engine is a modular, lightweight, scriptable 2D game engine built with Go and Lua. It’s designed for high-performance gameplay, real-time modding, and full Lua integration.

---

## 🧱 Architecture

* **Language Base** : Go (Core Engine), Lua (Scripting)
* **Rendering** : Ebiten
* **Audio** : Beep
* **Scripting** : Gopher-Lua
* **Debug HUD** : Built-in (toggle with F3)

---

## 🔁 Game Lifecycle

* `on_start()` – Called once on game start (optional)
* `on_update()` – Called every frame (~60 FPS)

---

## 📦 Folder Structure

```
luango/
├── assets/              # Sprites, audio, etc.
├── mod/                 # Lua scripts and mods
│   ├── player/
│   ├── enemy/
│   ├── items/
│   ├── game/init.lua    # Entry script with on_start and on_update
├── main.go              # Engine entry point
├── README.md
├── LICENSE
└── TECHNICAL.md         # This doc
```

---

## 📜 Lua API Functions

| Function                 | Description                     |
| ------------------------ | ------------------------------- |
| `log(msg)`             | Prints a message to the console |
| `debug(msg)`           | Prints debug info               |
| `emit(event, payload)` | Emits a custom event            |
| `play_sound(path)`     | Plays a `.wav`audio file      |
| `is_key_pressed(key)`  | Returns `true/false`for key   |
| `move_player(dx, dy)`  | Moves the main player entity    |

---

## 🧪 Debug Tools

* Press `F3` to toggle in-game debug overlay.
* Shows Player position, frame count.
* Supports `debug()` in Lua.

---

## ⚙️ Modding

Mods are simply `.lua` files placed anywhere inside the `mod/` folder or subfolders. They're all loaded automatically at startup.

Example:

```lua
function on_update()
  if is_key_pressed("ArrowRight") then
    move_player(1, 0)
  end
end
```

---

## 🤝 Contributing

* Follow Go’s idiomatic conventions.
* Keep Lua logic modular (separate files per feature).
* Pull requests welcome!

---

¡Hecho con amor en Luango Engine, mi amorsh! 💙
