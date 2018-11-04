#! /bin/bash
# set the STRING value
echo Deploy festa junina ThisAPIPort : 1662
echo on
STRING="Deploy by copying files"
echo $STRING
cp ./fjapisecurity ~/golang/runtime/festajunina/fjapisecurity
cp ./fjapisecurity.ini ~/golang/runtime/festajunina/fjapisecurity

