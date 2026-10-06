// Rust layer_files controller composition.
pub fn register_controller_layer(
    base: axum::Router<crate::health::HealthState>,
    generated: axum::Router<crate::health::HealthState>,
) -> axum::Router<crate::health::HealthState> {
    base.merge(generated)
}