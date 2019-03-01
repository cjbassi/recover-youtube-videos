use crate::Result;

pub fn json_to_file<T>(filename: &str, json: &T) -> Result<()>
where
    T: serde::ser::Serialize,
{
    let j = serde_json::to_string_pretty(json)?;
    std::fs::write(filename, j)?;
    Ok(())
}
