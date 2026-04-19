# GOPM - Go Package Manager

GOPM is a comprehensive package manager for Go projects, with special support for GoScript ecosystem components including Gocsx, WebGPU, GoUIX, and GoScale.

## Features

- **Complete Package Management**: Install, update, and manage dependencies
- **Gocsx CSS Framework Support**: Build, optimize, and manage CSS themes
- **WebGPU and 3D Support**: Initialize WebGPU projects, build shaders, and manage 3D assets
- **2D Canvas Support**: Create and manage sprites, animations, and sprite atlases
- **GoUIX Support**: Create and test UI components
- **GoScale API Support**: Create and deploy APIs with edge computing capabilities
- **GoScale DB Support**: Manage databases with time series and NoSQL features

## Installation

```bash
go install github.com/davidjeba/goscript/cmd/gopm@latest
```

## Basic Usage

### Package Management

```bash
# Install packages
gopm get package1 package2

# Update packages
gopm update

# Run a script
gopm run start

# List installed packages
gopm list

# Check for vulnerabilities
gopm audit
```

### Configuration

```bash
# View configuration
gopm config

# Set configuration
gopm config registry https://registry.gopm.dev
```

## Gocsx CSS Framework Commands

```bash
# Build CSS
gopm css:build

# Watch and rebuild CSS
gopm css:watch

# Optimize CSS
gopm css:optimize

# Analyze CSS usage
gopm css:analyze

# Create a theme
gopm css:theme create dark
```

## WebGPU and 3D Commands

```bash
# Initialize WebGPU project
gopm webgpu:init

# Build WebGPU shaders
gopm webgpu:build

# Create a 3D scene
gopm 3d:scene

# Import a 3D model
gopm 3d:model model.glb

# Export a 3D model
gopm 3d:export model.glb model.obj

# Optimize a 3D model
gopm 3d:optimize model.glb

# Convert between 3D formats
gopm 3d:convert model.glb model.obj
```

## 2D Canvas Commands

```bash
# Initialize 2D canvas project
gopm 2d:init

# Create a sprite
gopm 2d:sprite player

# Create an animation
gopm 2d:animation walk

# Create a sprite atlas
gopm 2d:atlas game-sprites

# Optimize 2D canvas performance
gopm 2d:optimize
```

## GoUIX Commands

```bash
# Initialize UIX project
gopm uix:init

# Create a UIX component
gopm uix:component Button

# Test UIX components
gopm uix:test

# Start UIX storybook
gopm uix:storybook

# Build UIX project
gopm uix:build
```

## GoScale API Commands

```bash
# Initialize API project
gopm api:init

# Create API schema
gopm api:schema User

# Deploy API
gopm api:deploy

# Deploy to edge network
gopm api:edge

# Test API
gopm api:test

# Generate API documentation
gopm api:doc
```

## GoScale DB Commands

```bash
# Initialize database
gopm db:init

# Run database migrations
gopm db:migrate

# Seed database
gopm db:seed

# Backup database
gopm db:backup

# Restore database
gopm db:restore

# Create database schema
gopm db:schema users

# Enable time series features
gopm db:timeseries metrics
```

## Command Reference

### Basic Commands

| Command | Description |
|---------|-------------|
| `get` | Install packages |
| `update` | Update packages |
| `clean` | Clean project |
| `run` | Run a script |
| `audit` | Check for vulnerabilities |
| `publish` | Publish a package |
| `version` | Show version information |
| `cache-clear` | Clear the cache |
| `list` | List installed packages |
| `verify` | Verify package integrity |
| `dedupe` | Remove duplicate packages |
| `prune` | Remove unused packages |
| `config` | Manage configuration |
| `help` | Show help |
| `auth` | Authenticate with registry |
| `setup` | Setup project |
| `sync` | Sync dependencies |
| `doctor` | Diagnose and fix issues |
| `migrate` | Migrate to a new version |
| `rollback` | Rollback to a previous version |

