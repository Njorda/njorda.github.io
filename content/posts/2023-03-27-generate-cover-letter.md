---
layout: post
title: "How to faster apply for a new job"
subtitle: "Use Chatgpt3 to generate your cover letter"
date: 2023-03-28
author: "Niklas Hansson"
URL: "/2023/03/23/chatgpt3-cover-letter-generator"
---

Applying for jobs are fun, but it is not fun to write cover letters at least not compare to hacking on some new open source tool or testing out something new. Therefore the goal is to reduce the time to apply for a job thus the goal today is to generate cover letters. For this we will need the following: 

- CV, I will copy paste mine from linkedin(minimal effort and I spend way to much time on linkedin so it is up to date)
- Job add(will take one from linkedin that sounds interesting)


# Get the CV. 

Go to your profile, select more(mine is in swedish but you get it) and download as pdf, find [a guide here](https://www.linkedin.com/help/linkedin/answer/a541960/save-a-profile-as-a-pdf?lang=en). 

# Find a job

I found the following job which actually sound pretty interesting

`Data Backend Engineer (Tools and Simulations) - Games`

> As a Data Backend Engineer in the data team at Embark you will be a gamemaker focused on building the tools and simulation models needed to support the creation of the most fun, balanced, and engaging games possible.

# Prompt time

In order to give Chatgp3 some instructions I structured my chat as follows: 


First instruction: 

```bash
Based upon this profile: 

-------------------------------------------------------------------------------------------------
```

I used `-` to separate the different sections. The next step where to copy paste every thing from the pdf of my profile, no time to clean up just take it all. I finished it up with some `-` similar to above to make it a bit more structured. After that I added: 


```bash
Generate a cover letter for this job: 
```

followed by the job add. The following where the cover letter generated: 


```text
Dear Niklas Hansson,

I am excited to apply for the position of Data Backend Engineer at Embark. With my extensive experience in machine learning and data engineering, I am confident that I can make significant contributions to your team.

At Northvolt, I have been working on applications for electrode manufacturing where I wrote code in Python, Go, and SQL to aggregate data. As a Data Scientist at Sandvik, I worked on predictive models using deep learning techniques. Moreover, I have contributed to Jupyter Hub, Apache Beam, Kubeflow, Argo, and Tekton, and I am an active open-source enthusiast.

I am skilled in distributed systems, cloud infrastructure, and building data-driven applications, which are some of the core requirements for the Data Backend Engineer role at Embark. I have also worked on database optimization and written ad-hoc SQL queries for data analysis, which will come in handy while collaborating with other teams to build solutions.

I am passionate about interactive experiences, games, and new technologies, and I would love to apply my skills to the gaming industry. I have experience building services that are easy to deploy, scale, and maintain, which aligns with Embark's engineering practices.

Thank you for considering my application. I am excited about the opportunity to join your team and contribute to the development of fun, engaging games.

Sincerely,
[Your Name]
```

# Conclusion

The cover letter is pretty good and I would probably not make it better myself after all I am a engineer not a write. It managed to screw up my name vs the recruiting manager but that it. Will be interesting to see how this affects the importance of cover letters in the future. Also what is really cool is that I can now chat about my cover letter to get it improved how ever I fancy. I think this will be a great tool to work smarter and reduce the time spend on boring things, however it will also flood the internet with generated content and what is original and real will be harder and harder to distinguish. 