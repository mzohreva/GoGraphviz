#!/bin/bash

cwd=$(pwd)

cd $cwd/linked-list  && make && make run
cd $cwd/binary-tree  && make && make run
cd $cwd/linked-list2 && make && make run

open $cwd/*/*.png
