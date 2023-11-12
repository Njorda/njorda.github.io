---
layout: post
title: "How to set up Neon, serverless postgres on k8s"
subtitle: "Run serverless postgres"
date: 2023-11-12
author: "Niklas Hansson"
URL: "/2023/11/12/serverless-postgres"
---


All the code can be found [here](https://github.com/Njorda/neon-setup)

In this blog post we will dive in to how to set up [neon](https://github.com/neondatabase/neon) on your k8s cluster. We will use [Minikube](https://minikube.sigs.k8s.io/docs/start/) but feel free to use the setup of your choice of k8s. The first step is to define the k8s resources. In this case we will take a short cut and start of with generating them from the [docker compose files](https://github.com/neondatabase/neon/blob/release-4179/docker-compose/docker-compose.yml) used for testing. For this we will use [kompsoe](https://kubernetes.io/docs/tasks/configure-pod-container/translate-compose-kubernetes/) from the compose directory: 

```bash
kompose --file docker-compose.yml convert
```


This is a great start but sadly it will not be a one hit wonder. First of all we need to understand a bit more what the docker compose set up is actually doing(in order to understand how neon works) and specifically the scripts inside [/docker-compose/compute_wrapper](https://github.com/neondatabase/neon/tree/release-4179/docker-compose/compute_wrapper). The script `shell/compute.sh` tells most of the story: 


```
#!/bin/bash
set -eux

# Generate a random tenant or timeline ID
#
# Takes a variable name as argument. The result is stored in that variable.
generate_id() {
    local -n resvar=$1
    printf -v resvar '%08x%08x%08x%08x' $SRANDOM $SRANDOM $SRANDOM $SRANDOM
}

PG_VERSION=${PG_VERSION:-14}

SPEC_FILE_ORG=/var/db/postgres/specs/spec.json
SPEC_FILE=/tmp/spec.json

echo "Waiting pageserver become ready."
while ! nc -z pageserver 6400; do
     sleep 1;
done
echo "Page server is ready."

echo "Create a tenant and timeline"
generate_id tenant_id
PARAMS=(
     -sb 
     -X POST
     -H "Content-Type: application/json"
     -d "{\"new_tenant_id\": \"${tenant_id}\"}"
     http://pageserver:9898/v1/tenant/
)
result=$(curl "${PARAMS[@]}")
echo $result | jq .

generate_id timeline_id
PARAMS=(
     -sb 
     -X POST
     -H "Content-Type: application/json"
     -d "{\"new_timeline_id\": \"${timeline_id}\", \"pg_version\": ${PG_VERSION}}"
     "http://pageserver:9898/v1/tenant/${tenant_id}/timeline/"
)
result=$(curl "${PARAMS[@]}")
echo $result | jq .

echo "Overwrite tenant id and timeline id in spec file"
sed "s/TENANT_ID/${tenant_id}/" ${SPEC_FILE_ORG} > ${SPEC_FILE}
sed -i "s/TIMELINE_ID/${timeline_id}/" ${SPEC_FILE}

cat ${SPEC_FILE}

echo "Start compute node"
/usr/local/bin/compute_ctl --pgdata /var/db/postgres/compute \
     -C "postgresql://cloud_admin@localhost:55433/postgres"  \
     -b /usr/local/bin/postgres                              \
     -S ${SPEC_FILE}
```

To summarize what is happening here is that we do the following: 

1) We create a tenant, user, company, customer this is a unique database. 
2) We create a timeline for the user. 
3) We update the spec file that will be sent to the compute node in order for it to start up. 

These are steps that normally would not be part fo the compute node but an orchestering layer however since the docker-compose files describe the test step this kind of make sens but is not what we want do do. 

Since I have worked on this project from time to time, i dont know exactly how I changed the k8s resouce but of course leave them in the repo so you can check the diffs if you like to. 

Next we need a k8s cluster, we will use minikbue and thus: 


```bash
minkube start
```

Something that was really tricky to figure out and that took me a long time to understand is how the communication actually works between the service and why it did not work out for me initially and it seems like the pageserver tried to reach the safekeeper over localhost(0.0.0.0). From the neon [Neon architecture blog post](https://neon.tech/docs/introduction/architecture-overview) the architecture is described as: 

![neon architecture](static/img/neon_architecture.avif)

however one pod that is used in the docker-compose files are not in the diagram, the [storage-broker](https://github.com/neondatabase/neon/blob/3710c32aaed4d699451c850fcf7a0dc21520539e/docker-compose/docker-compose.yml#L149). The storage broker turns out to play an important role. From the [docs](https://github.com/neondatabase/neon/blob/release-4179/docs/storage_broker.md) we can understand that the storage broker helps the safekeepers and pageservers learn which nodes also hold their timelines, and timeline statuses there. However the information is based upon the `--listen-pg` and `--listen-http` however these are assumed to be localhost in order to handled this the `--advertise-pg` allows for adding the information what the address should be when we use a service like k8s to run it.

Check out that nothing is running(at least not something you don't like running) using `kubetl get pods --all-namespaces`. The nest step is to deploy the resouces: 

```
kubectl apply -f 
```


After that we need to start do set up. Part of the docker-compose set up is also to create the bucket which we will use to backup our data. This step was done automatic [here](https://github.com/neondatabase/neon/blob/3710c32aaed4d699451c850fcf7a0dc21520539e/docker-compose/docker-compose.yml#L27) we will instead do this through the UI. To do this we need to port-forward and login to minio(user: minio, password: password is the default in the setup). 

```bash
kubectl port-forward svc/minio 9001:9001
```

Then you can just jump to local host and login and create the bucket `minio`. Next step is to create the `tenant` and the `tenantid` to do this we need to comunnicate to the `pageserver` which we will do port-forwarding again: 

```bash
kubectl port-forward svc/pageserver  9898:9898
```

I will create the tenant:  and timeline: but replace with what ever you like, however REMEMBER to update the spec.json that you will later use in the `compute node`. 

tenant: 
```bash
curl -v -sb -X POST -H "Content-Type: application/json" -d '{"new_tenant_id": "de200bd42b49cc1814412c7e592dd6e9"}' http://localhost:9898/v1/tenant/
```

timeline_id:
```bash
 curl  -X POST -H "Content-Type: application/json" -d '{"new_timeline_id": "de200bd42b49cc1814412c7e592dd6e7"}' http://localhost:9898/v1/tenant/de200bd42b49cc1814412c7e592dd6e9/timeline/
```

We are now ready to start the compute node, I will build the `spec.json` file into the container but that is completely up to you how you like to do it! 

```
docker build -t compute -f Dockerfile .
```

If you do this in any otherway remember to updat the k8s resource with the correct container. To allow minikube to find the container is use: 

```bash
eval $(minikube docker-env)
```

The way i set it up now is so the compute container will busy wait in order to not die. 

```bash
kubectl exec -it $YOUR_COMPUTE_NODE -- /bin/bash
```

this step can of course be replaced with just having the correct cmd/args but since I had to hack around to get it to work this was the easiest for me. 


```bash
compute_ctl --pgdata /var/db/postgres/compute \
     -C "postgresql://cloud_admin@localhost:55433/postgres"  \
     -b /usr/local/bin/postgres                              \
     -S spec.json
```

Now you should be ready to connect to your postgres instance, to do it port-forward the compute node

```bash
kubectl port-forward compute-5dc56c7fd9-7cs94  55433:55433
```
and then use psql: 

```bash 
psql -p55432 -h 127.0.0.1 -U cloud_admin postgres
```

Also listening to some of the neon talks it is not clear or not if the compute nodes are running in k8s or as VM:s. 

Happy coding!!!!