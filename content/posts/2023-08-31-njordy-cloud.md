---
layout: post
title: "Njordy web services"
subtitle: "How to build a cloud provider - high level plan"
date: 2023-08-31
author: "Niklas Hansson"
URL: "/2023/08/31/build-a-cloud"
---

This is the first blog in the series - `How to build your own cloud provider` the aim is to build a bare minimum k8s set up where a user can create the following resources: 

- A VM
- A postgres DB
- A bucket

The goal is to run the different services on to of a k8 cluster. A key concept of a cloud provider is also to make money, in this case the goal is not to make money but we do want to keep track of "users" spend. Thus we will also deploy [kubectl-cost](https://github.com/kubecost/kubectl-cost). THe goal is that this could be run on any k8s cluster, however since we want to replicate how it would look for a cloud provider the goal is to actually have some actual hardware. In this case we will use [Raspberry Pi](https://en.wikipedia.org/wiki/Raspberry_Pi).

# Things to buy


## Hardware

For our cluster we will use 3 Raspberry Pi 3+ since that was what I had lying around and could buy from friends. The set up described is inspired by [this blog post](https://anthonynsimon.com/blog/kubernetes-cluster-raspberry-pi/). 

Hardware: 

- 3 x Raspberry Pi( will use 3+)
- 3 x Raspberry Pi PoE(Power over Ethernet) heads
- 3 x SD cards
- 5 Ethernet cables
- Cloudlet case
- 1x TP-Link 5-Port Gigabit PoE Switch
- 1x TP-Link Nano Router WLAN

## Domain

After googling a lot I settled for [ArcticCompute.com](http://arcticcompute.com/). 


# Setting up the hardware

# Installing Ubuntu on the Raspberry Pi:s

# Setting up K8s

# Next step

Now when we have a complete K8s cluster running the next step is to start to deploy the different cloud compute parts and the webpage and servers keeping it together. 
