# K8S watcher skeleton

### Create user and get token
    oc create serviceaccount watcher
    oadm policy add-cluster-role-to-user cluster-admin system:serviceaccount:default:watcher
   

 If you want find framework libary, copy it from github, instead of k8s.io kubernetes project
 
### build:
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

### flag:
    --host - hostname of master (with port, without api prefix)
    --token - token. TODO: insect why config does not take config from config
    --username - username
    --password - password

    You can use token OR Username/password. 