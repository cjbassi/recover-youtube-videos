use std::path::PathBuf;
use structopt::StructOpt;

#[derive(StructOpt)]
pub struct Args {
    #[structopt(short = "d", long = "debug")]
    pub debug: bool,

    pub directory: PathBuf,
}
