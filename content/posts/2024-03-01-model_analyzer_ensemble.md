---
layout: post
title: "Cuda programming using python and ChatGPT"
subtitle: "Massivly distributed systems"
date: 2024-03-01
author: "Niklas Hansson"
URL: "/2024/03/01/model_analyzer_ensemble"
---

In this blog post we will deep dive in to the triton [model analyzer](https://github.com/triton-inference-server/model_analyzer/tree/main) and try to optimize a ensemble model, this is a continuation of the previous blog post [Triton shared memory and pinned memory]({{< ref "/posts/2024-03-01-model_analyzer_ensemble.md" >}}).

The model analyzer is a tool from Nvidia for better understanding memory and compute requirements of a Triton server model. Using the model analyzer can also optimizing for through put and latency. Lets get started with deploying our ensemble model, in this case we will run everything on the CPU. 


# Deploy model 

The first step is to deploy a model we can look at. We will use the examples from the [model-analyzer-repo](https://github.com/triton-inference-server/model_analyzer/tree/main/examples/quick-start). We will have the following structure:

```
models
|-- add
    |-- 1
    |   |-- model.py
    |-- config.pbtxt
|-- add_sud
    |-- 1
    |   |-- model.pt
    |-- config.pbtxt
    |-- output0_labels.txt
|-- ensemble_add_sub
    |-- config.pbtxt
...
```

The next step is to run from within, [model_analyzer/tree/main/examples/quick-start](https://github.com/triton-inference-server/model_analyzer/tree/main/examples/quick-start): 

```bash
    docker run -it --shm-size=1gb --rm -p8000:8000 -p8001:8001 -p8002:8002 -v$(pwd):/workspace/ -v/$(pwd):/models nvcr.io/nvidia/tritonserver:24.02-py3 bash
```



This is the link we care about in the end: https://github.com/NVIDIA/TensorRT/blob/c0c633cc629cc0705f0f69359f531a192e524c0f/quickstart/deploy_to_triton/README.md?plain=1#L34



