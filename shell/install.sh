#!/bin/sh

#mongodb
wget https://repo.mongodb.org/yum/redhat/7/mongodb-org/4.4/x86_64/RPMS/mongodb-org-mongos-4.4.14-1.el7.x86_64.rpm
wget https://repo.mongodb.org/yum/redhat/7/mongodb-org/4.4/x86_64/RPMS/mongodb-org-server-4.4.14-1.el7.x86_64.rpm
wget https://repo.mongodb.org/yum/redhat/7/mongodb-org/4.4/x86_64/RPMS/mongodb-org-shell-4.4.14-1.el7.x86_64.rpm
wget https://fastdl.mongodb.org/tools/db/mongodb-database-tools-rhel70-x86_64-100.5.2.rpm
rpm -ivh mongodb-org-mongos-4.4.14-1.el7.x86_64.rpm
rpm -ivh mongodb-org-server-4.4.14-1.el7.x86_64.rpm
rpm -ivh mongodb-org-shell-4.4.14-1.el7.x86_64.rpm
rpm -ivh mongodb-database-tools-rhel70-x86_64-100.5.2.rpm  --force --nodeps

systemctl enable mongod
systemctl start mongod