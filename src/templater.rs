use liquid::{self, Renderable};
use types::*;
use errors::{Result, ErrorKind};
use std::path::PathBuf;

#[derive(Debug)]
pub struct RenderedResource {
    pub file_name: String,
    pub rendered: String,
}

#[derive(Debug)]
pub struct RenderedResourceSet {
    pub name: String,
    pub resources: Vec<RenderedResource>,
}

pub fn template_resource_set(rs: &ResourceSet) -> Result<RenderedResourceSet> {
    // Single-file resource sets just get rendered directly:
    if rs.path.is_file() {
        return Ok(RenderedResourceSet {
            name: rs.name.clone(),
            resources: vec![template_file(&rs, &rs.path)?],
        })
    }

    let dir = rs.path.read_dir()?;
    let mut resources = Vec::new();

    for file_result in dir {
        let file = file_result?;
        let path = file.path();
        let file_name = path.file_name()
            .and_then(|n| n.to_str())
            .expect("Invalid file name");

        if !is_default_file(file_name) {
            resources.push(template_file(&rs, &path)?)
        }
    }

    Ok(RenderedResourceSet{ name: rs.name.clone(), resources})
}

fn template_file(rs: &ResourceSet, file_path: &PathBuf) -> Result<RenderedResource> {
    let file_name = format!("{}", file_path.to_str().unwrap());
    let template = liquid::parse_file(file_path, Default::default())?;
    let mut context = liquid::Context::with_values(rs.values.clone());
    let result = template.render(&mut context)?;

    match result {
        Some(rendered) => Ok(RenderedResource {
            file_name,
            rendered,
        }),
        None => Err(ErrorKind::EmptyTemplate(file_name).into()),
    }
}

/// Filters out files that provide default values from templating.
fn is_default_file(file: &str) -> bool {
    vec!["default.yml", "default.yaml", "default.json"].contains(&file)
}
