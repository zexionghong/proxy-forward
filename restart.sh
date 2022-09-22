#!/bin/bash

#***************************************
# Auth: John
# Email: zy1234500@outlook.com
# Ver: 0.1
# Date: 2022.9.22
# Args:
# Description:   Rstart the proxy-forward program and set nofile 
#***************************************


# 1 check supervisor
which supervisorctl >> /dev/null 2>&1
if [ $? != 0 ]
then 
    echo '==============================='
    echo '[ERROR] supervisorctl not found'
    echo '==============================='
    exit 2
fi 


# 2 restart http-proxy-forward
echo '[INFO] do restart http-proxy-forward'
supervisorctl restart http-proxy-forward:
if [ $? != 0 ]
then
    echo '========================================='
    echo '[ERROR] restart http-proxy-forward failed'
    echo '========================================='
    exit 3
fi


# 3 restart socks-proxy-forward
echo '[INFO] do restart socks-proxy-forward'
supervisorctl restart socks-proxy-forward:
if [ $? != 0 ]
then
    echo '========================================='
    echo '[ERROR] restart http-proxy-forward failed'
    echo '========================================='
    exit 3
fi


# 4 get socks-proxy-forward pid
pids=` supervisorctl  status | grep socks-proxy-forward | grep -v grep | awk '{print $4}' `    
echo '[DEBUG] socks-proxy-forward pids:' $pids


# 5 set prlimit nofile 655350
for p in $pids
do
    pid=${p/','/''}
    echo '[INFO] do set prlimit nofile pid:' $pid
    prlimit --pid $pid --nofile=655350
    if [ $? != 0 ]
    then
        echo '================================='
        echo '[ERROR] set prlimit nofile failed'
        echo '================================='
        exit 3
    fi
    cat /proc/$pid/limits | grep 'Max open files'
    echo '' 
done

# END 
echo '[SUCCEEC] done'