#!/bin/bash
cd ../BulletinBoard
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid

cd ../collector
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid

cd ../generator
cat result/running.pid | xargs -IX kill -9 X
:> result/running.pid

for i in $(seq 0 9)
do
kill -9 `ps -ef |grep exe/main\ $i|awk '{print $2}'`
done

ps -ef | grep start.sh | grep -v grep | awk '{print $2}' | xargs kill