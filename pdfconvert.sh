#!/bin/sh
echo "running chrome headless"

alias chrome="/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome"

chrome --headless --disable-gpu --print-to-pdf input.txt
ls

