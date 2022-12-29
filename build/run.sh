#!/bin/bash
chmod +x /data/abc
while sleep 10; do
  ps -ef | grep "abc" | grep -v "grep"
  if [ "$?" -eq 1 ]; then
    sleep 3
    cd /data
    nohup ./abc >>/dev/null 2>&1 &
  fi
done