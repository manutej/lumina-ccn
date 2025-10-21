# LUMINA Configuration Management

## Overview

LUMINA uses a clear separation between version-controlled project configuration and user-specific runtime configuration.

## Configuration Locations

### Version Controlled (Git)

These files are tracked in git and shared across all users:

```
ccn/
├── version.go                    # Application version info
├── go.mod                        # Go dependencies
├── go.sum                        # Dependency checksums
├── .gitignore                    # Git ignore rules
├── .goreleaser.yml              # Release configuration
├── .github/workflows/           # CI/CD workflows
│   ├── ci.yml
│   └── release.yml
└── docs/                        # Documentation
    ├── GIT_WORKFLOW.md
    └── CONFIGURATION_MANAGEMENT.md
```

### User-Specific (NOT Version Controlled)

These files are created at runtime and customized per user:

```
~/.config/lumina/
└── keybindings.json             # Custom keybindings (auto-generated)
```

### Build Artifacts (Ignored)

These files are generated during build and ignored by git:

```
ccn                              # Binary (macOS/Linux)
ccn.exe                          # Binary (Windows)
dist/                            # GoReleaser output
build/                           # Build directory
*.test                           # Test binaries
*.out                            # Coverage files
```

## .gitignore Configuration

Current `.gitignore` rules:

```gitignore
# Binaries
ccn
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Build artifacts
dist/
build/
```

## User Configuration File Format

### keybindings.json

Location: `~/.config/lumina/keybindings.json`

This file is auto-generated on first run with default keybindings. Users can customize it without recompiling.

**Example:**
```json
{
  "keybindings": {
    "quit": ["q", "ctrl+c"],
    "help": ["?"],
    "up": ["k", "up"],
    "down": ["j", "down"],
    "page_up": ["u", "ctrl+u", "pgup"],
    "page_down": ["d", "ctrl+d", "pgdown"],
    "half_page_up": ["ctrl+b"],
    "half_page_down": ["ctrl+f"],
    "top": ["g", "home"],
    "bottom": ["G", "end"],
    "copy": ["y"],
    "toggle_preview": ["tab"],
    "search": ["/"],
    "next_match": ["n"],
    "prev_match": ["N"]
  }
}
```

**Rules:**
- ✅ Auto-created with defaults on first run
- ✅ User can modify freely
- ✅ NOT version controlled (.gitignore excludes ~/.config/)
- ✅ Application validates format on load
- ✅ Falls back to defaults if invalid

## Managing Configuration Across Environments

### Development Environment

```bash
# Development uses local config
cd ~/projects/lumina/ccn
go run . ~/.config

# User config: ~/.config/lumina/keybindings.json
```

### Production Build

```bash
# Build ignores user config, creates binary
go build -o lumina

# First run creates default user config
./lumina --init-config

# User customizes
vim ~/.config/lumina/keybindings.json
```

### Testing Environments

```bash
# Use custom config directory for tests
export LUMINA_CONFIG_DIR=/tmp/lumina-test-config
go test ./...
```

## Configuration Migration Strategy

### Version 1.0.x → 1.1.x

If keybindings format changes:

```go
// In keybindings.go
func LoadKeybindings() (*KeyBindings, error) {
    // Try loading user config
    userConfig := loadUserConfig()

    // Migrate old format if needed
    if userConfig.Version == "1.0" {
        userConfig = migrateV1ToV2(userConfig)
        saveUserConfig(userConfig)
    }

    return userConfig, nil
}
```

**Migration Steps:**
1. Detect old config format
2. Convert to new format
3. Save migrated config
4. Log migration message to user

### Breaking Changes

If configuration format changes incompatibly:

```bash
# Backup old config
mv ~/.config/lumina/keybindings.json ~/.config/lumina/keybindings.json.v1.0.backup

# New version creates new default config
./lumina --version  # v1.1.0
# Auto-creates new format at ~/.config/lumina/keybindings.json

# User can manually migrate custom keybindings
```

## Environment Variables

