---
layout:     post 
title:      "Triton model analyzer"
subtitle:   "Model analyzer for ensemble model"
date:       2023-02-28
author:     "Niklas Hansson"
URL: "/2023/02/25/triton_model_analyzer/"
---


In this blog post we will try out the model analzer from [nvidia triton](https://github.com/triton-inference-server/server) on an [ensemble model](https://github.com/triton-inference-server/python_backend). In order to have an example to optimize we will use this [example](https://github.com/triton-inference-server/python_backend/tree/main/examples/preprocessing) with the modification that we will run it without GPU, in my case I will run this on my laptop but it should preferably be done on a production machine in order to get the best possible optimizations. 

First step is to clone the library: 
```bash 
git clone https://github.com/triton-inference-server/python_backend.git
```

from `python_backend/examples/preprocessing/` run the following(we will go over it quickly but for more details check out the [repo](https://github.com/triton-inference-server/python_backend/tree/main/examples/preprocessing)): 

1. Set up the folder structures
```bash
$ mkdir -p model_repository/ensemble_python_resnet50/1
$ mkdir -p model_repository/preprocess/1
$ mkdir -p model_repository/resnet50_trt/1

# Copy the Python model
$ cp model.py model_repository/preprocess/1
```

2. Converting PyTorch Model to ONNX format:
```bash
$ docker run --gpus=1 -it -v $(pwd):/workspace nvcr.io/nvidia/pytorch:25.01-py3 bash
$ pip install numpy pillow torchvision
$ python onnx_exporter.py --save model.onnx
$ trtexec --onnx=model.onnx --saveEngine=./model_repository/resnet50_trt/1/model.plan --explicitBatch --minShapes=input:1x3x224x224 --optShapes=input:1x3x224x224 --maxShapes=input:256x3x224x224 --fp16
```

3.
```
$ docker run --gpus=all -it --rm -p8000:8000 -p8001:8001 -p8002:8002 -v$(pwd):/workspace/ -v/$(pwd)/model_repository:/models nvcr.io/nvidia/tritonserver:23.01-py3 bash
$ pip install numpy pillow torchvision
$ tritonserver --model-repository=/models
```





# https://github.com/triton-inference-server/model_analyzer/issues/400

