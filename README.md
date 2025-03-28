# Luango Engine

**Luango Engine** is a lightweight, modular game engine written in Go with built-in scripting support using Lua. It is designed to be simple, flexible, and powerful — ideal for game developers who want performance and modding capabilities without the complexity.

## 🚀 Features

* 🔄 Written in Go for performance and concurrency
* 🧠 Embedded Lua scripting with [Gopher-Lua](https://github.com/yuin/gopher-lua)
* 🧩 Modular architecture for flexibility and expansion
* 🎮 Custom mod loading system
* 🧱 ECS-friendly foundation (coming soon)

## 📂 Project Structure

```
luango/
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
git clone https://github.com/yourname/luango.git
cd luango
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

## ❤️ Made with amor by...

A crazy paisa who codes with café in hand ☕🎮

---

¡Gracias por usar  **Luango Engine** , mi amor!
