Dotstribute: A dotfile manager tool
---

Dotstribute will recursively place symlinks of files in your current working
directory to your home directory. Only the files themselves will be symlinked.
Any sort of directory structure will be recreated as part of a run, but the
directories will not be symlinked.

The dotfiles stored in current working directory do not need to start with a
dot. Dotstribute will add the dot to the symlink when running. This lets you
store them as non-hidden files in a VCS folder. This only applies to the top
level of files and directories. Anything further down the directory tree that
require files start with a dot will need to be stored with dots.

Dotstribute will let you undo the symlinks after they're placed. Run
`dotstribute` in your dotfile directory with the `--unlink` flag to remove them
from your home directory. This currently does not delete the directories
themselves if they only contain the user's dotfiles but I could see that being
a useful feature in the future.

Dotstribute is aware you might not want to symlink everything in a directory. A
`.git` file probably doesn't need to be symlinked to a home directory. You can
add exclusions to a special filename called `.dotignore`. This file should
contain a list of file or directory names you do not want symlink separated by
a new line.

Installation:

```
go install
```

Usage:

```
> dotstribute --help
Usage of dotstribute:
  -a, --ask
  -p, --preview
  -u, --unlink
```


