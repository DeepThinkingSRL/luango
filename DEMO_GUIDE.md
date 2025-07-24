# 🎮 Demo del Agente Generativo - Luango Engine

## Instrucciones de Prueba

### 1. Ejecutar el Motor
```bash
cd /home/ubuntu/apps/luango
go run main.go
```

### 2. Activar el Agente
- Presiona **F4** en la ventana del juego para activar el agente generativo
- Se abrirá la interfaz de línea de comandos en la consola

### 3. Pruebas Básicas

#### Generar un Enemigo
```bash
generate script "crea un slime venenoso que salte y tenga 75 de vida"
```

#### Ver Recursos Pendientes
```bash
pending
```

#### Previsualizar un Recurso
```bash
preview script_1640123456  # Usar el ID que te dé el comando anterior
```

#### Aplicar un Recurso
```bash
apply script_1640123456    # Esto creará el archivo en mod/enemy/
```

#### Cambiar Modo a Automático
```bash
mode automatic
```

#### Generar Item (se aplicará automáticamente)
```bash
generate script "poción de vida que cure 50 HP"
```

#### Ver Historial
```bash
history
```

#### Estado del Agente
```bash
status
```

### 4. Pruebas desde Lua

En el juego:
- Presiona **G** para generar un enemigo desde Lua
- Presiona **M** para cambiar el modo del agente
- Los resultados aparecerán en la consola

### 5. Verificar Archivos Generados

Los scripts generados aparecerán en:
```
mod/
├── enemy/     # Enemigos generados
├── items/     # Items generados  
├── player/    # Scripts de jugador
└── behavior/  # Comportamientos AI
```

### 6. Integración con el Juego

Los scripts generados se cargan automáticamente la próxima vez que ejecutes:
```bash
go run main.go
```

## Ejemplos de Prompts

### Enemigos
- "crea un dragón de hielo que congele al jugador"
- "goblin archer que dispare flechas desde lejos"
- "fantasma que atraviese paredes y aparezca aleatoriamente"

### Items
- "espada mágica que aumente el daño en 25"
- "escudo que reduzca el daño en 50%"
- "poción de velocidad que dure 30 segundos"

### Comportamientos
- "sistema de patrullaje para guardias"
- "AI que persiga al jugador cuando esté cerca"
- "comportamiento de huida cuando la vida esté baja"

## Comandos Útiles

```bash
# Ver ayuda completa
help

# Limpiar pantalla
clear

# Salir del agente
exit

# Cambiar modos
mode interactive
mode automatic  
mode manual

# Rechazar un recurso pendiente
reject script_1640123456
```

¡Disfruta experimentando con el motor generativo! 🚀
