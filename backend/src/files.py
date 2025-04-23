from .models.model import FileStruct
from fastapi import APIRouter, UploadFile, File, Form
from pathlib import Path
import os
from fastapi.responses import JSONResponse
from fastapi import Request
from file_job.task import (
    send_for_download_RabbitMQ,
    get_for_download_RabbitMQ,
    send_for_upload_RabbitMQ
    )
import shutil
import os

files_router = APIRouter(
    prefix="/api/files",
    tags=["navigation"],
)


@files_router.post("/download")
async def download_file(file: FileStruct):
    send_for_download_RabbitMQ(str(file.path), str(file.name))
    print("Получен payload:", file)
    url = get_for_download_RabbitMQ()
    return {"download_url": url}


@files_router.post("/upload")
async def upload_file(
    file: UploadFile = File(...),
    path: str = Form(...)
):
    file_data = await file.read()
    send_for_upload_RabbitMQ(file_data, path, file.filename)
    return {"name": file.filename, "message": f"Файл отправлен в очередь, путь: {path}"}
