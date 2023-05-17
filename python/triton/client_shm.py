import sys
import numpy as np
import tritongrpcclient
import tritonclient.utils.shared_memory as shm


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
    

    input_byte_size = 16 * 4 * 2
    output_byte_size = 16 * 4 * 2
    shm_ip = shm.create_shared_memory_region("input_data",
                                                    "/input_simple",
                                                    input_byte_size)

    input0 = np.expand_dims(np.asarray([i for i in range(0,16,1)], dtype=np.int32))
    input1 = np.expand_dims(np.asarray([i for i in range(0,16,1)], dtype=np.int32))


    shm.set_shared_memory_region(shm_ip, [input0])
    shm.set_shared_memory_region(shm_ip, [input1],
                                 offset=input_byte_size/2)

    triton_client.register_system_shared_memory("input_data", "/input_simple",
                                                input_byte_size)
    inputs[0].set_shared_memory("input_data", input_byte_size)
    inputs[1].set_shared_memory("input_data", input_byte_size)


    shm_op_handle = shm.create_shared_memory_region("output_data",
                                                    "/output_simple",
                                                    output_byte_size)


    outputs.append(tritongrpcclient.InferRequestedOutput("OUTPUT0"))
    outputs.append(tritongrpcclient.InferRequestedOutput("OUTPUT1"))
    outputs[0].set_shared_memory("output_data", output_byte_size/2)
    outputs[1].set_shared_memory("output_data",
                                  output_byte_size/2,
                                  offset=output_byte_size/2)



    results = triton_client.infer(model_name="shared_memory",
                                  inputs=inputs,
                                  outputs=outputs)

    output0_data = shm.get_contents_as_numpy(
                        shm_op_handle, 
                        utils.triton_to_np_dtype(output0.datatype),
                        output0.shape)
    output1_data = shm.get_contents_as_numpy(shm_op_handle,
                        utils.triton_to_np_dtype(
                        output1.datatype),
                        output1.shape,
                        offset=output_byte_size)
    print(output_0_data, output_1_data)

if __name__ == "__main__":
    main()