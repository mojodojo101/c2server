# c2server
## A simple command and control server written in golang
---
#### The core is working now, together with the angular client
#### and the [testbeacon.go](https://github.com/mojodojo101/c2server/tree/master/internal_resources/beacons) you should be able to get a demo going on a local server 
#### the token is just a stub and this will change later on
#### same with the entire authentication at the moment

TODO:
* Add better authentication for client and active beacon
* Add Beacon Templates
* Add better Server-Beacon Communication
* Add better Client-Server Communication



#### Setting up the database

* get postgresql repo

```sh
sudo apt update
sudo apt install postgresql postgresql-contrib
```
* switch user to postgresql and create roles and database
```sh

sudo su postgres

createuser --interactive 
#----->add the “c2admin” user


createdb c2db
```

* set up priviliges and password
```sh
psql
grant all privileges on database c2db to c2admin;
alter user c2admin password 'mojodojo101+';
```

* make sure your pg_hba.conf allows the correct connection types
```
# Database administrative login by Unix domain socket
local   all             postgres                                peer
# TYPE  DATABASE        USER            ADDRESS                 METHOD

# "local" is for Unix domain socket connections only
local   all             all             	                peer
# IPv4 local connections:
host    all             all             127.0.0.1/32            md5
# IPv6 local connections:
host    all             all             ::1/128                 md5
# Allow replication connections from localhost, by a user with the
# replication privilege.
local   replication     all                                     peer
host    replication     all             127.0.0.1/32            md5
host    replication     all             ::1/128                 md5

```

```sh

systemctl restart postgresql.service

```

if you used different accounts names and passwords change the [config.go](https://github.com/mojodojo101/c2server/blob/master/config/config.go)
