#! /bin/bash
# set the STRING value
echo Deploy restaurante ThisAPIPort : 1661
echo on
STRING="Deploy by copying files"
echo $STRING
cp ./fjapisecurity ~/golang/runtime/restaurante/restapisecurity
cp ./fjapisecurity.ini ~/golang/runtime/restaurante/restapisecurity

