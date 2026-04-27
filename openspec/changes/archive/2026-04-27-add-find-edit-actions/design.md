## Context

The gocovtui application is a Go TUI for viewing coverage statistics. Currently, the file-picker component allows users to navigate the file system but lacks advanced features like search/filtering and in-TUI editing capabilities. Users must exit the TUI to edit files with external editors, disrupting workflow.

The existing file-picker likely uses a Go TUI framework (e.g., Bubble Tea, tcell) that handles keybindings and view updates. The application structure suggests modular action handlers that can be extended.

## Goals / Non-Goals

**Goals:**
- Add 'f' keybinding to trigger interactive file search/filter by name and folder name
- Add 'e' keybinding to open selected file in `$EDITOR` environment variable
- Maintain existing file-picker functionality and keybindings
- Ensure search is responsive and filters in real-time or with minimal latency

**Non-Goals:**
- Regex or advanced search syntax (simple substring matching)
- Preview functionality within the TUI
- Multi-file operations
- Modifying or creating files from within the TUI

## Decisions

**1. Find Action Implementation: Real-time Filter vs Dialog**
- **Decision**: Implement as a modal search dialog within the TUI (similar to Vim's `/` search)
- **Rationale**: Maintains focus on the TUI without spawning external tools; provides better UX with live filtering
- **Alternative Rejected**: Shell out to `find` command (less integrated, harder to control)

**2. Edit Action: Shell Execution with EDITOR**
- **Decision**: Use Go's `os/exec` package to spawn `$EDITOR` with the selected file path
- **Rationale**: Standard approach for TUI apps; respects user's editor preference; minimal code
- **Alternative Rejected**: Attempt in-TUI editing (complex, requires buffer management; conflicts with non-goal)

**3. Keybinding Placement**
- **Decision**: Add 'f' and 'e' keybindings to the file-picker's key handler
- **Rationale**: Aligns with existing architecture; only available when a file is selected
- **Open Question**: Do we want 'f' to work without selection to search globally? (Current scope: assumes selection exists)

**4. Filter Storage During Search**
- **Decision**: Maintain a separate filter state that temporarily overrides the displayed file list
- **Rationale**: Non-destructive; users can clear search and return to original view
- **Implementation**: Add `filterActive` and `filterTerm` fields to file-picker state

## Risks / Trade-offs

- **[Risk] Search Performance**: If directory contains thousands of files, real-time filtering may lag
  - *Mitigation*: Implement debouncing (e.g., 100ms) for filter updates; consider lazy loading for large trees
  
- **[Risk] EDITOR Not Set**: If `$EDITOR` is unset, spawning editor fails silently or shows cryptic error
  - *Mitigation*: Validate `EDITOR` at startup or when edit action is triggered; fallback to `vim` or show user-friendly error
  
- **[Risk] Modal Disruption**: Showing search modal may interfere with existing TUI layout
  - *Mitigation*: Use overlay modal design; ensure ESC or Ctrl+C cancels search gracefully

- **[Trade-off] Filter Scope**: Filtering only current directory vs entire tree
  - *Choice*: Start with current directory (simpler, faster); can expand later if needed

## Migration Plan

1. **Phase 1**: Implement find-action modal and filter logic
2. **Phase 2**: Implement edit-action with EDITOR spawning
3. **Phase 3**: Integration testing and keybinding refinement
4. **Phase 4**: Release with documentation noting new shortcuts

Rollback strategy: Remove new keybindings from file-picker handler; revert action functions.

## Open Questions

- Should 'f' search only the current folder or recursively through subdirectories?
- Should search be case-sensitive or case-insensitive?
- How should the filter modal be styled/visually distinct from the main view?
