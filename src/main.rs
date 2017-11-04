#[macro_use] extern crate error_chain;
#[macro_use] extern crate serde_derive;

extern crate serde_json;
extern crate serde;
extern crate serde_yaml;
extern crate liquid;

mod errors {
    error_chain! {
        foreign_links {
            Io(::std::io::Error);
            Serialization(::serde_yaml::Error);
            Templating(::liquid::Error);
        }

        errors {
            EmptyTemplate(path: String) {
                description("an empty template was passed"),
                display("the template in '{}' was empty", path),
            }
        }
    }
}

mod context;
mod types;
mod templater;

use std::path;

fn main() {
    let rs = types::ResourceSet{
        name: "test".into(),
        path: "test".into(),
        values: Default::default(),
        include: Vec::new(),
        parent: "foo".into(),
    };
    let res = templater::template_resource_set(&rs).expect("Oioi");

    for resource in res.resources {
        println!("name: {}\n{}", resource.file_name, resource.rendered)
    }
}
