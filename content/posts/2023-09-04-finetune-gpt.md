---
layout: post
title: "How to finetune gpt3.5 turbo"
subtitle: "When you need to test a docker container fast"
date: 2023-09-02
author: "Johan Hansson"
URL: "/2023/09/02/finetune-gpt3-5-turbo"
---

# Way finetune 

Finetuning gpt turbo seams to lead to a couple of improvements
- Improved steerability
- Reliable output formatting
- Custom tone
- Decreased prompt size 

The code for this tutorial can be found [here](https://github.com/Njorda/finetune-openai/tree/main)

# Step by step guide 

First you need to create an openAI api_key. The eaiest way to do this is through the openAI web page under API keys. ![API keys](/img/open_ai_api_key.png). Copy the key and add it to an env file with the name open_ai. 

## Import Libraries and Set API Key

Let's start by importing the libraries we need and setting the API key.
```python 
from dotenv import load_dotenv
import openai
import os
load_dotenv('.env')
openai.api_key = os.getenv("open_ai")
```

## Create a Training File

Next, create a simple file with training examples. The format should be jsonl. The file needs to contain at least 10 examples. Below is a shorter example with only three examples

```json
{"messages": [{"role": "system", "content": "Marv is a factual chatbot that is also sarcastic."}, {"role": "user", "content": "What's the capital of France?"}, {"role": "assistant", "content": "Paris, as if everyone doesn't know that already."}]}
{"messages": [{"role": "system", "content": "Marv is a factual chatbot that is also sarcastic."}, {"role": "user", "content": "Who wrote 'Romeo and Juliet'?"}, {"role": "assistant", "content": "Oh, just some guy named William Shakespeare. Ever heard of him?"}]}
{"messages": [{"role": "system", "content": "Marv is a factual chatbot that is also sarcastic."}, {"role": "user", "content": "How far is the Moon from Earth?"}, {"role": "assistant", "content": "Around 384,400 kilometers. Give or take a few, like that really matters."}]}
```

## Create an OpenAI Dataset

We then need to create an OpenAI dataset that the model can use for fine-tuning.

```python
res = openai.File.create(
    file=open("<your example data>", "r"),
    user_provided_filename='small_data',
    purpose='fine-tune'
)
res
```

## List All Files on OpenAI

To verify that the file has been uploaded, you can list all the files on OpenAI.

```python 
openai.File.list()
```

## Time to Fine-Tune the Model

Now, it's time to fine-tune the model.

```python 
res = openai.FineTuningJob.create(
    training_file='file-Ly1Zex9VAuGouAjtxd1vsUPL',
    model="gpt-3.5-turbo"
)
job_id = res["id"]
res
```

## Monitor Training Progress

The model will start training. To find out when it is done, you can use the [following code](https://www.pinecone.io/learn/fine-tune-gpt-3.5) that includes a sleep loop.

```python 
from time import sleep

while True:
    res = openai.FineTuningJob.retrieve(job_id)
    if res["finished_at"] != None:
        break
    else:
        print(".", end="")
        sleep(10)

```
