---
layout:     post 
title:      "Triton shared memory and pinned memory"
subtitle:   "Optimize triton serving"
date:       2023-02-25
author:     "Niklas Hansson"
URL: "/2023/02/25/triton_shared_memory/"
iframe: "https://nikenano.github.io/shinylive/"
---


This blog post will go in to depth how to use [shared memory](https://github.com/triton-inference-server/server/blob/main/docs/protocol/extension_shared_memory.md) together with [nvidia triton](https://developer.nvidia.com/nvidia-triton-inference-server) and [pinned memory](https://github.com/triton-inference-server/server/blob/dccc3b43df120f0339ad9c8d262ddddbce7b0ba3/src/main.cc#L593) for model serving. This will continue to build further on the other blog posts related to triton. First we will focuse on shared memory and then move over to also look in to pinned memory and why it matters. 


# Shared memory 

In the triton [examples(python)](https://github.com/triton-inference-server/client/tree/main/src/python/examples) shared memory is often abbreviated as shm. But what is shared memory and why does it matter? The documentation describes the benefits simply as: 

>
> The shared-memory extensions allow a client to communicate input and output tensors by system or CUDA shared memory. Using shared memory instead of sending the tensor data over the GRPC or REST interface can provide significant performance improvement for some use cases.[link](https://github.com/triton-inference-server/server/blob/main/docs/protocol/extension_shared_memory.md)

Thus can be summarize as it allows us to send a reference to memory instead of the data around. A more in depth blog post for shared memory for cuda can be found [here](https://developer.nvidia.com/blog/using-shared-memory-cuda-cc/). 

In this cas wel will us a single example that can be run locally and we will focus on system shared memory and not GPU shared memory but the logic will be the same. 


## Analyze performance using [pref_analyzer](https://github.com/triton-inference-server/client/blob/main/src/c++/perf_analyzer/README.md#shared-memory)

In order to see if this is actually worth doing [pref_analyzer](https://github.com/triton-inference-server/client/blob/main/src/c++/perf_analyzer/README.md#shared-memory)(tool to analyze performance for triton) can be used. 


If you have issue setting it up check [this github issue](https://github.com/triton-inference-server/server/issues/4479).



To start the docker container(since we dont like to install things and it is 2023 and docker is making life easier): 

```bash 
docker run -it --network host --shm-size=10gb --ipc host --ulimit memlock=-1 -v $(pwd):/workspace/src nvcr.io/nvidia/tritonserver:23.01-py3-sdk /bin/bash
```

Command to run: 
```bash 
perf_analyzer -m YOU_MODEL_NAME_HERE -u localhost:8001 -i gRPC --input-data YOUR_INPUT_DATA.json input.json
```


# Pinned memory 

[Nvidia blog post](https://developer.nvidia.com/blog/how-optimize-data-transfers-cuda-cc/) about pinned memory and why it matters. 