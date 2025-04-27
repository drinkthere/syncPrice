#!/bin/bash

# 指定日志文件所在的目录
spot_log_directory="/data/dc/syncPrice"
+ls -t $spot_log_directory/syncprice.log* | tail -n +5 | xargs rm -f