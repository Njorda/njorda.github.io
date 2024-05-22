---
layout: post
title: "Profile nvidia triton model server"
subtitle: "Profile GPU"
date: 2024-03-01
author: "Niklas Hansson"
URL: "/2024/04/07/nsys_profile"
---

In order to profile the triton on the GPU we will use [NVIDIA Nsight Systems](https://developer.nvidia.com/nsight-systems). Installation instructions can be found [here](https://developer.nvidia.com/nsight-systems/get-started). Older version of Nsigh systems can be found [here](https://developer.nvidia.com/gameworksdownload#?dn=nvidia-nsight-graphics--2023-4).The first step is to generate a nsys report. This can be done using the following command(you need to continue reading if you want to get all the traces and build NVIDIA triton image):

```bash
nsys profile --output /MY_OUTPUT_FOLDER/tmp.nsys-rep tritonserver --model-repository
```

In this blog post we will assume that the GPU is connected to a remote machine that you can run docker or kubernetes on. You will NVIDIA Nsight Systems locally and load the saved profiling report from triton. This can be fetched from k8 with something like the following command: 

```bash
kubectl --context YOUR_CONTEXT-n YOUR_NAMESPACE cp TRITON_POD:/MY_OUTPUT_FOLDER/tmp.nsys-rep out.nsys-rep
```

The command will sporadically output an output error but should still work. 

The NVIDIA Triton docker images is built with support for `nsys` cli, however, they do not contain the cuda NVTX markers by default. If you want to use NVTX markers, you have to build Triton with `build.py`, using the “--enable-nvtx” flag. This will provide details around some phases of processing a request, such as queueing, running inference, and handling outputs. Thus we will now work over how to build the triton docker image with NVTX markers added. More info on debugging can be found [here](https://github.com/triton-inference-server/server/blob/main/docs/user_guide/debugging_guide.md). 


Nest step is to check out branch of triton. We will build the `2.41.0` version specifically. 


```bash
git checkout git@github.com:triton-inference-server/server.git
git checkout r24.01
```

We stay on this version since some of the more recent once have had issues with ONNX which we will use for our example. After failing multiple times on the build I decided to jump on to a beefy machine(76 GM RAM and a lot of cores with 400 GB disc) on AWS with Ubuntu to make the build which turned out to be an amazing decision since it removed the errors and made the build fast and successful. 

Time to build the actual image, there are a lot of options and  ended up having to set `--enable-all` in order to get it to work, might be something more streamlined but this is how I did it: 

```bash
sudo ./build.py --enable-all
```

In this blog post we will dive deep in to the memory profiling specifically, [this video](https://www.youtube.com/watch?v=GCkdiHk6fUY) is a really good introduction to memory Analysis with NVIDIA Nsight Compute, which can be installed from [here](https://developer.nvidia.com/tools-overview/nsight-compute/get-started). NVIDIA Nsight Compute vs NVIDIA Nsight systems can be found [here](https://giahuy04.medium.com/introduction-to-nsight-systems-nsight-compute-642ff9578f9f).


The next step is to run triton with nsys: 

```bash
nsys profile --output --force-overwrite=true /MY_OUTPUT_FOLDER/tmp.nsys-rep tritonserver --model-repository
```

In order to get the report you need to shut down the server "ctrl + c" just press once and wait for nsys to collect the traces you should get something like: 

```bash
I0419 19:21:20.191207 159401 grpc_server.cc:2519] Started GRPCInferenceService at 0.0.0.0:8001
I0419 19:21:20.191732 159401 http_server.cc:4623] Started HTTPService at 0.0.0.0:8000
I0419 19:21:20.235119 159401 http_server.cc:315] Started Metrics Service at 0.0.0.0:8002
^CSignal (2) received.
I0419 19:23:17.167877 159401 server.cc:307] Waiting for in-flight requests to complete.
I0419 19:23:17.167903 159401 server.cc:323] Timeout 30: Found 0 model versions that have in-flight inferences
I0419 19:23:17.168044 159401 server.cc:338] All models are stopped, unloading models
I0419 19:23:17.168061 159401 server.cc:345] Timeout 30: Found 1 live models and 0 in-flight non-inference requests
Generating '/tmp/nsys-report-bacf.qdstrm'
I0419 19:23:17.169166 159401 onnxruntime.cc:2838] TRITONBACKEND_ModelInstanceFinalize: delete instance state
I0419 19:23:17.379429 159401 onnxruntime.cc:2762] TRITONBACKEND_ModelFinalize: delete model state
I0419 19:23:17.477649 159401 model_lifecycle.cc:612] successfully unloaded 'model' version 202404091557
[1/1] [=16%                        ] test.nsys-repI0419 19:23:18.168230 159401 server.cc:345] Timeout 29: Found 0 live models and 0 in-flight non-inference requests
[1/1] [========================100%] test.nsys-rep
Generated:
    /test.nsys-rep
```


If you see strange errors that involve: 

```bash
...
magic number mismatch
...
```

- possible there is a msissmatch between the triton container nsys `nsys --version` and NVIDIA Nsight Systems. 
- It could also be as suppressed error and you fail to write `Failed to create '/opt/tritonserver/report1': Permission denied. ` and needs to add a volume to which you have write access. 

Take the file down locally using kubectl cp, example below. DONT BELEVE IN THE ERROR MESSAGE IT IS A SCAM: 


```bash
kubectl --context CONTEXT -n NAMESPACE  cp POD_NAME:/models/out.nsys-rep out.nsys-rep
```

Now it is time to load up the report in Nsight.

Good luck and hope you manage to squeeze out some extra performance from the GPU!