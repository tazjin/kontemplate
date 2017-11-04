use std::path::PathBuf;
use serde_json::{Map, Value};
use liquid;

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct ResourceSet {
    /// Name of the resource set. This can be used in include/exclude statements during
    /// kontemplate runs.
    pub name: String,

    /// Path to the folder containing the files for this resource set. This defaults to the
    /// value of the 'name' field if unset.
    pub path: PathBuf,

    /// Values to include when interpolating resources from this resource set.
    pub values: liquid::Object,

    /// Nested resource sets to include
    pub include: Vec<ResourceSet>,

    // Parent resource set for flattened resource sets. Should not be manually specified.
    // TODO: pub parent: String,
    // TODO: Note - parent is used for include/exclude checks.
}

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub struct Context {
    /// The name of the kubectl context
    pub name: String,

    /// Global variables that should be accessible by all resource sets
    pub global: liquid::Object,

    /// File names of YAML or JSON files including extra variables that should be globally
    /// accessible.
    pub imports: Vec<String>,

    /// The resource sets to include in this context
    #[serde(rename = "include")]
    pub resource_sets: Vec<ResourceSet>,

    // This field represents the absolute path to the context base directory and should not be
    // manually specified.
    // TODO: base_dir: PathBuf,
}
