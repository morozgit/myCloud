from .models.model import FileStruct
from fastapi import APIRouter, UploadFile, File
from pathlib import Path
import os
from fastapi.responses import JSONResponse
from fastapi import Request
from file_job.task import connectRabbitMQ
import shutil
import os

files_router = APIRouter(
    prefix="/api/files",
    tags=["navigation"],
)

UPLOAD_DIR = "/home/user/Downloads/MyCloudFiles"

@files_router.post("/download")
async def download_file(file: FileStruct):
    connectRabbitMQ(str(file.path), str(file.name))
    print("Получен payload:", file)
    return {"message": "Задание на скачивание получено"}


@files_router.post("/upload")
async def upload_file_from_go(file: UploadFile = File(...)):
    # Создаем только имя файла, а не полный путь
    filename = file.filename.split("/")[-1]  # Оставляем только имя файла
    filepath = os.path.join(UPLOAD_DIR, filename)
    print(f"Saving file to: {filepath}")  # Логируем путь, куда сохраняем файл

    # Создаем директорию, если она не существует
    os.makedirs(os.path.dirname(filepath), exist_ok=True)

    # Сохраняем файл в указанную директорию
    with open(filepath, "wb") as buffer:
        shutil.copyfileobj(file.file, buffer)

    return JSONResponse(content={"message": "Файл получен", "filename": filename})