LUMINA supports environment variables for configuration:

```bash
# Custom config directory
export LUMINA_CONFIG_DIR=~/.lumina

# Disable custom keybindings (use defaults)
export LUMINA_USE_DEFAULT_KEYS=1

# Debug mode
export LUMINA_DEBUG=1
```

## Configuration Validation

### Startup Validation

```go
// Application validates config on load:
func init() {
    // 1. Check config directory exists
    if !configDirExists() {
        createConfigDir()
    }

    // 2. Check keybindings file exists
    if !keybindingsFileExists() {
        createDefaultKeybindings()
    }

    // 3. Validate keybindings format
    kb, err := LoadKeybindings()
    if err != nil {
        log.Warn("Invalid keybindings, using defaults")
        kb = DefaultKeybindings()
    }
}
```

### User Feedback

```bash
$ ./lumina

# First run:
Created default configuration at ~/.config/lumina/keybindings.json
Customize your keybindings by editing this file.

# Subsequent runs:
# (silent - config loaded successfully)

# Invalid config:
Warning: Invalid keybindings format in ~/.config/lumina/keybindings.json
Using default keybindings. Fix config file or delete to regenerate.
```

## Best Practices

### For Developers

1. **Never commit user config files**
   ```bash
   # .gitignore should exclude:
   ~/.config/
   *.local.json
   .env
   ```

2. **Provide clear defaults**
   ```go
   // Always have working defaults in code
   func DefaultKeybindings() *KeyBindings {
       return &KeyBindings{ /* ... */ }
   }
   ```

3. **Validate user input**
   ```go
   // Validate all user config before using
   func (kb *KeyBindings) Validate() error {
       // Check for conflicts, invalid keys, etc.
   }
   ```

4. **Document config format**
   ```bash
   # Include example config in docs/
   docs/config-examples/
   └── keybindings.example.json
   ```

### For Users

1. **Backup your config before upgrading**
   ```bash
   cp ~/.config/lumina/keybindings.json ~/backups/
   ```

2. **Use version control for your dotfiles**
   ```bash
   # Keep user configs in dotfiles repo
   ~/dotfiles/
   └── lumina/
       └── keybindings.json

   # Symlink to config location
   ln -s ~/dotfiles/lumina/keybindings.json ~/.config/lumina/
   ```

3. **Check default config after updates**
   ```bash
   # Regenerate default config to see new options
   mv ~/.config/lumina/keybindings.json{,.old}
   ./lumina --init-config
   diff ~/.config/lumina/keybindings.json{.old,}
   ```

## Troubleshooting

### Config Not Loading

```bash
# Check config file location
ls -la ~/.config/lumina/

# Check file permissions
chmod 644 ~/.config/lumina/keybindings.json

# Validate JSON format
cat ~/.config/lumina/keybindings.json | jq .

# Delete and regenerate
rm ~/.config/lumina/keybindings.json
./lumina --init-config
```

### Conflicting Keybindings

```bash
# Application should detect conflicts on load
# and show helpful error message:

Error: Conflicting keybinding 'j' mapped to both 'down' and 'jump'
Edit ~/.config/lumina/keybindings.json to resolve conflict.
```

### Missing Config Directory

```bash
# Create config directory manually
mkdir -p ~/.config/lumina

# Run with init flag
./lumina --init-config
```

## Future Enhancements

### Planned Config Features

1. **Color Schemes**
   ```json
   {
     "theme": "monokai",
     "colors": {
       "active_pane": "#00D084",
       "inactive_pane": "#666666"
     }
   }
   ```

2. **Custom Actions**
   ```json
   {
     "custom_commands": {
       "edit_in_vim": {
         "key": "e",
         "command": "vim {file}"
       }
     }
   }
   ```

3. **Layout Preferences**
   ```json
   {
     "layout": {
       "file_tree_width": 30,
       "preview_width": 40,
       "show_line_numbers": true
     }
   }
   ```

---

**Last Updated**: 2025-10-21
