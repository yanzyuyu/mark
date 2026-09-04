# mark ◈

Version-controlled code bookmarks. Store bookmarks in your git repo, share them with your team, open them in any editor.

## Install

```sh
go install github.com/yanzyuyu/mark/cmd/mark@latest
```

Or download a binary from [Releases](https://github.com/yanzyuyu/mark/releases).

## Usage

```sh
mark add "entry point"  cmd/main.go:1
mark add "hot path"     pkg/vm/vm.go:142
mark add "auth logic"   internal/auth/handler.go:89

mark ls
mark go "hot path"
mark rm "entry point"
mark check
```

## Commands

| Command | Description |
|---|---|
| `mark add <label> <file>:<line>` | Add a bookmark |
| `mark ls` | List all bookmarks |
| `mark go <label>` | Open bookmark in `$EDITOR` |
| `mark rm <label>` | Remove a bookmark |
| `mark check` | Verify bookmarks are still valid |

## How It Works

`mark` stores all bookmarks in a `.marks` file at your git repo root:

```json
{
  "version": 1,
  "bookmarks": [
    {
      "id": "a3f7b2c1",
      "label": "hot path",
      "file": "pkg/vm/vm.go",
      "line": 142,
      "content": "func (vm *VM) Run() error {",
      "created_at": "2026-09-04T11:30:00Z",
      "author": "yanzyuyu"
    }
  ]
}
```

Commit `.marks` to share bookmarks with your team. Every editor, every machine, same bookmarks.

## Editor Support

Set `$EDITOR` to your preferred editor:

```sh
export EDITOR=vim      # vim +<line> <file>
export EDITOR=nano     # nano +<line> <file>
export EDITOR=code     # code -g <file>:<line>
export EDITOR=subl     # subl <file>:<line>
```

If `$EDITOR` is unset, `mark go` prints the path:line to stdout.

## mark check

Detects when bookmarked lines have been edited or the file has changed:

```
  ✓  hot path
  ~  auth logic — content changed at line 89
     expected: func handleLogin(w http.ResponseWriter, r *http.Request) {
     current:  func handleLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
  ✕  entry point — file not found or too short

  2/3 bookmarks valid
```
