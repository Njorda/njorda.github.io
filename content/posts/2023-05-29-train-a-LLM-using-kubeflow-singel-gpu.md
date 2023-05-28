---
layout: post
title: "Training a LLM using kubeflow and a singel gpu"
subtitle: ""
date: 2023-05-29
author: "Johan Hansson"
URL: "2023/05/29-train-a-LLM-using-kubeflow-singel-gpu"
image: "/img/background_2022_07_17.png"
---


This blog post will go over how to train a LLM using kubeflow and one singel GPU, insperationa AND a lot of code has is from this huggingface [blog post](https://huggingface.co/blog/4bit-transformers-bitsandbytes).


The pipeline has two components, the first one ingests data the second starts the training. We will go over them first and then the actual pipeline. The code runs in a dev container and can befound here. To be able to run it you need to have a gcp project and intialized to it. With the comande below since the dockercontainer takes in your default credentials. 

```bash
gcloud default auth
```

You also need to make sure you have a Quota on gpu and more specifycly NVIDIA_TESLA_T4 in the region europe-west1. If not you need to request in from the IAM porta on gcp and quotas, this normaly takes two works days. And in my case I hade to pre pay 10 euros to google. 




## Ingest component 


The ingests components consists of one file ingest_data.py

```python 
from kfp import compiler
from kfp.dsl import component, OutputPath, Output, Dataset, Input, Model, Metrics


def inngest_data(df_path:Output[Dataset],model_id:str,dataset:str='Abirate/english_quotes'): 
    import torch
    import logging
    from datasets import load_dataset
    from transformers import AutoTokenizer
    logging.warning(f'model name is {model_id}')
    tokenizer = AutoTokenizer.from_pretrained(model_id)
    data = load_dataset(dataset)
    data = data.map(lambda samples: tokenizer(samples["quote"]), batched=True)
    data.save_to_disk(df_path.path)

```



## Train component 

```python 

from kfp.dsl import component, OutputPath, Output, Dataset, Input, Model, Metrics

def train_model(df_path:Input[Dataset],model_id:str): 
    import torch
    import sys
    import transformers
    import logging
    from transformers import AutoTokenizer, AutoModelForCausalLM, BitsAndBytesConfig
    from peft import prepare_model_for_kbit_training,LoraConfig, get_peft_model
    from datasets import load_from_disk


    data = load_from_disk(df_path.path)

    logging.warning(f'model id is {model_id}')
    bnb_config = BitsAndBytesConfig(
        load_in_4bit=True,
        bnb_4bit_use_double_quant=True,
        bnb_4bit_quant_type="nf4",
        bnb_4bit_compute_dtype=torch.bfloat16
    )

    tokenizer = AutoTokenizer.from_pretrained(model_id)
    model = AutoModelForCausalLM.from_pretrained(model_id, quantization_config=bnb_config, device_map={"":0})


    model.gradient_checkpointing_enable()
    model = prepare_model_for_kbit_training(model)
    logging.warning('done with all loading')

    def print_trainable_parameters(model):
        """
        Prints the number of trainable parameters in the model.
        """
        trainable_params = 0
        all_param = 0
        for _, param in model.named_parameters():
            all_param += param.numel()
            if param.requires_grad:
                trainable_params += param.numel()
        print(
            f"trainable params: {trainable_params} || all params: {all_param} || trainable%: {100 * trainable_params / all_param}"
        )
    
    config = LoraConfig(
        r=8, 
        lora_alpha=32, 
        target_modules=["query_key_value"], 
        lora_dropout=0.05, 
        bias="none", 
        task_type="CAUSAL_LM"
    )

    model = get_peft_model(model, config)
    print_trainable_parameters(model)



    # needed for gpt-neo-x tokenizer
    tokenizer.pad_token = tokenizer.eos_token
    trainer = transformers.Trainer(
        model=model,
        train_dataset=data["train"],
        args=transformers.TrainingArguments(
            per_device_train_batch_size=1,
            gradient_accumulation_steps=4,
            warmup_steps=2,
            max_steps=10,
            learning_rate=2e-4,
            fp16=True,
            logging_steps=1,
            output_dir="outputs",
            optim="paged_adamw_8bit"
        ),
        data_collator=transformers.DataCollatorForLanguageModeling(tokenizer, mlm=False),
    )
    logging.warning('starting to train')

    model.config.use_cache = False  # silence the warnings. Please re-enable for inference!
    trainer.train()
    logging.warning('----------------done-----------------------')

```


## The pipeline


```python 
ingest_data_component = component(inngest_data,
                                  packages_to_install=["bitsandbytes","peft==0.3.0","transformers==4.29.2", "accelerate==0.19.0", "datasets==2.12.0"],
                                  base_image='europe-docker.pkg.dev/vertex-ai/training/pytorch-gpu.1-13.py310:latest'
                                  ) # https://cloud.google.com/vertex-ai/docs/training/pre-built-containers

train_model_component = component(train_model,
                                  packages_to_install=["bitsandbytes","git+https://github.com/huggingface/peft.git","transformers==4.29.2", "accelerate==0.19.0", "datasets==2.12.0"],
                                  base_image='europe-docker.pkg.dev/vertex-ai/training/pytorch-gpu.1-13.py310:latest'
                                  ) # https://cloud.google.com/vertex-ai/docs/training/pre-built-containers
                                  
```


```python 

# Define the pipeline using the Kubeflow Pipelines SDK
@dsl.pipeline(
    name="test-train",
)
def add_pipeline():
    # Instantiate the ingest_data_component and store its output
    model_id = "EleutherAI/gpt-neox-20b"
    ingest_data = ingest_data_component(model_id=model_id)

    # remberer that you manually have to update yout qouta so you can use gpu's https://stackoverflow.com/questions/53415180/gcp-error-quota-gpus-all-regions-exceeded-limit-0-0-globally
    train_model_component(df_path=ingest_data.outputs['df_path'],model_id=model_id).set_cpu_limit('4').set_memory_limit('60G').add_node_selector_constraint('NVIDIA_TESLA_T4').set_gpu_limit('1') 

    #https://cloud.google.com/vertex-ai/docs/pipelines/machine-types
    #https://cloud.google.com/vertex-ai/pricing#europe


# Compile the pipeline to generate a JSON file for execution
compiler.Compiler().compile(pipeline_func=add_pipeline, package_path="local_run.yaml")

```



```python 
bucket,gcp_project,gcp_service_account

job = aip.PipelineJob(
    #job_id='test' # TODO se in the future
    display_name="First kubeflow pipeline",
    template_path="local_run.yaml",
    pipeline_root=bucket,
    location="europe-west1",
    project=gcp_project,
)

job.submit(service_account=gcp_service_account)
```