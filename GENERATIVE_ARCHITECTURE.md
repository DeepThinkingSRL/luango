# 🤖 Luango Generative Engine - Arquitectura del Motor Generativo

## 🎯 Visión General

**Luango Generative Engine** es un motor de videojuegos revolucionario que combina:
- 🔧 **Motor base en Go** para rendimiento y estabilidad
- 🧠 **Scripting en Lua** para flexibilidad y modding
- 🤖 **Agente Generativo** para crear recursos automáticamente
- 🎮 **Modo híbrido** que permite tanto generación automática como edición manual

## 🏗 Arquitectura

```
Luango Engine
├── 🎮 Core Engine (Go)
│   ├── Entity Management
│   ├── Graphics (Ebiten)
│   ├── Audio System
│   └── Input Handling
│
├── 🧠 Scripting Layer (Lua)
│   ├── Game Logic
│   ├── Entity Behaviors
│   ├── Item Systems
│   └── World Scripts
│
└── 🤖 Generative Agent
    ├── Resource Generators
    ├── AI Integration
    ├── Template System
    └── Interactive CLI
```

## 🤖 Sistema de Agente Generativo

### Modos de Operación

1. **🔄 Interactive Mode** (Predeterminado)
   - Genera recursos y espera confirmación
   - Permite previsualización antes de aplicar
   - Control total sobre qué se aplica

2. **⚡ Automatic Mode**
   - Genera y aplica recursos automáticamente
   - Ideal para prototipado rápido
   - Menos control, mayor velocidad

3. **✋ Manual Mode**
   - Solo genera, nunca aplica automáticamente
   - Para desarrolladores que prefieren control total
   - Útil para inspección y aprendizaje

### Tipos de Recursos Generables

- 📜 **Scripts**: Enemigos, jugadores, sistemas, comportamientos
- 🎨 **Sprites**: Personajes, objetos, efectos (futuro)
- 🔊 **Audio**: Música, efectos de sonido (futuro)
- 🗺️ **Niveles**: Mapas, terrenos, puzzles (futuro)
- 💬 **Diálogos**: Conversaciones, narrativa (futuro)
- 🎭 **Animaciones**: Movimientos, efectos visuales (futuro)

## 🎮 Controles del Motor

### En el Juego
- **WASD/Flechas**: Mover jugador
- **F3**: Toggle información de debug
- **F4**: Activar/desactivar agente generativo

### En la Consola del Agente
```bash
# Comandos principales
generate <tipo> <prompt>     # Genera un recurso
apply <id>                   # Aplica un recurso pendiente
reject <id>                  # Rechaza un recurso
preview <id>                 # Previsualiza un recurso

# Gestión
pending                      # Lista recursos pendientes
history                      # Historial de generaciones
status                       # Estado del agente
mode <modo>                  # Cambia modo (interactive/automatic/manual)

# Ayuda
help                         # Muestra ayuda completa
exit                         # Salir del agente
```

## 📝 Ejemplos de Uso

### Desde la Consola del Agente

```bash
# Generar un enemigo
generate script "crea un dragón de fuego que escupa llamas y tenga 200 de vida"

# Generar un item
generate script "poción de maná que restaure 50 puntos de magia"

# Generar comportamiento
generate script "sistema de patrullaje para guardias"

# Cambiar modo
mode automatic

# Ver recursos pendientes
pending

# Aplicar un recurso
apply script_1640123456
```

### Desde Scripts Lua

```lua
-- Generar un enemigo desde Lua
local resultID = generate("script", "create a ice wizard enemy")

-- Aplicar resultado si está en modo interactivo
if resultID then
    apply_generated(resultID)
end

-- Cambiar modo del agente
set_agent_mode("automatic")
```

## 🔧 Sistema de Templates

El generador de scripts Lua utiliza templates inteligentes:

### Template de Enemigo
```lua
{{.Name}} = {
    name = "{{.Name}}",
    health = {{.Health}},
    speed = {{.Speed}},
    damage = {{.Damage}},
    state = "idle"
}

function {{.Name}}.update(dt)
    -- AI behavior logic
end

function {{.Name}}.attack()
    -- Attack logic
end
```

