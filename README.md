Dotstribute
===

###Keep your original dotfiles into a git directory, and symlink them to your home directory.

This script uses an optional file called the ".dotignore" file. Create this .dotignore file in the directory with all of your dotfiles. The purpose of this file is to specify other files that will not be symlinked to your home directory. Think of it like a .gitignore file, for dotsribute. This repository comes with a .dotignore file. Please use that as an example.

dotstribute takes an optional argument that should specify the folder your dotfiles are contained in. dotstribute will symlink the contents of that folder. If no argument is provided, dotstribute will use the current directory.


Instalation
===

`$ python setyp.py install`

Flags explained
===

####-d / --dotignore
Specify a different location for your .dotignore file

example:

`./dotstribute.py -d /home/user/I/put/my/file/somewhere/else/.dotignore`

####-u / --unlink
Unlink (remove) the links in your home directory, based on the contents of your dotfiles folder. See example image.

Example:

`./dotstribute.py -u`

####-a / --ask
Ask before linking or unlinking each file.

Example:

```
./dotstribute.py -a
./dotstribute.py -au
```

####-p / --preview
Preview the changes before they happen.

```
./dotstribute.py -p
./dotstribute.py -pu
```

####-g / --generate-dotignore
Generate a .dotignore file containing generic files you may find in a git directory, that you may not want to symlink.

```
./dotstribute -g /path/to/dotfiles
```

![Example usage](https://raw.githubusercontent.com/ProfOak/dotstribute/master/media/example.png)
