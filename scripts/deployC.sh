#!/bin/bash

#mysql Params
my_sql_container_name="mysql-rubis"
my_sql_image_name="aniketdhuri/rubis-mysql"
network_name="my-net"
my_sql_password="aniketd9"

mysql_host_ip="152.46.18.63"
mysql_host_username="root"
mysql_ip_adress=""

#Rubis Params
#sudo docker  run --name rubis --network=my-net -p 81:80 -e "DBIP=10.0.9.2" -e "DBUSER=root" -e "DBPASSWORD=aniketd9" -e "MYHOSTNAME=host1" -d aniketdhuri/rubis-docker
rubis_container_name="rubis"
rubis_image_name="aniketdhuri/rubis-docker"
MYHOSTNAME="host2"

rubis_host_ip="152.46.18.63"
rubis_host_username="root"



log(){ echo [$(date -u)] : $1;}
log_input() { echo -n [$(date -u)] : $1; }

setup_mysql(){
log "Setting up $my_sql_container_name container"

# Get user input for ip
log_input "Enter Host Ip where you want to install $my_sql_container_name [default: $mysql_host_ip]"
read userInput

if [[ -z "$userInput" ]]; then
   a=$mysql_host_ip
else
   # If userInput is not empty show what the user typed in and run ls -l
   mysql_host_ip=$userInput
fi

# Get user input for username
log_input "Enter username for $mysql_host_ip[default :root ]"
read userInput

if [[ -z "$userInput" ]]; then
   mysql_host_username="adhuri"
else
   # If userInput is not empty show what the user typed in and run ls -l
   mysql_host_username=$userInput
fi


log "Stopping existing container [if exists] on host $mysql_host_ip "
log "ssh $mysql_host_username@$mysql_host_ip sudo docker stop $my_sql_container_name"
ssh $mysql_host_username@$mysql_host_ip "sudo docker stop $my_sql_container_name" >/dev/null


log "Removing existing container [if exists] on host $mysql_host_ip "
log "ssh $mysql_host_username@$mysql_host_ip sudo docker rm $my_sql_container_name"
ssh $mysql_host_username@$mysql_host_ip "sudo docker rm $my_sql_container_name" >/dev/null


log "Starting container on host $mysql_host_ip "
log "ssh $mysql_host_username@$mysql_host_ip sudo docker  run --name $my_sql_container_name --network=$network_name -p 3306:3306 -e MYSQL_ROOT_PASSWORD=$my_sql_password  -d $my_sql_image_name"
command="ssh $mysql_host_username@$mysql_host_ip sudo docker  run --name $my_sql_container_name --network=$network_name -p 3306:3306 -e MYSQL_ROOT_PASSWORD=$my_sql_password  -d $my_sql_image_name"
mysql_container_id=`$command`
if [ $? -ne 0 ]
then
    log "docker run failed. exiting..."
    exit $?
else
  log "docker run success"
fi


#Inspecting Ip address

log "Inspecting IP address of started $my_sql_container_name container with containerID $mysql_container_id"
log "ssh $mysql_host_username@$mysql_host_ip sudo docker inspect $mysql_container_id|grep IPv4Address |cut -d \":\" -f2"
mysql_ip_adress=`ssh $mysql_host_username@$mysql_host_ip sudo docker inspect $mysql_container_id|grep IPv4Address |cut -d ":" -f2`
if [ $? -ne 0 ]
then
    log "docker inspect failed. exiting..."
    exit $?
else
  log "docker inspect success - Ip Address of $my_sql_container_name is $mysql_ip_adress"
fi

if [ -z "$mysql_ip_adress" ]
then
      log "mysql_ip_adress is empty"
      exit $?
fi

}

# Setup rubis

setup_rubis(){
#rubis_container_name="rubis"
#rubis_image_name="aniketdhuri/rubis-docker"
#MYHOSTNAME="host2"

#rubis_host_ip="152.46.18.63"
#rubis_host_username="root"


log "Setting up $rubis_container_name container"

# Get user input for ip
log_input "Enter Host Ip where you want to install $rubis_container_name [default: $rubis_host_ip]"
read userInput

if [[ -z "$userInput" ]]; then
   a=$rubis_host_ip
else
   # If userInput is not empty show what the user typed in and run ls -l
   rubis_host_ip=$userInput
fi

# Get user input for username
log_input "Enter username for $rubis_host_ip[default :root ]"
read userInput

if [[ -z "$userInput" ]]; then
   rubis_host_username="adhuri"
else
   # If userInput is not empty show what the user typed in and run ls -l
   rubis_host_username=$userInput
fi


log "Stopping existing container [if exists] on host $rubis_host_ip "
log "ssh $rubis_host_username@$rubis_host_ip sudo docker stop $rubis_container_name"
ssh $rubis_host_username@$rubis_host_ip "sudo docker stop $rubis_container_name" >/dev/null


log "Removing existing container [if exists] on host $rubis_host_ip "
log "ssh $rubis_host_username@$rubis_host_ip sudo docker rm $rubis_container_name"
ssh $rubis_host_username@$rubis_host_ip "sudo docker rm $rubis_container_name" >/dev/null

#rubis_container_name="rubis"
#rubis_image_name="aniketdhuri/rubis-docker"
#MYHOSTNAME="host2"

#rubis_host_ip="152.46.18.63"
#rubis_host_username="root"

log "Starting container on host $rubis_host_ip "
log "ssh $rubis_host_username@$rubis_host_ip sudo docker  run --name $rubis_container_name --network=$network_name -p 81:80 -e DBIP=$mysql_ip_adress -e DBUSER=root -e DBPASSWORD=$my_sql_password -e MYHOSTNAME=$MYHOSTNAME -d $rubis_image_name"
command="ssh $rubis_host_username@$rubis_host_ip sudo docker  run --name $rubis_container_name --network=$network_name -p 81:80 -e DBIP=$mysql_ip_adress -e DBUSER=root -e DBPASSWORD=$my_sql_password -e MYHOSTNAME=$MYHOSTNAME -d $rubis_image_name"
rubis_container_id=`$command`
if [ $? -ne 0 ]
then
    log "docker run failed. exiting..."
    exit $?
else
  log "docker run success"
fi


#Inspecting Ip address

log "Inspecting IP address of started $rubis_container_name container with containerID $rubis_container_id"
log "ssh $rubis_host_username@$rubis_host_ip sudo docker inspect $rubis_container_id|grep IPv4Address |cut -d \":\" -f2"
rubis_ip_adress=`ssh $rubis_host_username@$rubis_host_ip sudo docker inspect $rubis_container_id|grep IPv4Address |cut -d ":" -f2`
if [ $? -ne 0 ]
then
    log "docker inspect failed. exiting..."
    exit $?
else
  log "docker inspect success - Ip Address of $rubis_container_name is $rubis_ip_adress"
fi

if [ -z "$rubis_ip_adress" ]
then
      log "rubis_ip_adress is empty"
      exit $?
fi

}

setup_mysql
setup_rubis
