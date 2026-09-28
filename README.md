<p align="center">
  <img src="https://raw.githubusercontent.com/tplAIter/.github/main/assets/banner.png?v=20260928" alt="tplAIter — Build with blocks. Spend fewer tokens." width="100%">
</p>

<h1 align="center">template-rust</h1>

<p align="center">A Rust service template rendered through tplAIter's existing Go engine.</p>

<p align="center"><strong>Status: public development preview · local manifest</strong></p>

<p align="center"><a href="https://github.com/tplAIter/tplaiter">core CLI</a> · <a href="https://github.com/tplAIter/template-base">base template</a> · <a href="https://github.com/tplAIter/template-go">Go template</a> · <a href="https://github.com/tplAIter/tplaiter/blob/main/docs/template-validation.md">validation workflow</a></p>

`template-rust` provides a Rust service layout rendered by Go's `text/template` engine. It covers service, repository, controller, worker, migration, and observability assets; integration and lifecycle work is still in progress.

In the intended MCP-assisted workflow, an agent selects parameters and blocks
instead of retyping a service foundation. The rendered result remains ordinary
Rust, ready for local checks and review. This preview does not claim that the
live MCP or project lifecycle is complete.

## Capabilities

- Manifest-driven database settings, including PostgreSQL options.
- Service, repository, controller, worker, migration, and observability render assets.
- A native template contract with no declared external dependencies.

## Verification

The repository workflow runs on pushes, pull requests, and manual dispatch. It first uses the pinned core template-check action to validate the manifest and render every fixture combination. Each rendered Rust project is then checked with:

```sh
cargo test --all-targets --locked
cargo check --all-targets --locked
```

The workflow pins Rust 1.98.0 for CI validation; that pin does not establish a minimum supported Rust version for generated projects. The shared checker validates template files and does not run template hooks or manifest commands; its boundaries are documented in the [core validation contract](https://github.com/tplAIter/tplaiter/blob/main/docs/template-validation.md).

## Status

This is a local `0.0.0-local` preview. Further template validation and lifecycle integration are still in progress; this repository does not claim complete production readiness or a direct `new`/`update` workflow by itself.
