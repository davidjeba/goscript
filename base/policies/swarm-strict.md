# Swarm Strict Policy

Use this policy when `sw` mode runs with `strict` behavior.

## Rules

- A sensitive module may have only one active authorized server identity at a time
- Server identity must be cryptographic, not only an IP address
- IP address is routing metadata, not the trust anchor
- If an authorized server changes IP, the module must be revalidated before it runs again
- High-trust modules such as `pay` should stop running until the new endpoint is approved
- Every rebind, revoke, or renewal must be audit logged

