# wingbaas
Blockchain as a service

# clone project
git clone https://github.com/wingchain/wingbaas.git

# platformsrv build and run prerequest 
1, install golang ,version >=1.14.2
2, config GOPATH enviroment
3, clone project to $GOPATH/src/github.com
4, install nfs client
nfs搭建
1）在NFS服务器上安装NFS
[root@youxi1 ~]# yum -y install rpcbind nfs-utils
2）启动NFS，并开机自启
[root@youxi1 ~]# systemctl start rpcbind
[root@youxi1 ~]# systemctl enable rpcbind
[root@youxi1 ~]# systemctl start nfs-server　　//NFS依赖rpcbind进行通讯，所以要先启动rpcbind
[root@youxi1 ~]# systemctl enable nfs-server
[root@youxi1 ~]# netstat -antup | grep 2049
netstat -anp|grep 2049


# platformsrv build and run
1, cd  $GOPATH/src/github.com/wingbaas/platformsrv
2, go build
3, ./platformsrv