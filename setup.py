#!/usr/bin/env python

import os
from setuptools import setup, find_packages

setup(
    name='Dotstribute',
    version='1.1',
    description='Make your dotfiles easier to manage',
    author='ProfOak',
    author_email='OpenProfOak@gmail.com',
    url='https://github.com/ProfOak/dotstribute',
    packages=find_packages(),
    scripts=[
        'dotstribute.py',
    ],
    keywords=['dotfile github version control'],
)
