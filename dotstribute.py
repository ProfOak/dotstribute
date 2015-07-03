#!/usr/bin/env python
import os
from optparse import OptionParser

def main():
    parser = OptionParser()
    parser.add_option("-d", "--dotexe", dest = "dot_exclude",
            help = "Exclude files, given by .dotexclude file")
    parser.add_option("-f", "--force", dest = "force", default = False,
            action = "store_true", help = "Force overwrite the previous links")

    (options, args) = parser.parse_args()

    # option for custom dotexclude file
    dotexclude = ".dotexclude"
    if options.dot_exclude:
        dotexclude = options.dot_exclude

    # prepare list of files not to link to $HOME
    EXCLUDE = []
    if os.path.exists(dotexclude):
        with open(dotexclude) as f:
            EXCLUDE = f.read().split()
    # still never add the .dotexclude file to $HOME
    EXCLUDE.append(dotexclude)

    # get files not in exclude list
    git_dir = "."
    if len(args) == 1:
        if os.path.exists(args[0]):
            git_dir = args[0]
        elif not os.path.exists(args[0]):
            print "That directory does not exist"
            print "Use the current working directory? (y/N)"
            if raw_input("> ").lower() != "y":
                print "Now exiting"
                return

    files = [f for f in os.listdir(git_dir) if f not in EXCLUDE]

    # by default place dotfiles in your $HOME
    for f in files:
        to = os.environ["HOME"] + "/"
        if not f.startswith("."):
            to += "."
        to += f

        # for correct symlinking
        f = os.getcwd() + "/" + git_dir + f

        # add option: replace (ask)
        # add option: force replace (no ask)
        if not os.path.exists(to):
            os.symlink(f, to)
        else:
            print "skipping", f

        # add option: chmod symlink


if __name__ == "__main__":
    main()

