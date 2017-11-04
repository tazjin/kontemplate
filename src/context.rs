use std::path::PathBuf;
use types::*;
use errors::Result;
use std::fs::File;

use serde_json;
use serde_yaml;
use liquid::Object;

/// This intermediate struct is used for deserialising context objects.
/// Several fields in the actual context type used later on are computed based on other
/// values and not directly parsed from the files, but I don't want to deal with them
/// as optional values.
#[derive(Debug, Deserialize)]
struct SerialisedContext {
    #[serde(rename = "context")]
    name: String,
    global: Object,
    // TODO: imports: Vec<String>,
    resource_sets: Vec<ResourceSet>,
}

pub fn merge(mut obj1: Object, obj2: Object) -> Object {
    obj2.into_iter().for_each(|(k, v)| {
        obj1.insert(k, v);
    });

    obj1
}

pub fn load_context_from_file(file_path: PathBuf) -> Result<Context> {
    let file = File::open(file_path)?;
    let context: SerialisedContext = serde_yaml::from_reader(file)?;

    unimplemented!()
}
