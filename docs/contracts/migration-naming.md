# SQLx migration naming

Both migration and crud require an explicit table matching `^[a-z_][a-z0-9_]*$`.
The title is never parsed or pluralized. Numeric sequences precede the title; the retained `numbered: goose` identifier selects numbering only.
SQL bodies are forward-only SQLx migrations. There is no Goose Down directive or DROP appended to the forward file.
Rollback needs a separately approved migration; destructive rollback is never generated implicitly.
