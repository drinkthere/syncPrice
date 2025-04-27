#!/bin/bash

cpuidx=0

start() {
    echo "Starting the process..."
    nohup  taskset -c "$cpuidx" ./syncprice ../config/config.json >> /data/dc/syncPrice/nohup.log 2>&1 &
    echo "Process started."
}

stop() {
    echo "Stopping the process..."
    pid=$(pgrep -f "syncprice ../config/config.json")
    if [ -n "$pid" ]; then
        kill -SIGINT $pid
        echo "Process stopping: "
        sleep 1
        echo "Process stopped."
    else
        echo "Process is not running."
    fi
}

restart() {
    stop
    sleep 10
    start
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
        ;;
esac

exit 0