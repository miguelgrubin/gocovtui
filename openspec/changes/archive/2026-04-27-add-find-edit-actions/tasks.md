## 1. Setup and Investigation

- [x] 1.1 Review existing file-picker implementation to understand current keybinding architecture
- [x] 1.2 Identify file-picker component/module location and its state management
- [x] 1.3 Examine how existing keybindings (e.g., navigation arrows, Enter) are handled
- [x] 1.4 Understand the TUI framework being used (Bubble Tea, tcell, etc.)
- [x] 1.5 Create a new file for find/edit action handlers (e.g., `pkg/actions/find.go`, `pkg/actions/edit.go`)

## 2. Find Action - Modal and UI

- [x] 2.1 Design find mode modal component with search input field (e.g., "Find: ")
- [x] 2.2 Add find mode state to file-picker (e.g., `findModeActive bool`, `searchTerm string`, `filteredFiles []File`)
- [x] 2.3 Implement keybinding for 'f' key to activate find mode in file-picker key handler
- [x] 2.4 Create rendering logic to display find modal overlaid on file picker
- [x] 2.5 Implement ESC/Ctrl+C handling to exit find mode without selecting

## 3. Find Action - Filter Logic

- [x] 3.1 Implement substring matching function (case-insensitive) to filter files by name
- [x] 3.2 Add debouncing (100-200ms) to filter updates for performance with large directories
- [x] 3.3 Implement real-time filtering as user types in search input
- [x] 3.4 Handle updating filtered file list and selection highlighting during search
- [x] 3.5 Implement logic to clear search and restore full file list
- [x] 3.6 Implement navigation through filtered results using arrow keys

## 4. Find Action - Selection and Confirmation

- [x] 4.1 Implement Enter key handling to select and confirm file from filtered results
- [x] 4.2 Implement logic to close find mode after selection and return to normal view
- [x] 4.3 Add deletion of characters in search input (backspace handling)
- [x] 4.4 Test find mode exit behavior with various search states

## 5. Edit Action - EDITOR Validation and Setup

- [x] 5.1 Add EDITOR environment variable validation at application startup
- [x] 5.2 Implement fallback logic for missing EDITOR (default to "vim" or show user error)
- [x] 5.3 Add user-friendly error messaging if EDITOR command is invalid or missing
- [x] 5.4 Test EDITOR validation with various configurations

## 6. Edit Action - Implementation

- [x] 6.1 Implement 'e' keybinding in file-picker key handler
- [x] 6.2 Add pre-edit validation (file must be selected, file must exist)
- [x] 6.3 Implement spawning external editor process using `os/exec` with selected file path
- [x] 6.4 Handle TUI suspension while editor is active (ensure TUI is not responsive during edit)
- [x] 6.5 Implement logic to resume TUI after editor process exits
- [x] 6.6 Add error handling and display for failed editor spawning

## 7. Integration and Testing

- [x] 7.1 Integration test: 'f' keybinding activates find mode
- [x] 7.2 Integration test: Find mode filters files correctly with various search terms
- [x] 7.3 Integration test: Find mode keyboard navigation works (arrows, backspace, ESC, Enter)
- [x] 7.4 Integration test: 'e' keybinding opens selected file in EDITOR
- [x] 7.5 Integration test: TUI resumes correctly after exiting editor
- [x] 7.6 Edge case testing: Empty directory, very large directory, special characters in filenames
- [x] 7.7 Test EDITOR fallback behavior when EDITOR variable is unset

## 8. Documentation and User Guidance

- [x] 8.1 Update application help/usage text to document 'f' and 'e' keybindings
- [x] 8.2 Add inline help in find mode showing available commands (ESC to cancel, Enter to select, etc.)
- [ ] 8.3 Create or update README with new keyboard shortcuts
- [x] 8.4 Add code comments documenting find/edit action implementation

## 9. Final Verification and Release

- [x] 9.1 Manual end-to-end testing of find action with various file trees
- [x] 9.2 Manual end-to-end testing of edit action with different EDITOR values
- [x] 9.3 Verify no regression in existing file-picker functionality
- [ ] 9.4 Performance testing with large directory structures (find latency)
- [ ] 9.5 Cross-platform testing (Linux, macOS, Windows compatibility)
- [ ] 9.6 Code review of changes
- [ ] 9.7 Merge and deploy
