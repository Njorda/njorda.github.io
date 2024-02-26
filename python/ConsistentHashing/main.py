from fastapi import FastAPI
from pydantic import BaseModel


# Example here: https://gist.github.com/soekul/a240f9e11d6439bd0237c4ab45dce7a2
# Example here: https://medium.com/@sent0hil/consistent-hashing-a-guide-go-implementation-fe3421ac3e8f
# https://levelup.gitconnected.com/binary-search-trees-in-go-58f9126eb36b
# https://github.com/stathat/consistent/blob/master/consistent.go


app = FastAPI()


nodes = {
    "node1": "8081",
    "node2": "8082",
    "node3": "8083",
}

class Object(BaseModel):
    key: str
    name: str | None = None
    value: float

    model_config = {
        "json_schema_extra": {
            "examples": [
                {
                    "key": "Foo",
                    "name": "A very nice Item",
                    "value": 3.2,
                }
            ]
        }
    }


# Here we add a node
@app.put("/node/{node_id}")
async def add(node_id: int, item: str):
    nodes[node_id] = item
    return nodes

@app.put("/add/{key_id}")
async def add(key_id: int, item: Object):
    results = {"item_id": key_id, "item": item}
    return results

@app.put("/get/{key_id}")
async def get(key_id: int,):
    return {"item_id": key_id, "item": ""}

@app.put("/delete/{key_id}")
async def get(key_id: int,):
    return "delete"