### Template de Item
```lua
{{.Name}} = {
    name = "{{.Name}}",
    type = "{{.Type}}",
    rarity = "{{.Rarity}}"
}

function {{.Name}}.use(player)
    -- Item usage logic
end
```

## 🎯 Flujo de Trabajo del Agente

1. **📝 Prompt del Usuario**
   - Usuario describe lo que quiere generar
   - Puede ser en lenguaje natural

2. **🧠 Análisis del Prompt**
   - El agente analiza el tipo de recurso
   - Extrae parámetros y características

3. **⚙️ Generación**
   - Selecciona template apropiado
   - Llena datos basado en el prompt
   - Genera código/contenido

4. **👀 Previsualización** (Modo Interactive)
   - Muestra preview del recurso
   - Permite inspección antes de aplicar

5. **✅ Aplicación**
   - Crea archivos en el proyecto
   - Integra con el sistema existente

6. **🔄 Feedback Loop**
   - Registra en historial
   - Ejecuta callbacks
   - Notifica resultados

## 🚀 Extensibilidad

### Agregar Nuevos Generadores

```go
// Crear generador personalizado
type SpriteGenerator struct {
    // Implementación
}

func (sg *SpriteGenerator) Generate(ctx context.Context, request GenerationRequest) (*GenerationResult, error) {
    // Lógica de generación de sprites
}

// Registrar en el agente
spriteGen := NewSpriteGenerator()
agent.RegisterGenerator(ResourceSprite, spriteGen)
```

### Agregar Callbacks

```go
// Registrar callback para notificaciones
agent.RegisterCallback(ResourceScript, func(result *GenerationResult) {
    log.Printf("Script generado: %s", result.Request.Prompt)
    // Enviar notificación, actualizar UI, etc.
})
```

## 🛠 Configuración del Proyecto

### Estructura de Archivos
```
luango/
├── main.go                 # Motor principal
├── engine/
│   └── generator/          # Sistema de agente generativo
│       ├── agent.go        # Núcleo del agente
│       ├── lua_generator.go # Generador de scripts Lua
│       └── cli.go          # Interfaz de línea de comandos
├── mod/                    # Scripts Lua del juego
│   ├── player/
│   ├── enemy/
│   ├── items/
│   └── world.lua
└── assets/                 # Recursos del juego
    ├── sprites/
    └── sounds/
```

### Dependencias
```go
require (
    github.com/hajimehoshi/ebiten/v2  // Motor gráfico
    github.com/yuin/gopher-lua        // Runtime Lua
    github.com/faiface/beep           // Sistema de audio
)
```

## 🌟 Características Únicas

1. **🔄 Integración Seamless**: El agente se integra nativamente con el motor
2. **🎮 Runtime Generation**: Genera recursos mientras el juego está ejecutándose
3. **🧠 Smart Templates**: Templates que se adaptan al contexto del prompt
4. **📜 Lua-First**: Los scripts generados son Lua nativo, totalmente modificables
5. **🔧 Híbrido**: Combina generación automática con edición manual
6. **📚 Historial**: Rastrea todas las generaciones para aprendizaje
7. **🎯 Contextual**: Entiende la estructura del proyecto para generar apropiadamente

## 🚧 Roadmap

### Próximas Características
- [ ] Generador de sprites con IA
- [ ] Generador de música y sonidos
- [ ] Sistema de diálogos generativo
- [ ] Generación de niveles procedurales
- [ ] Editor visual integrado
- [ ] Integración con modelos de IA externos
- [ ] Sistema de versionado de recursos
- [ ] Marketplace de templates
- [ ] Modo colaborativo multi-usuario

---

**Luango Generative Engine** representa el futuro del desarrollo de videojuegos, donde la creatividad humana se combina con la potencia de la generación automática para crear experiencias únicas e innovadoras. 🎮✨
