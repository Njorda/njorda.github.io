---
layout:     post 
title:      "Trition with post and pre processing"
subtitle:   "Ensemble model"
date:       2023-02-01
author:     "Niklas Hansson"
URL: "/2023/02/01/Trition_with_post_and_pre_processing/"
---


# Prompt engineering

With the rise of GPT-3 and Stable diffusion the concpet of promt engineering has gain more and more traction. According to wikipedia the task can be described as

>
>Prompt engineering typically works by converting one or more tasks to a prompt-based dataset and training a language model with what has been called "prompt-based learning" or just "prompt learning"
> - [Wikipedia](https://en.wikipedia.org/wiki/Prompt_engineering#:~:text=Prompt%20engineering%20is%20a%20concept,of%20it%20being%20implicitly%20given.)


Prompts are inputs for models that expect text as input however the output can very most famously images or text. 

In this tutorial we will play around with prompt engineering for images. 

# Example

In this example we will use stable diffusion using triton. 

- https://github.com/triton-inference-server/server/tree/main/docs/examples/stable_diffusion

Build something cool with this. 