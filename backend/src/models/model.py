from pydantic import BaseModel


class FileStruct(BaseModel):
    path: str
    name: str