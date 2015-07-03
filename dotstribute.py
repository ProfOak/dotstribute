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
            ignore_file = self.git_dir + "/" + dotignore

        # prepare list of files NOT to link to $HOME
        IGNORE = []
        if os.path.exists(ignore_file):
            with open(ignore_file) as f:
                IGNORE = f.read().split()
        # still never add the .dotignore file to $HOME
        IGNORE.append(dotignore)

        # compile a list of full paths to begin with
        # this will make other operations much simpler
        self.git_links  = []
        self.home_links = []

        for f in os.listdir(self.git_dir):
            if f in IGNORE:
                continue

            # by default place dotfiles in your $HOME
            to = os.environ["HOME"] + "/"
            if not f.startswith("."):
                to += "."
            to += f

            # for correct symlinking
            # of the form /this/path/to/file
            tmp_dir = self.git_dir
            if tmp_dir == ".":
                tmp_dir = ""
            # don't trust user input
            if not tmp_dir.endswith("/"):
                tmp_dir += "/"

            f = os.getcwd() + "/" + tmp_dir + f
            self.git_links.append(f)
            self.home_links.append(to)

    def link(self):
        for i, f in enumerate(self.git_links):
            # add option: replace (ask)
            # add option: force replace (no ask)
            if not os.path.exists(self.home_links[i]):
                os.symlink(f, self.home_links[i])
            else:
                print "skipping", f

    def unlink(self):
        for f in self.home_links:
            if os.path.exists(f):
                os.unlink(f)
            else:
                print "Does not exist:", f

    # add option: chmod symlinked files

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
    if options.unlink:
        d.unlink()
    else:
        d.link()


if __name__ == "__main__":
    main()

