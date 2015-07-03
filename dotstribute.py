#!/usr/bin/env python
import os
from optparse import OptionParser

class Dot():
    def __init__(self, git_dir):
        self.git_dir = git_dir

    def get_files(self, dotignore=".dotignore"):
        # use the .dotignore file IN the git directory, not in cwd
        ignore_file = dotignore
        if self.git_dir != ".":
            ignore__file = self.git_dir + dotignore

        # prepare list of files not to link to $HOME
        IGNORE = []
        if os.path.exists(ignore_file):
            with open(ignore_file) as ef:
                IGNORE = ef.read().split()
        # still never add the .dotignore file to $HOME
        IGNORE.append(dotignore)

        self.files = [f for f in os.listdir(self.git_dir) if f not in IGNORE]

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
    parser.add_option("-d", "--dotignore", dest = "dot_ignore",
            help = "Exclude files, given by .dotignore file")
    parser.add_option("-f", "--force", dest = "force", default = False,
            action = "store_true", help = "Force overwrite the previous links")
    parser.add_option("-u", "--unlink", dest = "unlink", default = False,
            action = "store_true", help = "Remove the previous links")

    (options, args) = parser.parse_args()

    if options.dot_ignore:
        dotignore = options.dot_ignore

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
    d.get_files(".dotignore")
    d.link()


if __name__ == "__main__":
    main()

