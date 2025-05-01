# Titanium Go Plugin Scaffolder - Initial Implementation

This PR implements the initial version of the Titanium plugin scaffolder, focusing on the basic plugin structure and development mode functionality.

## Related Issues
- [Create the scaffolder as a plugin using HashiCorp go-plugin](https://github.com/titan-syndicate/titanium/issues/49)

## Changes
- Added basic plugin structure using HashiCorp's go-plugin
- Implemented development mode for local testing
- Added interactive CLI using Cobra
- Set up Mage build system for development workflow
- Added proper logging and error handling
- Fixed template loading to use embedded filesystem
- Improved configuration handling with Viper
- Added test data in idiomatic Go project structure

## Development Workflow
The scaffolder can be run in two modes:

1. Development mode (for local testing):
```bash
mage dev
```

2. Plugin mode (when loaded by Titanium):
```bash
mage run
```

## Build Commands
- `mage build` - Build the plugin
- `mage dev` - Run in development mode
- `mage run` - Run as a plugin
- `mage clean` - Clean build artifacts
- `mage install` - Install the plugin
- `mage scaffold` - Run scaffold with test configuration

## Testing
The development mode allows for local testing of the CLI functionality without needing the full Titanium environment. The scaffolder currently uses test data from `testdata/scaffold.yaml` for configuration.

## Next Steps
- Implement plugin scaffolding logic
- Add more interactive prompts
- Create template files for generated plugins
- Add validation for plugin names
- Implement remaining features from the ticket