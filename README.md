**twkb** (short for **t**ask**w**arrior **k**an**b**an) is a kanban-style TUI for [taskwarrior](https://taskwarrior.org), written in Go.

## Features

- Kanban view with `todo`, `doing` and `done columns`
- Creation of new tasks
- Modifying existing tasks
- Block and unblock tasks
- Create recurring tasks
- Delete tasks

In development:

- [ ] Project tabs
- [ ] Complete info of a single task

## Installation

> [!NOTE]
> twkb needs [taskwarrior](https://taskwarrior.org) to be installed

Currently, the way to install **twkb** is to use `go install`:

```sh copy
go install github.com/DerTimonius/twkb
```

## Usage

| **Keys**         | **View**                    | **Functionality**                                    |
| ---------------- | --------------------------- | ---------------------------------------------------- |
| `Arrows`, `hjkl` | `normal`, `block form`      | Navigation in the different columns                  |
| `Space`          | `normal`                    | Start/Stop selected task                             |
| `Enter`          | `normal`                    | Finish selected task                                 |
| `n`              | `normal`                    | Create new task, enters `create form`                |
| `m`              | `normal`                    | Modify selected task, enters prefilled `create form` |
| `d`              | `normal`                    | Delete selected task, enters `confirmation screen`   |
| `b`              | `normal`                    | Block other tasks selected task, enters `block form` |
| `u`              | `normal`                    | Unblock selected task, enters `confirmation screen`  |
| `Tab`            | `create form`               | Go to next field                                     |
| `Enter`          | `create form`, `block form` | Confirm the form / selection                         |
| `Space`          | `block form`                | Select task that should be blocked                   |
| `Esc`            | `all forms`                 | Go back to `normal` view                             |
| `y`              | `confirmation screen`       | Confirm                                              |

## Contributing

Contributions are always welcome! Please open an issue or submit a pull request if you have any improvements, bug fixes, or new features to propose.

## License

twkb is licensed under the MIT License.
