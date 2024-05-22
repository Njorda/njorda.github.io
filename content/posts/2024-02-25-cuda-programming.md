---
layout: post
title: "Cuda programming using python and ChatGPT"
subtitle: "Massivly distributed systems"
date: 2024-02-16
author: "Niklas Hansson"
URL: "/2024/02/16/cuda_programming"
---

This blog post aims to be the first in a series about [CUDA](https://developer.nvidia.com/cuda-toolkit) programming. This will be based upon the awesome material from [CUDA MODE](https://github.com/cuda-mode). Today we will focuse on the material in the first two lectures and start to build some simples stuff in CUDA. 



# Background

Cuda is an API from Nvidia for programming GPU:s. GPU:s are pwoerful hardware for highly parallel processes. If the process is not highly parallels GPU:s is probably not the right hardware for you. However a lot of exciting problems are highly parallel. As an example we will today focus on vector search and try to do it as fast as possible. Something that is good to know is that GPU:s are possible to process a lot more data per second than it can load per second. 

# Implement a CUDA kernal

Lets write some code and try it out. I will write most of the code running from a jupyter notebook, but copy it in here to small python snippets. 


# How does GPU:s work

GPU:s are different then CPU:s in many ways but the core difference is that GPU:s are designed for doing a lot of simple computations utalising a lot of cores while CPU:s have a lot fewer cores doing. Further CPU:s are designed for sequential operations while GPU:s are optimized for parallel executions.The difference between GPU:s and CPU:s can be visualised as in figure 1 from [here](https://docs.nvidia.com/cuda/cuda-c-programming-guide/index.html#the-benefits-of-using-gpus) where a lot more transistors are devoted to compute per area compare to caching and control flow. 

IMAGE

Applications have a mix of sequential and parallel parts and thus a mix is needed for most processes. It is important to be aware of your problem and how well it suits for parallel work [Amdahls law](https://en.wikipedia.org/wiki/Amdahl%27s_law) describes this well and can be summarize as, see figure 2 for examples: 

IMAGE

| "the overall performance improvement gained by optimizing a single part of a system is limited by the fraction of time that the improved part is actually used"


# CUDA

| "a general purpose parallel computing platform and programming model"

CUDA was release in 2006 by NVIDIA and allows to write code for NVIDIA GPU:s using C++. We will cut some cornerns in the description of the programming model of CUDA so when you start to dive deep in to CUDA there might be some slight changes to what we learn here today. 

The fundamental building blocks of how to write your CUDA code is built on the concepts of:
- Stream multiprocessors known as SM
- Blocks, each block will run on one SM and be using a nbr of threads(usually a lot, don't be shy and make sure to use a lot of threads). We can use A LOT of blocks. 
- Threads(we wil use a lot of them). We usually want to stay with 256 or something similar in terms of threads per block. 

Using CUDA your code is written around the concept of blocks and threads executing kernels. Kernels are small functions that can be operated in parallel(the whole point) with no guarantee on order and no return. Instead kernels should only mutate memory passed in to the kernel, thus a input tensor and an output tensor. Compare to most programming(depending upon what you are used to) the kernels are written with the hardware in mind where blocks and threads very present in the kernels. Something to bear in mind for later is that blocks have the possibility to share memory that is faster then the memory between blocks which can be used to make you GPU go even more brrrrr(read faster). 


IMAGE. 


Blocks are organized in to one, two or three dimensional grids of threads. As mentioned before the thread blocks must be able to run independently in order to allow for any number of cores. Allowing the program to code to scale with the hardware, up or down. 


All the idx do is to calculate which index in the input and output we should do the operations on. However it is done in a way that is not very intuitive if we take a first look, however we will dig down and hopefully we will get through it so it make sense. 


Max 1024 threads per block