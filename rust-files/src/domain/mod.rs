//! Framework-independent domain contracts exposed to the application layer.

pub mod entity;
pub mod error;

pub use entity::Entity;
pub use error::DomainError;
