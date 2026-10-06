// Rust per_entity controller composition.
pub fn register_controller_entity(
    base: axum::Router<crate::health::HealthState>,
    generated: axum::Router<crate::health::HealthState>,
) -> axum::Router<crate::health::HealthState> {
    base.merge(generated)
}