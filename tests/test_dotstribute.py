import os
import shutil
import tempfile
from unittest import TestCase

from dotstribute import (
    generate_dotignore,
    get_files,
    # link,
    # preview,
    # unlink,
)


class DotstributeTest(TestCase):
    def setUp(self):
        self.home_path = tempfile.mkdtemp()
        self.git_path = tempfile.mkdtemp()

        ignore_files = [
            'ignore1',
            'ignore2',
        ]

        dotfiles = [
            'one',
            '.two',
        ]

        # create dotignore file
        with open(self.git_path + '/.dotignore', 'w') as f:
            for i in ignore_files:
                f.write(i + '\n')

        # create git directory
        for dot in dotfiles:
            with open(os.path.join(self.git_path, dot), 'w') as f:
                f.write('')

        # make a directory with both dotfiles and ignore files in them
        final_dotfiles = []
        for d in dotfiles + ignore_files:
            if d.startswith('.'):
                final_dotfiles.append(d)
            else:
                final_dotfiles.append('.' + d)


        self.home_paths = [os.path.join(self.home_path, d) for d in final_dotfiles]
        self.git_paths = [os.path.join(self.git_path, d) for d in dotfiles]

    def tearDown(self):
        shutil.rmtree(self.home_path)
        shutil.rmtree(self.git_path)

    def test_get_files(self):
        symlinks = get_files(self.git_path, self.home_path)
        for git_path, home_path in symlinks.items():
            assert git_path in self.git_paths
            assert home_path in self.home_paths

