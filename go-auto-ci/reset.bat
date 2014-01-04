@echo off
plink -ssh -pw * root@192.168.192.81 "echo 0 > /var/task/%1"
