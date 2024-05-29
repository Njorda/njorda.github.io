---
layout: post
title: "Mixtral of Experts"
subtitle: "Paper review: Mixtral of Experts"
date: 2024-01-21
author: "Niklas Hansson"
URL: "/2024/01/21/mixtral-of-expert-review.md"
---

In this blog post we will review the [Mixtral of Experts](https://arxiv.org/pdf/2401.04088.pdf)paper from Mistral Ai. A short blog post can be found [here](https://mistral.ai/news/mixtral-of-experts/). Lets go. 

A key part of the paper is the following section from the abstract:

| We introduce Mixtral 8x7B, a Sparse Mixture of Experts (SMoE) language
model. 

Naturally one might as what is a Mixture of Experts(MOE) and what is a Sparse MOE? A MOE model is a model build up of a set of models, experts. These expert models are then used for potentially different parts of the input to improve the out put while limiting the compute. Compute is a key component here to consider. And the goal with MOE is to achieve a higher result for a lower compute compare to a dense model of he same size. These models are are also pretrained a lot faster then a dense model. 

The Sparse MOE(SMoE) have two key features in terms of the mixture and sparsity: 

It sounds a lot like drop out to be honest ... 

The Hugging face [article](https://huggingface.co/blog/moe) is better. 