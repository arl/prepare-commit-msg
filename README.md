# prepare-commit-msg
Git prepare-commit-msg hook that prepends the commit message with the file or directory modified in the commit

## Installation

```
# Build and install the binary on your system
go install github.com/arl/prepare-commit-msg@latest

# Install it as git hook
ln -s $(which prepare-commit-msg) .git/hooks/prepare-commit-msg
```

## Examples

If git status shows:
```
On branch main
Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
        new file:   dir1/file1
        new file:   dir1/file2
```

Then the git commit message will be prepared with:
```
dir1: 
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Changes to be committed:
#	new file:   dir1/file1
#	new file:   dir1/file2
```

If git status shows:
```
On branch main
Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
        modified:   README.md
```

Then the git commit message will be prepared with:
```
README.md: 
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# On branch main
# Changes to be committed:
#	modified:   README.md
```

If there's no common denominator between staged files, then the commit message will be prepended with `all: `
