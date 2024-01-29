#!/bin/bash
printf '[\n'
for i in {1..30};do
    printf '    {
        "id": '${i}',
        "hostID": "host-1-'${i}'",
        "hostName": "mysql-'${i}'.dev.com",
        "hostIP": "192.168.1.'${i}'",
        "hostSSHPort": 22,
        "hostType": "裸金属",
        "createdTime": "2023-12-'${i}' 15:30:00",
        "updatedTime": "2023-12-'${i}' 15:30:00",
        "status": "在线"
    }'
    if [ ${i} -ne 30 ];then
        printf ','
    fi
    printf '\n'
done

printf ']\n'
