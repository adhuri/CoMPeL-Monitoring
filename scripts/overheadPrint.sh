#!/bin/bash

log(){ echo [$(date -u)] : $1;}

#Global vars

sum_cpu_monitoring_agent=0
sum_memory_monitoring_agent=0

sum_cpu_migration_agent=0
sum_memory_migration_agent=0

counter=1
#localVars

monitoring_agent="compel-monitoring-agent"
migration_agent="Compel-Migration-agent"


while :
do
clear

log "Checking stats of  $monitoring_agent and $migration_agent"

#Monitoring agent
stats_monitoring_agent=`ps aux|grep "$monitoring_agent" | grep -v grep |grep -v sudo |head -n 1 `

pid_monitoring_agent=`echo $stats_monitoring_agent|awk '{print $2}'`
cpu_monitoring_agent=`echo $stats_monitoring_agent|awk '{print $3}'`
memory_monitoring_agent=`echo $stats_monitoring_agent|awk '{print $4}'`
sum_cpu_monitoring_agent=`echo "$sum_cpu_monitoring_agent + $cpu_monitoring_agent" |bc`
sum_memory_monitoring_agent=`echo "$sum_memory_monitoring_agent + $memory_monitoring_agent" |bc`


#Migration agent
stats_migration_agent=`ps aux|grep "$migration_agent" | grep -v grep |grep -v sudo |head -n 1 `

pid_migration_agent=`echo $stats_migration_agent|awk '{print $2}'`
cpu_migration_agent=`echo $stats_migration_agent|awk '{print $3}'`
memory_migration_agent=`echo $stats_migration_agent|awk '{print $4}'`

sum_cpu_migration_agent=`echo "$sum_cpu_migration_agent + $cpu_migration_agent" |bc`
sum_memory_migration_agent=`echo "$sum_memory_migration_agent + $memory_migration_agent" |bc`



if [[ -z "$stats_monitoring_agent" ]]; then
   log "[ERROR] Stats for $monitoring_agent not found . Check if $monitoring_agent is running"
  
else
   # Found PID of compel-monitoring-monitoring_agent
   echo "[$(date -u)] [INFO] Stats for $monitoring_agent found. | PID : $pid_monitoring_agent | CPU : $cpu_monitoring_agent | Memory : $memory_monitoring_agent |" >> /tmp/`basename "$0"`.tmp
fi


if [[ -z "$stats_migration_agent" ]]; then
   log "[ERROR] Stats for $migration_agent not found . Check if $migration_agent is running"
else
   # Found PID of compel-monitoring-monitoring_agent

   echo "[$(date -u)] [INFO] Stats for $migration_agent found. | PID : $pid_migration_agent | CPU : $cpu_migration_agent | Memory : $memory_migration_agent |" >> /tmp/`basename "$0"`.tmp
fi

log "-----------------------------------"
log "Average CPU% Utilization of $monitoring_agent `echo "scale=2; $sum_cpu_monitoring_agent/$counter"|bc`%"
log "Average MEM% Utilization of $monitoring_agent `echo "scale=2; $sum_memory_monitoring_agent/$counter"|bc`%"
log " "

log "Average CPU% Utilization of $migration_agent `echo "scale=2; $sum_cpu_migration_agent/$counter"|bc`%"
log "Average MEM% Utilization of $migration_agent `echo "scale=2; $sum_memory_migration_agent/$counter"|bc`%"

(( counter++ ))


sleep 5
done

