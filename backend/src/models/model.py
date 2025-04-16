from pydantic import BaseModel


class File(BaseModel):
    path: str
    name: str