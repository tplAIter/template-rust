//! Application use-case orchestration over a repository port.

use crate::{domain::{DomainError, Entity}, repository::EntityRepository};

pub struct GetEntity<R> {
    repository: R,
}

// CODEGEN:CRUD_APPLICATIONS

impl<R> GetEntity<R>
where
    R: EntityRepository,
{
    pub fn new(repository: R) -> Self { Self { repository } }

    pub async fn execute(&self, id: &str) -> Result<Entity, DomainError> {
        if id.trim().is_empty() {
            return Err(DomainError::Validation("id must not be empty".into()));
        }
        self.repository.find(id).await?.ok_or(DomainError::NotFound)
    }
}
