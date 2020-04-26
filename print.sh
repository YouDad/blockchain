#!/bin/bash

echo blockchain --port 1101 print --group 0 > .a
blockchain --port 1101 print --group 0 2>> .a
echo blockchain --port 1102 print --group 0 >> .a
blockchain --port 1102 print --group 0 2>> .a
echo blockchain --port 1103 print --group 0 >> .a
blockchain --port 1103 print --group 0 2>> .a
echo blockchain --port 1111 print --group 0 >> .a
blockchain --port 1111 print --group 0 2>> .a
echo blockchain --port 2201 print --group 0 >> .a
blockchain --port 2201 print --group 0 2>> .a
echo blockchain --port 2202 print --group 0 >> .a
blockchain --port 2202 print --group 0 2>> .a
echo blockchain --port 2203 print --group 0 >> .a
blockchain --port 2203 print --group 0 2>> .a
echo blockchain --port 2222 print --group 0 >> .a
blockchain --port 2222 print --group 0 2>> .a
echo blockchain --port 1101 print --group 1 >> .a
blockchain --port 1101 print --group 1 2>> .a
echo blockchain --port 1102 print --group 1 >> .a
blockchain --port 1102 print --group 1 2>> .a
echo blockchain --port 1103 print --group 1 >> .a
blockchain --port 1103 print --group 1 2>> .a
echo blockchain --port 1111 print --group 1 >> .a
blockchain --port 1111 print --group 1 2>> .a
echo blockchain --port 2201 print --group 1 >> .a
blockchain --port 2201 print --group 1 2>> .a
echo blockchain --port 2202 print --group 1 >> .a
blockchain --port 2202 print --group 1 2>> .a
echo blockchain --port 2203 print --group 1 >> .a
blockchain --port 2203 print --group 1 2>> .a
echo blockchain --port 2222 print --group 1 >> .a
blockchain --port 2222 print --group 1 2>> .a
vim .a
