package main

import (
        "fmt"
        "os"
        "strings"

        "github.com/davidjeba/goscript/cmd/gopm/commands"
        "github.com/davidjeba/goscript/pkg/gopm"
)

func main() {
        if len(os.Args) < 2 {
                printHelp()
                os.Exit(1)
        }

        command := os.Args[1]
        args := os.Args[2:]

        pm := gopm.NewPackageManager()

        switch command {
        // Jetpack commands
        case "jetpack":
                commands.JetpackCommand(args)
                return
        // Basic package management
        case "get":
                pm.Get(args)
        case "update":
                pm.Update(args)
        case "clean":
                pm.Clean(args)
        case "run":
                pm.Run(args)
        case "audit":
                pm.Audit(args)
        case "publish":
                pm.Publish(args)
        case "version":
                pm.Version(args)
        case "cache-clear":
                pm.CacheClear(args)
        case "list":
                pm.List(args)
        case "verify":
                pm.Verify(args)
        case "dedupe":
                pm.Dedupe(args)
        case "prune":
                pm.Prune(args)
        case "config":
                pm.ConfigCmd(args)
        case "help":
                pm.Help(args)
        case "auth":
                pm.Auth(args)
        case "setup":
                pm.Setup(args)
        case "sync":
                pm.Sync(args)
        case "doctor":
                pm.Doctor(args)
        case "migrate":
                pm.Migrate(args)
        case "rollback":
                pm.Rollback(args)

        // Gocsx CSS framework commands
        case "css:build":
                pm.CSSBuild(args)
        case "css:watch":
                pm.CSSWatch(args)
        case "css:optimize":
                pm.CSSOptimize(args)
        case "css:analyze":
                pm.CSSAnalyze(args)
        case "css:theme":
                pm.CSSTheme(args)

        // WebGPU and 3D commands
        case "webgpu:init":
                pm.WebGPUInit(args)
        case "webgpu:build":
                pm.WebGPUBuild(args)
        case "webgpu:optimize":
                pm.WebGPUOptimize(args)
        case "3d:scene":
                pm.Scene3DCreate(args)
        case "3d:model":
                pm.Model3DImport(args)
        case "3d:export":
                pm.Model3DExport(args)
        case "3d:optimize":
                pm.Model3DOptimize(args)
        case "3d:convert":
                pm.Model3DConvert(args)

        // 2D Canvas commands
        case "2d:init":
                pm.Canvas2DInit(args)
        case "2d:sprite":
                pm.SpriteCreate(args)
        case "2d:animation":
                pm.AnimationCreate(args)
        case "2d:atlas":
                pm.AtlasCreate(args)
        case "2d:optimize":
                pm.Canvas2DOptimize(args)

        // GoUIX commands
        case "uix:init":
                pm.UIXInit(args)
        case "uix:component":
                pm.UIXComponentCreate(args)
        case "uix:test":
                pm.UIXTest(args)
        case "uix:storybook":
                pm.UIXStorybook(args)
        case "uix:build":
                pm.UIXBuild(args)

        // GoScale API commands
        case "api:init":
                pm.APIInit(args)
        case "api:schema":
                pm.APISchemaCreate(args)
        case "api:deploy":
                pm.APIDeploy(args)
        case "api:edge":
                pm.APIEdgeDeploy(args)
        case "api:test":
                pm.APITest(args)
        case "api:doc":
                pm.APIDocGenerate(args)

        // GoScale DB commands
        case "db:init":
                pm.DBInit(args)
        case "db:migrate":
                pm.DBMigrate(args)
        case "db:seed":
                pm.DBSeed(args)
        case "db:backup":
                pm.DBBackup(args)
        case "db:restore":
                pm.DBRestore(args)
        case "db:schema":
                pm.DBSchemaCreate(args)
        case "db:timeseries":
                pm.DBTimeSeriesEnable(args)

        default:
                fmt.Printf("Unknown command: %s\n", command)
                printHelp()
                os.Exit(1)
        }
}

func printHelp() {
        help := `
GOPM - Go Package Manager

Usage: gopm [command] [options]

Basic Commands:
  get           Install packages
  update        Update packages
  clean         Clean project
  run           Run a script
  audit         Check for vulnerabilities
  publish       Publish a package
  version       Show version information
  cache-clear   Clear the cache
  list          List installed packages
  verify        Verify package integrity
  dedupe        Remove duplicate packages
  prune         Remove unused packages
  config        Manage configuration
  help          Show help
  auth          Authenticate with registry
  setup         Setup project
  sync          Sync dependencies
  doctor        Diagnose and fix issues
  migrate       Migrate to a new version
  rollback      Rollback to a previous version

Gocsx CSS Framework Commands:
  css:build     Build CSS
  css:watch     Watch and rebuild CSS
  css:optimize  Optimize CSS
  css:analyze   Analyze CSS usage
  css:theme     Manage themes

WebGPU and 3D Commands:
  webgpu:init     Initialize WebGPU project
  webgpu:build    Build WebGPU shaders
  webgpu:optimize Optimize WebGPU performance
  3d:scene        Create 3D scene
  3d:model        Import 3D model
  3d:export       Export 3D model
  3d:optimize     Optimize 3D model
  3d:convert      Convert between 3D formats

2D Canvas Commands:
  2d:init         Initialize 2D canvas project
  2d:sprite       Create sprite
  2d:animation    Create animation
  2d:atlas        Create sprite atlas
  2d:optimize     Optimize 2D canvas performance

GoUIX Commands:
  uix:init        Initialize UIX project
  uix:component   Create UIX component
  uix:test        Test UIX components
  uix:storybook   Start UIX storybook
  uix:build       Build UIX project

GoScale API Commands:
  api:init        Initialize API project
  api:schema      Create API schema
  api:deploy      Deploy API
  api:edge        Deploy to edge network
  api:test        Test API
  api:doc         Generate API documentation

GoScale DB Commands:
  db:init         Initialize database
  db:migrate      Run database migrations
  db:seed         Seed database
  db:backup       Backup database
  db:restore      Restore database
  db:schema       Create database schema
  db:timeseries   Enable time series features

Jetpack Performance Monitoring:
  jetpack         Performance monitoring and optimization:
    init          Initialize Jetpack
    monitor       Monitor performance
    lighthouse    Run Lighthouse audits
    panel         Manage performance panel
    metrics       Manage metrics
    security      Security monitoring
    export        Export metrics
    report        Generate reports
    chrome        Manage Chrome extension

For more information, run: gopm help [command]
`
        fmt.Println(strings.TrimSpace(help))
}