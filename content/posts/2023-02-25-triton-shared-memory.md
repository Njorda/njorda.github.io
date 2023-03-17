---
layout:     post 
title:      "Triton shared memory and pinned memory"
subtitle:   "Optimize triton serving"
date:       2023-02-25
author:     "Niklas Hansson"
URL: "/2023/02/25/triton_shared_memory/"
---


This blog post will go in to depth how to use [shared memory](https://github.com/triton-inference-server/server/blob/main/docs/protocol/extension_shared_memory.md) together with [nvidia triton](https://developer.nvidia.com/nvidia-triton-inference-server) and [pinned memory](https://github.com/triton-inference-server/server/blob/dccc3b43df120f0339ad9c8d262ddddbce7b0ba3/src/main.cc#L593) for model serving. This will continue to build further on the other blog posts related to triton. First we will focuse on shared memory and then move over to also look in to pinned memory and why it matters. 


# Shared memory 

In the triton [examples(python)](https://github.com/triton-inference-server/client/tree/main/src/python/examples) shared memory is often abbreviated as shm. But what is shared memory and why does it matter? The documentation describes the benefits simply as: 

>
> The shared-memory extensions allow a client to communicate input and output tensors by system or CUDA shared memory. Using shared memory instead of sending the tensor data over the GRPC or REST interface can provide significant performance improvement for some use cases.[link](https://github.com/triton-inference-server/server/blob/main/docs/protocol/extension_shared_memory.md)

Thus can be summarize as it allows us to send a reference to memory instead of the data around. A more in depth blog post for shared memory for cuda can be found [here](https://developer.nvidia.com/blog/using-shared-memory-cuda-cc/). 

In this case we will us a single example that can be run locally and we will focus on system shared memory and not GPU shared memory but the logic will be the same. 

For this example we will build upon [one of the example models](https://github.com/triton-inference-server/server/tree/main/docs/examples/model_repository/simple) available from triton in order to focus on the relevant parts and deep dive in to the client. It is important to be aware of that no changes are required on model side or the configurations in order to user shared memory. It is fully handled on the client side. 

```config.pbtxt
name: "shared_memory"
platform: "tensorflow_graphdef"
max_batch_size: 8
input [
  {
    name: "INPUT0"
    data_type: TYPE_INT32
    dims: [ 16 ]
  },
  {
    name: "INPUT1"
    data_type: TYPE_INT32
    dims: [ 16 ]
  }
]
output [
  {
    name: "OUTPUT0"
    data_type: TYPE_INT32
    dims: [ 16 ]
  },
  {
    name: "OUTPUT1"
    data_type: TYPE_INT32
    dims: [ 16 ]
  }
]
```

The model can be found [here](https://github.com/triton-inference-server/server/blob/main/docs/examples/model_repository/simple/1/model.graphdef)


In this example we will run triton locally and the only requirement will be to have docker installed, we will not configure the model to use a GPU but this could easily be changed. In order to run the code the following folder structure is required: 

```
.
└── shared_memory
    ├── 1
    │   └── model.graphdef
    └── config.pbtxt
```


To start the model server you can run: 

```bash
docker run  -it --shm-size=3gb --rm -p8000:8000 -p8001:8001 -p8002:8002 -v$(pwd):/workspace/ -v/$(pwd):/models nvcr.io/nvidia/tritonserver:23.01-py3 bash
```

followed by: 

```bash
tritonserver --model-repository=/models
```

The next step is to write the code for the client. We will start of with a client implementation that don't use any shared memory and then update it. Also we will focuse on [grpc](https://grpc.io/). 


The client we start out looks as follows: 

```python
def main():
    try:
        triton_client = tritongrpcclient.InferenceServerClient(
            url="0.0.0.0:8001", verbose=0)
    except Exception as e:
        print("channel creation failed: " + str(e))
        sys.exit(1) 

    inputs = []
    outputs = []

    inputs.append(
        tritongrpcclient.InferInput("INPUT0", np.asarray([1, 16], dtype=np.int64), "INT32"))
    inputs.append(
        tritongrpcclient.InferInput("INPUT1", np.asarray([1, 16], dtype=np.int64), "INT32"))
    outputs.append(tritongrpcclient.InferRequestedOutput("OUTPUT0"))
    outputs.append(tritongrpcclient.InferRequestedOutput("OUTPUT1"))


    inputs[0].set_data_from_numpy(np.expand_dims(np.asarray([i for i in range(0,16,1)], dtype=np.int32), axis=0))
    inputs[1].set_data_from_numpy(np.expand_dims(np.asarray([i for i in range(0,16,1)], dtype=np.int32), axis=0))
    results = triton_client.infer(model_name="shared_memory",
                                  inputs=inputs,
                                  outputs=outputs)

    output_0_data = results.as_numpy("OUTPUT0")
    output_1_data = results.as_numpy("OUTPUT1")
    print(output_0_data, output_1_data)
```

For an official example you can check [here](https://github.com/triton-inference-server/client/blob/main/src/python/examples/simple_grpc_shm_client.py) but we will try to go more in depth what is happening on the way. 

The first thing we will have to add to our client implementation is:

```python
    triton_client.unregister_system_shared_memory()
    triton_client.unregister_cuda_shared_memory()
```

this is only needed in order to make sure we have no memory regions are register with the server. I have not found if this is actually for the specific client or not but hope so, will do further test on this and update later.  The next step is to register the memory for the inputs and outputs that will be shared between triton and the client. This is slightly different between GPU and CPU, in this case we will focuse on the CPU, system memory. First step is to create a memory region: 

```python 
shm_ip_handle = shm.create_shared_memory_region("input_data",
                                                "/input_simple",
                                                input_byte_size * 2)
```

The first input to the function represent the name, the second the key to the underlying memory region and the third is the memory byte size. The byte size can be calculated based upon the input and the input data type size(int64 vs int32 and so on), in this case it is equal to: `16 * 4 *2` (16 variables of type 32 int, 4 bytes and we have 2 inputs).The name `shm_ip_handle` comes from `share memory input handler` and we will have to make one for the input and output separately. When we have the region we will set the data in the regions as follows

```python
# Put input data values into shared memory
input0_data = np.expand_dims(np.asarray([i for i in range(0,16,1)]
input1_data = np.expand_dims(np.asarray([i for i in range(0,16,1)]
shm.set_shared_memory_region(shm_ip_handle, [input0_data])
shm.set_shared_memory_region(shm_ip_handle, [input1_data],offset=16 * 4)
```

the offset for the input1_data is in order to add the bytes sequentially after the first input. The last step in order to have the shared memory set up is: 

```python
triton_client.register_system_shared_memory("input_data", "/input_simple",
                                            16* 4 * 2)
```

which will register the memory with the triton server. The same has to be done with the output data and then in the end to get the output we do: 

```python 
output0_data = shm.get_contents_as_numpy(
    shm_op_handle, utils.triton_to_np_dtype(output0.datatype),
    output0.shape)
```
where we fetch the data from the the shared memory. Finally in order to clean up:

```python 
triton_client.get_system_shared_memory_status()
triton_client.unregister_system_shared_memory()
shm.destroy_shared_memory_region(shm_ip_handle)
shm.destroy_shared_memory_region(shm_op_handle)
```

Where we fist unregister the memory with the triton server and then destroy it. 

## Analyze performance using [pref_analyzer](https://github.com/triton-inference-server/client/blob/main/src/c++/perf_analyzer/README.md#shared-memory)

In order to see if this is actually worth doing [pref_analyzer](https://github.com/triton-inference-server/client/blob/main/src/c++/perf_analyzer/README.md#shared-memory)(tool to analyze performance for triton) can be used. 


If you have issue setting it up check [this github issue](https://github.com/triton-inference-server/server/issues/4479).



To start the docker container(since we dont like to install things and it is 2023 and docker is making life easier): 

```bash 
docker run -it --network host --shm-size=10gb --ipc host --ulimit memlock=-1 -v $(pwd):/workspace/src nvcr.io/nvidia/tritonserver:23.01-py3-sdk /bin/bash
```

remember that the server also needs to be run with `--ipc host` to work. 

Command to run: 
```bash 
perf_analyzer -m YOU_MODEL_NAME_HERE -u localhost:8001 -i gRPC --input-data YOUR_INPUT_DATA.json input.json
```


# Pinned memory 

[Nvidia blog post](https://developer.nvidia.com/blog/how-optimize-data-transfers-cuda-cc/) about pinned memory and why it matters. 