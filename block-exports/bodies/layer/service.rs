// Rust layer_files service composition.
pub struct ServiceComposition<S> {
    pub service: S,
}

impl<S> ServiceComposition<S> {
    pub fn new(service: S) -> Self {
        Self { service }
    }

    pub fn into_inner(self) -> S {
        self.service
    }
}

pub fn register_service_layer<S>(service: S) -> ServiceComposition<S> {
    ServiceComposition::new(service)
}