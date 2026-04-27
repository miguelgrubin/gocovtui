## Why

The TUI coverage viewer currently lacks efficient file navigation and editing capabilities. Users must rely on external tools to find files and edit them, breaking their workflow. Adding integrated find and edit actions will improve user productivity by keeping users within the TUI interface for common file operations.

## What Changes

- Add a **Find action** triggered by pressing 'f' to interactively search and filter files by name and folder name from the current directory
- Add an **Edit action** triggered by pressing 'e' to open the selected file in the configured `EDITOR` environment variable
- Integrate both actions with the existing file selection context

## Capabilities

### New Capabilities
- `find-action`: Interactive file search and filter functionality triggered by 'f' key, allowing users to search files and folders by name
- `edit-action`: Open selected file in external editor using `EDITOR` env var, triggered by 'e' key

### Modified Capabilities
- `file-picker`: Extend file picker to support new keybindings ('e' for edit, 'f' for find) alongside existing navigation

## Impact

- **Code**: Modifications to the file-picker component/handler to add keybinding logic and action handlers
- **UX**: New keyboard shortcuts ('e', 'f') available during file browsing
- **Dependencies**: Will rely on environment configuration for the `EDITOR` variable; existing Go standard library exec package
- **Systems**: File browsing and selection workflow enhanced with new capabilities
