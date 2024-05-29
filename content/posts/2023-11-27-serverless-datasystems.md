---
layout: post
title: "Serverless data systems"
subtitle: "Separate compute and storage"
date: 2023-11-27
author: "Niklas Hansson"
URL: "/2023/11/27/serverless-data-systems"
---

This blog post is a summary of these three blog posts: 

- [The Architecture Of Serverless Data Systems](https://jack-vanlightly.com/blog/2023/11/14/the-architecture-of-serverless-data-systems)
- [building-and-operating-a-pretty-big-storage-system](https://www.allthingsdistributed.com/2023/07/building-and-operating-a-pretty-big-storage-system.html) video can be found [here](https://www.youtube.com/watch?v=sc3J4McebHE)

One of the core concepts of most modern database systems is separation of storage and compute and the possibility to scale these independently. Especially in relation to serverless applications. 

CSP - cloud solution provider. 
MT - Multi tenant, multiple users. 


Multi tenancy is about resource sharing - good point

Shared processes - logic means it is handled with code I guess? 
Containers - it is handled through k8s(faregate uses firecracker though .. )
Virtualization - Multiple VM:s per host


Consumption based pricing. 


Consistency requirements is interesting! 

Data splitting based upon hot spotting(it is all about avoiding hotspots)

Kafka is just a storage API

| The declarative nature of SQL is a major strength, but also a common source of operational problems. This is because SQL obscures one of the most important practical questions about running a program: how much work are we asking the computer to do?

This quote is good! 

Tenant isoloation - nosy neighbour problems. Bare metal is a solution though

|  that concurrent tenants that are served from shared hardware resources appear to be served from their own dedicated services.

| “Everything got faster, but the relative ratios also completely flipped.”


| The separation of storage and compute is becoming increasingly a reality, partly because the network no longer presents the bottleneck it used to.

- DDR - ram
- PCIe transfer on motherboard


High latency on cloud storage, hard to use with OLTP. 

| Engineers can choose to include object storage in their low-latency system but counter the latency issues of object storage by placing a durable, fault-tolerant write-cache and predictive read-cache that sits in front of the slower object storage. This durable write-cache is essentially a cluster of servers that implement a replication protocol and write data to block storage. In the background, the cluster uploads data asynchronously to object storage obeying the economic pattern of writing fewer, larger files.


This is exactly what neon does ... 



Heat management - S3 are masters on this ... 



| Large datasets such as big databases or high throughput event streams must be sharded in order to spread the load effectively over the fleet. 

Shards in kinesis!

| High resource utilization through resource pooling is the name of the game, but doing so with solid tenant isolation and predictable performance is the challenge.

This is what the cloud is ... (and more)