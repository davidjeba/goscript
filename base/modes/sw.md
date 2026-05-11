# SW Mode

`sw` means swarm mode, defined as modular distributed server architecture.

This does not mean clone servers. It means modules are assigned to the best execution node for their job.

## Characteristics

- `app/task` may run on a remote server
- `app/pos` may run on a local offline server
- `app/pay` may run on a third-party secure server
- Each module has an execution owner, trust boundary, and fallback behavior

## Folder Extensions

`sw` keeps the shared tree and adds these folders:

```text
project/
  app/
    topology/
    sync/
    swarm-policies/
    trust/
```

## Rules

- Use `manifest.mode = "sw"`
- Do not describe the system as duplicate servers running the same app everywhere
- Define module placement explicitly
- Keep `strict` policies in `base/policies/swarm-strict.md`
- Reserve `swarm-policies/` for node assignment, trust, and routing rules

