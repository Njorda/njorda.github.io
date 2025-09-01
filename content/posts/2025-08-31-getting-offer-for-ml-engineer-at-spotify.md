---
layout: post
title: "How I Prepared for and Got a Spotify ML Engineer Offer"
subtitle: "Sccing the interviews"
date: 2025-08-30
author: "Johan Hansson"
URL: "/2025/08/31/ml_engineer_spotify"
---

This blog post covers how I prepared for my interviews at Spotify and ultimately received an offer for an ML engineer role. For some people, the interview process at FAANG-style companies is as easy as a walk in the park, but for me, it required extensive preparation. The time you'll need is probably highly individual, but I really, really wanted this job, so I went all in! Note that I won't be sharing any specific interview questions.

## The Interview Process

My interview process consisted of:
1. Call with HR
2. A 1-hour Machine Learning Tech Screen
3. Assessment day with four interviews:
   - ML Data Interview 
   - Values Interview
   - ML System Design 
   - ML Depth Interview

I ended up doing an extra ML depth interview, but the normal process typically includes just the ones listed above.

I broke my preparation down into two parts: the tech screen and the assessment day, with different strategies for each. For both parts, I tried to read about others' experiences—not to find exact questions or scripts, but rather to get a feel for the interviews and start thinking about how I'd handle various questions. While there are countless videos and articles available, in hindsight, none of them really matched what I experienced at Spotify.

## Tech Screen Preparation

I knew there would be three parts: a short introduction, technical ML questions, and a coding section. Feeling confident about the first two parts, I focused all my energy on coding preparation.

