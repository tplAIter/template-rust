# Rust source acceptance boundaries

This source contains five generators and both layer_files/per_entity declarations, including SQLx create/read/update/delete and repository/use-case/service/controller bridges.
The root native v1 contract has empty dependencies. No Base composition, installed managed-export enrollment, Rust formatter grant, build action or database permission follows from these declarations.
Only the combined lifecycle descriptor is a prospective provider; layer/per-entity descriptors are references, not simultaneous providers. The observability profile remains planned and is not an active export.
The formatter options digest describes exact local edition2024 options and is content metadata, not a runtime profile approved by core.
Checks materialize public templates using the existing Go text/template syntax; this is source testing, not installed CLI/MCP execution.
Generated bridge tests use an in-memory test repository solely to exercise policy and real port implementations. SQLx database persistence requires separately approved disposable PostgreSQL acceptance, never an absent-DATABASE_URL silent pass.
Offline Cargo failures are reported as cache or compile refusals. No dependency acquisition or MSRV certification is implied by the installed Rust1.98 toolchain.
OpenAPI, OIDC, OTLP and brokers remain unavailable. Full P11 stays open pending installed source/tool/lifecycle and database acceptance.
