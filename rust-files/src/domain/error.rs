//! Domain errors deliberately contain no transport or persistence types.

#[derive(Clone, Debug, PartialEq, Eq, thiserror::Error)]
pub enum DomainError {
    #[error("entity was not found")]
    NotFound,
    #[error("entity validation failed: {0}")]
    Validation(String),
    #[error("entity storage failed")]
    Storage,
}
