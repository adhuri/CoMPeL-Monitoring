## To create swarm on master,host1,host2


egIP - 
master = 152.7.99.151,10.10.7.151
host1 = 152.1.13.166,10.10.3.166
host2 = 152.1.13.94,10.10.3.94


#On master init swarm - prefer public ip

sudo docker swarm init --advertise-addr 152.7.99.151

#Check status
sudo docker node ls

#On host1 and host2

#Find manager tocken
sudo docker swarm join-token manager
#This is important
docker swarm join  --token <token>152.7.99.151:2377

#Start network on master

sudo docker network create  --driver overlay  --subnet 10.0.9.0/24  my-net --attachable


# Promote host1 and host2 as manager so that you see MANAGER STATUS= in ```docker node ls```
sudo docker node promote <nodeid from node ls >



#Check on host1 and host2 if network my-net is available

sudo docker network ls
