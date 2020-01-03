# Validation

This package can be used to validate forensic artifact definition files.

Implemented validations:

- Files should end with ".yaml".
- First line of the file should be a comment.
- Lines should not end with whitespace.
- Artifact definitions in files with OS prefix should only be for that OS.
- Artifact definition names should be unique.
- Artifact definition names in files with OS prefix should have an OS prefix.
- Artifact definition names should be CamelCase.
- Artifact definition names should not contain spaces.
- Artifact definition names should end in type (except WMI and ARTIFACT_GROUPS).
- Long docs should contain an empty second line.
- Artifact definition should have sources.
- Types should be one of a predefined list.
- Label should be one of a predefined list.
- Supported OS should be one of a predefined list in artifact definitions and sources.
- Provides should match a parameter in artifact definitions and sources.
- Attributes should only be set according to source type.
- Some attributes are required by source type.
- REGISTRY_* and WMI source types are only allowed in windows sources.
- Paths in MacOS should have linked path as well (eg. /etc and /private/etc version).
- Some paths (e.g. "%%users.userprofile%%\\AppData\\Local") should use a shorter version (e.g. "%%users.localappdata%%").
- Paths should not contains "**".
- "%%users.homedir%%" should not be used on Windows.
- Registry keys should not start with %%CURRENT_CONTROL_SET%%.
- Registry keys do not support HKEY_CURRENT_USER.
- Registry keys should be unique.
- Registry values should be unique.
- Artifact group should not references itself.
- Artifact groups should not contain cyclic references.
- Artifact groups should only referece other existing artifacts.
- Parameter should be one of a predefined list.
- Parameter should be provided by at least one provides.

Missing [style guide](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc) validations:

- Multi-line documentation should use the YAML Literal Style as indicated by the | character ([2.2.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#22-long-docs-form)).
- Explicit newlines (\n) should not be used in doc ([2.2.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#22-long-docs-form)).
- Where sources take a single argument with a single value, the one-line {} form should be used to save on line breaks ([3.0.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#3-sources)).
- Require the args attribute for commands ([3.3.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#33-command-source)).
- Generally use the short [] format for single-item lists that fit inside 80 characters to save on unnecessary line breaks ([6.2.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#62-lists)).
- Quotes should not be used for doc strings, artifact names, and simple lists like labels and supported_os  ([6.3.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#63-quotes)).
- Paths and URLs should use single quotes to avoid the need for manual escaping ([6.3.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#63-quotes)).
- Double quotes should be used where escaping causes problems, such as regular expressions ([6.3.](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#63-quotes)).
- To minimize the number of artifacts in the list, combine them using the supported_os and conditions attributes where it makes sense ([6.4](https://github.com/ForensicArtifacts/artifacts/blob/master/docs/Artifacts%20definition%20format%20and%20style%20guide.asciidoc#64-minimize-the-number-of-definitions-by-using-multiple-sources)).
