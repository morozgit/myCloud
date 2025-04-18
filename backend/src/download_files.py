from .models.model import FileStruct
from fastapi import APIRouter, UploadFile, File
from pathlib import Path
import os
from fastapi.responses import JSONResponse
from fastapi import Request
from file_job.task import sendToRabbitMQ, getFromRabbitMQ
import shutil
import os

files_router = APIRouter(
    prefix="/api/files",
    tags=["navigation"],
)

@files_router.post("/download")
async def download_file(file: FileStruct):
    sendToRabbitMQ(str(file.path), str(file.name))
    print("Получен payload:", file)
    url = getFromRabbitMQ()
    return {"download_url": url}
