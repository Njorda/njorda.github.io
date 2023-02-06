---
layout:     post 
title:      "Shinylive app with hugo"
subtitle:   "hugo dynamic content"
date:       2023-02-06
author:     "Niklas Hansson"
URL: "/2023/02/06/hugo-with-python-rshiny/"
---
# Shinylive app with hugo

This blog post is based upon [RamiKrispin
/
shinylive](https://github.com/RamiKrispin/shinylive) where we will take a look in to using Shiny with python and leveraging WebAssembly to let it run in the browser with out a backend. This allows for interactive static webpages. 


<iframe src="https://nikenano.github.io/shinylive/"></iframe>


<!--shiny.html-->
<iframe src="https://nikenano.github.io/shinylive/{{ index .Params 0 }}/"
        style="height:{{ index .Params 1 }}px;width:100%;border:none;overflow:hidden;" scrolling="no"></iframe>