There are numerous platforms for developing coding skills. I used two:
1. [Leetcode.com](https://leetcode.com/)
2. [AlgoExpert](https://www.algoexpert.io/product)

What mattered most to me was training on similar problems and algorithms in different ways across multiple problems. This helped me learn patterns rather than memorize specific questions. I focused on easy to medium problems—I never attempted hard ones, to be honest.

For LeetCode, I used the freemium version, but I paid for AlgoExpert, which suited me well. I tried to balance deep understanding with volume—solving 1-3 medium problems daily and several easy ones.

However, the most crucial part was asking friends to conduct mock interviews. Solving problems is one thing, but doing it while explaining your thought process to another person is entirely different. I highly recommend this—and ask for feedback! If you don't have programming friends, that's fine. You should be able to explain your solution so anyone can understand it. There was a steep learning curve here, but I felt massive improvement after about five mock interviews.

Additionally, if you're not comfortable with it already, you need to master time and memory complexity. This is a key software engineering skill, and you will get questions about it.

## Assessment Day Preparation

For this part of the process, I focused on three areas: recommendation systems, ML system design, and reading articles (mainly from Spotify's blog).

### Recommendation Systems

I needed to refresh my knowledge of recommendation systems and spent considerable time studying through blogs, YouTube, and books. Some helpful resources:
- https://www.youtube.com/watch?v=9vBRjGgdyTY
- https://www.hopsworks.ai/dictionary/two-tower-embedding-model
- https://www.uber.com/en-TW/blog/innovative-recommendation-applications-using-two-tower-embeddings/
- https://cloud.google.com/blog/products/ai-machine-learning/scaling-deep-retrieval-tensorflow-two-towers-architecture
- https://medium.com/data-science/recommender-systems-a-complete-guide-to-machine-learning-models-96d3f94ea748

I focused on content-based filtering, collaborative filtering approaches, and two-tower networks. Within these areas, I also studied how real-time inputs would affect the system.

### Reading Technical Articles

I approached this in three ways:
1. Read relevant articles from Spotify's engineering blog
2. Researched my interviewers and read their published articles/papers
3. Studied the most defining papers in my specific ML area (most ML roles are specialized—for example, gen-AI, recommendations, etc.)

Some recommended engineering blogs:
- https://huggingface.co/blog
- https://blog.doordash.com/en-us
- https://engineering.atspotify.com/
- https://netflixtechblog.com/

### ML System Design

This was where I spent the most time. I started with [Machine Learning System Design Interview](https://www.amazon.com/Machine-Learning-System-Design-Interview/dp/1736049127) to establish a framework for breaking down ML system design interviews. There are numerous good approaches, and I ultimately used a mix of the book's methodology and the [MLEpath YouTube channel](https://www.youtube.com/watch?v=XN2ymraj27g).

I divided the process into 7 steps. Below are some examples of what you could discuss in the different steps. A key point for me and probably most people is to visualize your design:

**1. Clarifying Requirements**
- Supervised, unsupervised, or reinforcement learning problem?
- Structured or unstructured data?
- Scale requirements
- Latency requirements
- What data do I have access to?
- Batch or real-time processing?

This is the part where you ask a lot of questions to understand the problem. Take your time here—it will really help you. Feel free to take notes. When I felt I had all the information I needed, I summarized it for the interviewer and asked if I had understood everything correctly. I did this not only to avoid going in the wrong direction but also to show the interviewer I had understood the task.

**2. Framing the Problem as an ML Problem**
- What steps will be included in the ML pipeline?
- Example: 100 million songs → 10,000 songs (optimize for latency and recall)
- 10,000 songs → 100 songs (optimize for precision)

Not all problems are ML problems, so to be able to solve them with ML, you need to frame them as ML problems. An example could be transforming "Too many users are canceling their Spotify Premium subscriptions" into "Given a user's behavior over the last 30 days, predict the probability they'll cancel their subscription in the next 7 days."

**3. Data Preparation**
- Normalization
- Batching, even buckets, distribution-based buckets
- Handling outliers
- Creating embeddings
- Splitting datasets
- Handling imbalanced datasets
- Data distributions

No matter how good your model is, it will never be better than your data. In this step, I talked about how I would prep the data, how I would split it, and what considerations I would take.

**4. Model Development**
- Baseline model
- Optimizers
- Loss functions
- ML algorithms
- Model architecture (activation functions, skip connections, normalization, dropout)
- Hyperparameter tuning
- Batch size considerations

This is perhaps the most fun part: what model would be a good fit for the problem, how should we train it, etc. There are tons of different things to talk about, and you won't have time to cover it all.

**5. Evaluation**
- Evaluation metrics (confusion matrix, F1-score, accuracy, MAP, precision, recall, RMSE)
- Model loss during training
- Offline evaluation
- Online evaluation
- Shadow deployment vs. A/B testing

For me, evaluation has always been the important step to get right. Similar to model development, there's a lot to talk about—use your previous experience. I also spent a good amount of time discussing what metrics I believe are a good fit for the problem at hand and how they translate to the actual business problem we have.

**6. Model Monitoring**
- Data drift
- Model drift

If we finally decide to take the model out into the real world, we need to monitor it—distributions can change, etc. Here you can talk a lot about your previous experience and how you would automate things since it's hard to manually monitor a model 24/7.

**7. Deployment and Serving**
- CI/CD 
- Kubernetes


To gain hands-on experience, I conducted mock interviews with friends (don't forget to time them!). For a realistic experience, you need people who work in tech and your target area—practicing with a friend in accounting probably won't help much. When I felt ready, I used [Prepfully](https://prepfully.com/). Interviewing with someone from a top company gave me valuable experience and recreated the real interview feeling. I did two mock interviews through the platform. While not cheap, it was worth it for me. The feedback I received transformed how I approached interviews and really helped.

## Key Takeaways

A crucial piece of feedback was to always explain my thinking, present options, and then explain why I chose a specific approach. This demonstrates deep knowledge and reasoning behind selecting a particular optimizer, model, etc. Try to do this as much as possible.

**Time management is KEY!** You have limited time, and spending too much on early steps leaves insufficient time for topics like model monitoring. Deployment and serving often aren't prioritized, but be prepared to discuss them.

From this and other interview experiences, I've learned that you can sometimes get stuck or head in the wrong direction. Really listen to and understand the interviewer's questions—they often try to help guide you back on track.

Finally, learn the tools you'll use in the interview to avoid wasting time. For me, this meant learning Figma for the ML system design portion. I prepared by learning exactly which shapes to use for different components (databases, APIs, etc.) depending on the problem.

## Final Thoughts

Remember, preparation is highly individual. What worked for me might not work exactly the same for you, but I hope sharing my experience helps you in your journey. Good luck with your interviews!