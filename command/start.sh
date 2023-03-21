#!/bin/bash
cd ../collector
for i in $(seq 0 3)
do
go run main.go $i > result/result$i.txt &
echo "consensus node $i is running"
sleep 1
port=3000$i
PID=$(sudo netstat -nlp | grep "$port" | awk '{print $7}' | awk -F '[ / ]' '{print $1}')
echo ${PID} >> result/running.pid
done

cd ../generator
for i in $(seq 0 9)
do
go run main.go $i > result/result$i.txt &
echo "entropy node $i is running"
echo $! >> result/running.pid
sleep 1
done

sleep 5

cd ../BackBone
go run main.go 

wait
echo "all nodes are closed"