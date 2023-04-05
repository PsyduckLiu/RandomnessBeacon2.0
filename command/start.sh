#!/bin/bash
cd ../BulletinBoard
sudo go build main.go
sudo go run main.go > result/result.txt &
echo "board is running"
sleep 2
port=40000
PID=$(sudo netstat -nlp | grep "$port" | awk '{print $7}' | awk -F '[ / ]' '{print $1}')
echo ${PID} >> result/running.pid

cd ../collector
sudo go build main.go
for i in $(seq 0 3)
do
sudo go run main.go $i > result/result$i.txt &
echo "consensus node $i is running"
sleep 1
port=3000$i
PID=$(sudo netstat -nlp | grep "$port" | awk '{print $7}' | awk -F '[ / ]' '{print $1}')
echo ${PID} >> result/running.pid
done

cd ../generator
sudo go build main.go
for i in $(seq 0 4)
do
sudo go run main.go $i > result/result$i.txt &
echo "entropy node $i is running"
echo $! >> result/running.pid
sleep 1
done

wait
echo "all nodes are closed"