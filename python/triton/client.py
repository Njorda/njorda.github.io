import sys
import numpy as np
import tritongrpcclient



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

if __name__ == "__main__":
    main()
