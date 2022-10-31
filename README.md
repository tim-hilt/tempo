# Book

A cli to manage bookings with the Tempo-API

## Requirements

### General

- CLI commands and arguments
- Logging

### Functional

| Requirement                                                          | Priority (0-5) |
| -------------------------------------------------------------------- | -------------- |
| Create worklogs for a specific day by parsing daily note             | 5              |
| Watch directory of daily notes and update debounced (5min) on change | 4              |

### Libraries to help with achieving the requirements

| Requirement         | Package  |
| ------------------- | -------- |
| CLI parsing         | cobra    |
| Logging             | zerolog  |
| Accessing Tempo API | resty    |
| File Watching       | fsnotify |

## Usage

- `book day 2022-10-30` submits all bookings on the given day
- `book watch` watches for file changes in the daily-notes-directory and synchronizes the changes to Tempo
