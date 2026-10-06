// Rust layer_files repository composition.
pub struct RepositoryComposition<R> {
    pub repository: R,
}

impl<R> RepositoryComposition<R> {
    pub fn new(repository: R) -> Self {
        Self { repository }
    }

    pub fn into_inner(self) -> R {
        self.repository
    }
}

pub fn register_repository_layer<R>(repository: R) -> RepositoryComposition<R> {
    RepositoryComposition::new(repository)
}