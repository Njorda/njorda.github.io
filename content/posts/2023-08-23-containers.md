---
layout: post
title: "What is a container"
subtitle: "Lets build it"
date: 2023-08-23
author: "Niklas Hansson"
URL: "/2023/08/23/what-is-a-container"
---


Today we will deep dive in to containers a subject that is constantly brought up in different posts all over the internet. However the goal with this post is to implement it from scratch in order to learn what it actually is, which is not a light weight VM. 

Containers are build on three core concepts: 
- Namespaces - What can you see. A process in one namespace can not peak in to another namespace
- Cgroups - What can you use, limit on CPU, memory. 
- Chroot - Where our home is and what we can see in the filesystem. 


In order to implement this we will have to know one or two things about system calls. So what is a system call? A system call is simply a call to the operating system. We will need to do this since the things mentioned above is handled by the operating system. Also important to understand is that in unix everything is a file, we will not deep dive in to this to much but feel free to dig down the [rabbit hole](https://en.wikipedia.org/wiki/Everything_is_a_file#:~:text=Everything%20is%20a%20file%20is,through%20the%20filesystem%20name%20space.)

This blog post is heavily based upon the talk: [Container from scratch]()

We will do this in go but the programming language will not matter to much. In order to demo the different concept we will go step by step. A good first step would be to `echo`
something, lets aim for "Hello container". 

```go

```


