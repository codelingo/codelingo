# Tenets

Tenets are patterns and anti-patterns for your software stack. All Tenets are grouped in Bundles.

Tenet Bundles are groups of related tenets that can usefully be applied to new projects, or serve as exemplars for the CodeLingo community.

Bundles using a single language are in that language's directory, cross domain bundles have their own directories.

## Writing a bundle

Bundles have a README to explain their purpose with an ideas section for tenets yet to be implemented. Each tenet should have an associated src file (or other example of queried domain) and an expected output at `./testdata/<tenet_name>.<language_extension>` and `./testdata/<tenet_name>.json` respectively.
<!-- TODO: build simple  `lingo test-tenet <dir>` command -->
