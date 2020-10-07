
/// To build the images, navigate to "GolangEmail" folder in the command prompt and do the following commands

docker build -t msaimage ./msa
docker build -t mtaimage ./mta
docker build -t bluebookimage ./bluebook

/// Create the subnet
docker network create --subnet 192.168.1.0/24 emailservice

/// To set up the containers do the following
docker run --name msacontainer --net emailservice --ip 192.168.1.8 --detach --publish 3000:8888 --security-opt apparmor=unconfined msaimage
docker run --name mtacontainer --net emailservice --ip 192.168.1.7 --detach --publish 3001:8989 --security-opt apparmor=unconfined mtaimage
docker run --name bluebookcontainer --net emailservice --ip 192.168.1.6 --detach --publish 3002:9090 --security-opt apparmor=unconfined bluebookimage



//The two email addresses in the bluebook are "Billy" and "Bob" these two can use the email service with one another curling the MSA

///Example commands you can use are: 
//please note that for me the ip is 192.168.99.100, you can check yours by doing the command "docker-machine env NAME" where name is the name of the machine.

// To send a email from Billy to Bob
curl -v -X POST -d "{\"To\":\"Bob\", \"Body\":\"Body of the Message\" }" 192.168.99.100:3000/MSA/Outbox/Billy

// To List Bobs Inbox
curl -v GET 192.168.99.100:3000/MSA/Inbox/Bob

// To List Billys Outbox
curl -v GET 192.168.99.100:3000/MSA/Outbox/Billy

// To read a specific email in Bobs inbox (where X is an int value)
curl -v GET 192.168.99.100:3000/MSA/Inbox/Bob/X

// To delete a specific email in Bobs inbox (where X is an int value)
curl -v -X DELETE 192.168.99.100:3000/MSA/Inbox/Bob/X

