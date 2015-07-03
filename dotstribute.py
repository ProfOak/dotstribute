#!/usr/bin/env python
import os
from optparse import OptionParser

class Dot():
    def __init__(self, gd):
        self.git_dir = gd

    def exclude(self, dotexclude=".dotexclude"):
        # use the .dotexclude file IN the git directory, not in cwd
        exclude_file = dotexclude
        if self.git_dir != ".":
            exclude_file = self.git_dir + dotexclude

        # prepare list of files not to link to $HOME
        EXCLUDE = []
        if os.path.exists(exclude_file):
            with open(exclude_file) as ef:
                EXCLUDE = ef.read().split()
        # still never add the .dotexclude file to $HOME
        EXCLUDE.append(dotexclude)

        self.files = [f for f in os.listdir(self.git_dir) if f not in EXCLUDE]

    def link(self):
        # by default place dotfiles in your $HOME
        for f in self.files:
            to = os.environ["HOME"] + "/"
            if not f.startswith("."):
                to += "."
            to += f

            # for correct symlinking
            f = os.getcwd() + "/" + self.git_dir + f

            # add option: replace (ask)
            # add option: force replace (no ask)
            if not os.path.exists(to):
                os.symlink(f, to)
            else:
                print "skipping", f

            # add option: chmod symlink

def main():
    parser = OptionParser()
    parser.add_option("-d", "--dotexe", dest = "dot_exclude",
            help = "Exclude files, given by .dotexclude file")
    parser.add_option("-f", "--force", dest = "force", default = False,
            action = "store_true", help = "Force overwrite the previous links")
    parser.add_option("-u", "--unlink", dest = "unlink", default = False,
            action = "store_true", help = "Remove the previous links")

    (options, args) = parser.parse_args()

    if options.dot_exclude:
        dotexclude = options.dot_exclude

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

    d = Dot(git_dir)
    d.exclude(".dotexclude")
    d.link()


if __name__ == "__main__":
    main()

