---
layout:     post 
title:      "Postgress extention using ChatGPT3"
subtitle:   "Create feature embeddings inside postgres"
date:       2023-03-11
author:     "Niklas Hansson"
URL: "/2023/03/11/postgress_extention_using_chatgpt3/"
---


# pgx rust Postgres extension

This blog post will dive into building [postgres extensions](https://www.postgresql.org/docs/current/sql-createextension.html) using [pgx](https://www.postgresql.org/docs/current/sql-createextension.html) in rust. In order to do something existing we will ride the hype curve and integrate ChatGPT3 into postgres. 

TLDR: [repo](https://github.com/NikeNano/PostChat)

# Setup 

First step is to set up `pgx`: 

```bash
$ cargo install --locked cargo-pgx
$ cargo pgx init
```

We will then create a new create using: 

```bash
$ cargo pgx new my_extension
$ cd my_extension
```

Which should give you something like: 

```
.
├── Cargo.toml
├── postchat.control
├── sql
└── src
    └── lib.rs
```

to make sure everything is set up correct run: 

```bash
$ cargo pgx run
```

This will give you a `psql` shell, first step is to load your extensions: 

```
$ postchat=# CREATE EXTENSION postchat;
```

if the extension already exsist first drop it `DROP EXTENSION postchat;`

Useful command is also `\df` to show all functions. 

```SQL
SELECT * FROM hello_postchat();
```

## ChatGPT3


Ok, so the dummy test works. Next step is to make a client call to ChatGpt3 using rust, in this case we will use the [openai-api](https://github.com/deontologician/openai-api-rust/) crate. A simple example to invoke the API is: 


```rust 
use openai_api::Client;

#[tokio::main]
async fn main() {

    let api_token = std::env::var("OPENAI_SK").unwrap();
    let client = Client::new(&api_token);
    let prompt = String::from("Once upon a time,");
    println!(
        "{}{}",
        prompt,
        client.complete_prompt(prompt.as_str()).await.unwrap()
    );
}
```

before a token to OPENAI needs to be set:

``` bash 
$ export OPENAI_SK=YOUR_TOKEN
```
to run use: `cargo run` which gives: 

```bash 
cargo run
    Finished dev [unoptimized + debuginfo] target(s) in 0.31s
     Running `target/debug/postchat`
Once upon a time, the belief in chaos and entropy was dominant, a belief credited to the first French
```
how ever the input will differ between runs. Now when we know how to call rust we need to get it inside a pgx rust extension. This primarly builds around calling `async` code from a none `async` function and boils down to: 


```rust 
async fn prompt(input: &str) -> String {
    let api_token = std::env::var("OPENAI_SK").unwrap();
    let client = Client::new(&api_token);
    let prompt = String::from(input);
    let res = client.complete_prompt(prompt.as_str()).await.unwrap();
    return res.to_string()
}

#[pg_extern]
fn chaty(input: &str) -> String {
    // Create the runtime
    let mut rt = Runtime::new().unwrap();
    // Spawn a future onto the runtime
    return rt.block_on(prompt(input)); 
}
```


Lets create a table to try it out and add some data to it: 

```SQL
CREATE TABLE testy(
    id serial PRIMARY KEY,
    chat TEXT
);
```

```SQL
INSERT INTO testy (chat)
VALUES
    ('Hello'),
    ('tell me something cool'),
    ('why is that cool');
```

now we can run some chatgpt from Postgres

```SQL
SELECT chaty(chat) FROM testy;
```

This is a dummy example and will not scale due to multiple reasons out of which one is that that we create a new client for each entry in the table, but is still a fun hack :). 
