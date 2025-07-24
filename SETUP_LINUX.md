# 🐧 Guía de Setup para Linux - Luengo Engine

Esta guía te ayudará a configurar y ejecutar **Luengo Engine** en sistemas Linux (Ubuntu/Debian y distribuciones similares).

## 📋 Prerrequisitos

### 1. Go Programming Language
Asegúrate de tener Go instalado (versión 1.19 o superior):

```bash
# Verificar si Go está instalado
go version

# Si no está instalado, descárgalo desde:
# https://golang.org/dl/
```

### 2. Git
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install git

# Verificar instalación
git --version
```

## 🔧 Instalación de Dependencias del Sistema

Luengo Engine utiliza **Ebiten** para gráficos y **ALSA** para audio, por lo que necesitas instalar las siguientes dependencias:

### Ubuntu/Debian:
```bash
# Actualizar repositorios
sudo apt update

# Instalar dependencias necesarias
sudo apt install -y \
    libasound2-dev \
    libx11-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxinerama-dev \
    libxi-dev \
    libgl1-mesa-dev \
    libglu1-mesa-dev \
    libxxf86vm-dev \
    pkg-config
```

### Fedora/CentOS/RHEL:
```bash
# Fedora
sudo dnf install -y \
    alsa-lib-devel \
    libX11-devel \
    libXrandr-devel \
    libXcursor-devel \
    libXinerama-devel \
    libXi-devel \
    mesa-libGL-devel \
    mesa-libGLU-devel \
    libXxf86vm-devel \
    pkgconfig

# CentOS/RHEL (con EPEL)
sudo yum install -y \
    alsa-lib-devel \
    libX11-devel \
    libXrandr-devel \
    libXcursor-devel \
    libXinerama-devel \
    libXi-devel \
    mesa-libGL-devel \
    mesa-libGLU-devel \
    libXxf86vm-devel \
    pkgconfig
```

### Arch Linux:
```bash
sudo pacman -S \
    alsa-lib \
    libx11 \
    libxrandr \
    libxcursor \
    libxinerama \
    libxi \
    mesa \
    glu \
    libxxf86vm \
    pkgconf
```

## 🚀 Instalación y Ejecución

### 1. Clonar el repositorio
```bash
git clone https://github.com/tu-usuario/luengo.git
cd luengo
```

### 2. Descargar dependencias de Go
```bash
go mod download
```

### 3. Compilar y ejecutar
```bash
go run main.go
```

Si todo está configurado correctamente, deberías ver algo como:
```
[Engine] Starting Luengo Engine
[Mod] Loading: mod/enemy/slime.lua
[Lua]: 👾 Slime enemy module loaded!
[Mod] Loading: mod/game/init.lua
[Mod] Loading: mod/items/potions.lua
[Lua]: 🧪 Potions module ready!
[Mod] Loading: mod/player/init.lua
[Lua]: 👤 Player module loaded!
[Mod] Loading: mod/world.lua
[Lua]: 🌍 World module initialized!
[Lua]: 🌟 Game is starting, mi amor!
```

## 🛠 Compilación para Distribución

### Compilar binario optimizado:
```bash
go build -ldflags="-s -w" -o luengo main.go
```

### Crear binario con información de debug:
```bash
go build -o luengo-debug main.go
```

## ❌ Resolución de Problemas Comunes

### Error: "Package alsa was not found"
```bash
# Instalar ALSA development headers
sudo apt install libasound2-dev

# O en Fedora:
sudo dnf install alsa-lib-devel
```

### Error: "X11/extensions/Xrandr.h: No such file or directory"
```bash
# Instalar X11 development headers
sudo apt install libxrandr-dev libx11-dev

# O en Fedora:
sudo dnf install libXrandr-devel libX11-devel
```

### Error: "cannot find -lXxf86vm"
```bash
# Instalar libXxf86vm development headers
sudo apt install libxxf86vm-dev

# O en Fedora:
sudo dnf install libXxf86vm-devel
```

### Error: "cannot find -lGL"
```bash
# Instalar OpenGL development headers
sudo apt install libgl1-mesa-dev libglu1-mesa-dev

# O en Fedora:
sudo dnf install mesa-libGL-devel mesa-libGLU-devel
```

## 🎮 Configuración para Desarrollo

### Variables de entorno útiles:
```bash
# Para debug de audio
export ALSA_DEBUG=1

# Para debug de OpenGL
export MESA_DEBUG=1

# Para logs más detallados de Ebiten
export EBITEN_DEBUG=1
```

### Configuración de VSCode (opcional):
Si usas Visual Studio Code, crea `.vscode/settings.json`:
```json
{
    "go.buildTags": "debug",
    "go.testTimeout": "30s",
    "go.lintTool": "golangci-lint"
}
```

## 📝 Notas Adicionales

- **Rendimiento**: Si experimentas problemas de rendimiento gráfico, asegúrate de tener los drivers de tu tarjeta gráfica actualizados.
- **Audio**: ALSA es el sistema de audio predeterminado en la mayoría de distribuciones Linux.
- **Wayland**: Ebiten funciona tanto en X11 como en Wayland.

## 🆘 Obtener Ayuda

Si encuentras problemas:
1. Revisa los logs de salida para errores específicos
2. Verifica que todas las dependencias estén instaladas
3. Consulta la documentación de [Ebiten](https://ebitengine.org/)
4. Abre un issue en el repositorio con detalles del error

---

¡Disfruta desarrollando con **Luengo Engine**! 🎮✨
