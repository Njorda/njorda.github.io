---
layout: post
title: "DBOS: A Database-Oriented
Operating System"
subtitle: "Will the future OS be built on top of a database?"
date: 2023-04-14
author: "Niklas Hansson"
URL: "/2023/05/19/operating_system_databases"
---

 A group of researches are proposing a radical change of the future operating system. Replacing the fundamental idea from Unix that everything is a file and instead relying on concepts from the database world a operating system that supports large scale distributed applications in the cloud can be built. 

 | Everything is a file

The core principles suggested to achieve this is:
- Store all application in tables in a distributed database
- Store all OS state in tables in a distributed database. 

The suggested architecture(four levels) looks as follows: 
![Arch.](/img/DBOS.png)


Links: 
- [Video](https://www.youtube.com/watch?v=eB4bJqDzsU8)
- [Blog](https://dbos-project.github.io/blog/intro-blog.html)
- [Paper](https://arxiv.org/abs/2007.11112)

Will continue on this topic in a later blog post where we deep dive in to the research. 




