# Task Spec Index

Task specs are detailed, focused documents for complex or cross-cutting capabilities. They complement
the broader prompts in `prompts/`.

| # | Task | Description | Related Prompts |
|---|------|-------------|-----------------|
| T01 | [SenangPay Signature](senangpay-signature.md) | Emulated SenangPay MD5 signature scheme, charge/callback/webhook flow, config, and Go helper code | 12 |
| T02 | [OpenMuara-to-OpenMuara Migration Guide](openmuara-migration-guide.md) | Back up, migrate, and verify a legacy `muara` workspace under the current `openmuara` layout | 18 |

## When to Use a Task Spec

Use a task spec when a prompt needs deeper detail on a single complex topic before implementation.
Task specs should be referenced from their related prompt.