### Gocsx CSS Framework Commands

| Command | Description |
|---------|-------------|
| `css:build` | Build CSS |
| `css:watch` | Watch and rebuild CSS |
| `css:optimize` | Optimize CSS |
| `css:analyze` | Analyze CSS usage |
| `css:theme` | Manage themes |

### WebGPU and 3D Commands

| Command | Description |
|---------|-------------|
| `webgpu:init` | Initialize WebGPU project |
| `webgpu:build` | Build WebGPU shaders |
| `webgpu:optimize` | Optimize WebGPU performance |
| `3d:scene` | Create 3D scene |
| `3d:model` | Import 3D model |
| `3d:export` | Export 3D model |
| `3d:optimize` | Optimize 3D model |
| `3d:convert` | Convert between 3D formats |

### 2D Canvas Commands

| Command | Description |
|---------|-------------|
| `2d:init` | Initialize 2D canvas project |
| `2d:sprite` | Create sprite |
| `2d:animation` | Create animation |
| `2d:atlas` | Create sprite atlas |
| `2d:optimize` | Optimize 2D canvas performance |

### GoUIX Commands

| Command | Description |
|---------|-------------|
| `uix:init` | Initialize UIX project |
| `uix:component` | Create UIX component |
| `uix:test` | Test UIX components |
| `uix:storybook` | Start UIX storybook |
| `uix:build` | Build UIX project |

### GoScale API Commands

| Command | Description |
|---------|-------------|
| `api:init` | Initialize API project |
| `api:schema` | Create API schema |
| `api:deploy` | Deploy API |
| `api:edge` | Deploy to edge network |
| `api:test` | Test API |
| `api:doc` | Generate API documentation |

### GoScale DB Commands

| Command | Description |
|---------|-------------|
| `db:init` | Initialize database |
| `db:migrate` | Run database migrations |
| `db:seed` | Seed database |
| `db:backup` | Backup database |
| `db:restore` | Restore database |
| `db:schema` | Create database schema |
| `db:timeseries` | Enable time series features |

## Configuration

GOPM uses a configuration file located at `~/.gopm/config.json` or in the project directory as `.gopmrc.json`.

Example configuration:

```json
{
  "registry": "https://registry.gopm.dev",
  "cache-dir": "~/.gopm/cache",
  "global-dir": "~/.gopm/global",
  "proxy": {
    "enabled": true,
    "url": "https://proxy.gopm.dev"
  },
  "timeout": 60,
  "retry-count": 3,
  "max-concurrent": 10,
  "strict-ssl": true,
  "save-exact": false,
  "production": false,
  "development": true,
  "ignore-scripts": false,
  "force-fetch": false,
  "offline-mode": false,
  "compression-level": 6
}
```

## Project Configuration

GOPM uses a `gopm.json` file in the project directory to manage dependencies and project configuration.

Example `gopm.json`:

```json
{
  "name": "my-project",
  "version": "1.0.0",
  "description": "My awesome project",
  "main": "main.go",
  "scripts": {
    "start": "go run main.go",
    "build": "go build -o app main.go",
    "test": "go test ./..."
  },
  "dependencies": {
    "github.com/davidjeba/goscript": "^1.0.0",
    "github.com/davidjeba/gocsx": "^1.0.0"
  },
  "devDependencies": {
    "github.com/stretchr/testify": "^1.8.0"
  },
  "engines": {
    "go": ">=1.18"
  },
  "gocsx": {
    "theme": "default",
    "breakpoints": {
      "sm": "640px",
      "md": "768px",
      "lg": "1024px",
      "xl": "1280px"
    }
  },
  "webgpu": {
    "shaders": "./shaders"
  },
  "goscale": {
    "api": {
      "port": 8080,
      "edge-enabled": true
    },
    "db": {
      "connection-string": "localhost:5432",
      "time-series-enabled": true
    }
  }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License