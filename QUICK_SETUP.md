# 🚀 Setup Rápido - Linux

## Instalación de Dependencias (Una línea)

### Ubuntu/Debian:
```bash
sudo apt update && sudo apt install -y libasound2-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libgl1-mesa-dev libglu1-mesa-dev libxxf86vm-dev pkg-config
```

### Fedora:
```bash
sudo dnf install -y alsa-lib-devel libX11-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel mesa-libGL-devel mesa-libGLU-devel libXxf86vm-devel pkgconfig
```

### Arch Linux:
```bash
sudo pacman -S alsa-lib libx11 libxrandr libxcursor libxinerama libxi mesa glu libxxf86vm pkgconf
```

## Ejecutar el Juego
```bash
git clone <tu-repo>
cd luango
go run main.go
```

Para la guía completa, consulta [SETUP_LINUX.md](SETUP_LINUX.md)
