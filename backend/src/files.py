import os
from pathlib import Path

from fastapi import APIRouter, File, Form, UploadFile
from fastapi.responses import FileResponse, PlainTextResponse

from file_job.task import (
    get_for_download_RabbitMQ,
    send_for_download_RabbitMQ,
    send_for_upload_RabbitMQ,
)

from .models.model import FileStruct

files_router = APIRouter(
    prefix="/api/files",
    tags=["navigation"],
)

BASE_DIR = Path(os.getenv("BASE_DIR", "/home"))


@files_router.post("/download")
async def download_file(file: FileStruct):
    send_for_download_RabbitMQ(str(file.path), str(file.name))
    print("Получен payload:", file)
    url = get_for_download_RabbitMQ()
    return {"download_url": url}


@files_router.post("/upload")
async def upload_file(file: UploadFile = File(...), path: str = Form(...)):
    file_data = await file.read()
    send_for_upload_RabbitMQ(file_data, path, file.filename)
    return {"name": file.filename, "message": f"Файл отправлен в очередь, путь: {path}"}


@files_router.get("/{file_path:path}")
async def serve_file(file_path: str):
    try:
        target_path = (BASE_DIR / file_path).resolve()

        if not str(target_path).startswith(str(BASE_DIR)):
            return PlainTextResponse("Доступ запрещён", status_code=403)

        if not target_path.exists() or not target_path.is_file():
            return PlainTextResponse("Файл не найден", status_code=404)

        return FileResponse(target_path)

    except Exception as e:
        return PlainTextResponse(f"Ошибка открытия файла: {e}", status_code=500)
