# BO - Build Out

`bo` is the build-out exporter for GoScript.

It is not part of `gopm`. `gopm` can stay focused on project/tooling workflows, while `bo` turns a manifest into deployable output.

## What `bo` does

- Builds a selected tool, module, or app entrypoint into an executable
- Packages that executable into a portable `goe` bundle
- Generates packaging scaffolds for `apk`, `ipa`, and `dmg`
- Writes a normalized manifest and build plan beside the output so the export is inspectable by humans and agents

## Manifest format

`bo` expects a JSON manifest file, usually with the `.manifest` extension.

```json
{
  "name": "admin",
  "module": ".",
  "entrypoint": "./cmd/server",
  "output": "admin",
  "paths": ["/admin", "/admin/users"],
  "folders": ["pkg/components/admin"],
  "assets": ["static/admin"]
}
```

## Usage

```bash
bo admin.manifest - exe
bo admin.manifest - goe
bo calc.manifest - exe
bo admin.manifest - apk
bo admin.manifest - ipa
bo admin.manifest - dmg
```

## Output contract

- `exe` builds a host executable with `go build`
- `goe` builds a portable bundle that includes the executable plus manifest metadata
- `apk`, `ipa`, and `dmg` currently generate scaffolds that can be wired to native packagers later

## Why this exists

`bo` is for the ERP-style future where a developer can export just the tool they want, rather than splitting the whole system into a single monolith or forcing everything through a browser-only deployment model.

