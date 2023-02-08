---
layout:     post 
title:      "Shinylive app with hugo"
subtitle:   "hugo dynamic content"
date:       2023-02-06
author:     "Niklas Hansson"
URL: "/2023/02/06/hugo-with-python-rshiny/"
iframe: "https://nikenano.github.io/shinylive/"
---
# Shinylive app with hugo

This blog post is based upon [RamiKrispin/shinylive](https://github.com/RamiKrispin/shinylive) where we will take a look in to using Shiny with python and leveraging WebAssembly to let it run in the browser with out a backend. This allows for interactive static webpages. 

<iframe
    src="https://nikenano.github.io/shinylive/"
    style="height:800px;width:100%;"
></iframe>

In order to add the a shiny app it needs to be deployed, in this case that is handled through github pages and lives within [https://github.com/NikeNano/shinylive](https://github.com/NikeNano/shinylive). The second step is to add the iframe:

```html
<iframe 
    src="https://nikenano.github.io/shinylive/"
    style="height:800px;width:100%;"
></iframe>
```

How to integrat duckdb

<iframe
    src="https://shell.duckdb.org/"
    style="height:800px;width:100%;"
></iframe>
