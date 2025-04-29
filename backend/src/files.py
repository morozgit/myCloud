import os
from pathlib import Path
import zipfile

from fastapi import APIRouter, File, Form, UploadFile
from fastapi.responses import FileResponse, PlainTextResponse
import logging
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

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


@files_router.post("/download")
async def download_file(file: FileStruct):
    send_for_download_RabbitMQ(str(file.path), str(file.name))
    logger.info("Получен payload:", file)
    url = get_for_download_RabbitMQ()
    return {"download_url": url}


@files_router.post("/upload")
async def upload_file(file: UploadFile = File(...), path: str = Form(...)):
    file_data = await file.read()

    part_size = 8 * 1024 * 1024

    total_parts = (len(file_data) + part_size - 1) // part_size
    part_num = 1
    while file_data:
        part_data = file_data[:part_size]
        file_data = file_data[part_size:]
        send_for_upload_RabbitMQ(part_data, path, file.filename, part_num, total_parts)
        part_num += 1

    return {"name": file.filename, "message": f"Файл отправлен в очередь, путь: {path}"}


@files_router.get("/{file_path:path}")
async def open_file(file_path: str):
    try:
        target_path = (BASE_DIR / file_path).resolve()

        if not str(target_path).startswith(str(BASE_DIR)):
            return PlainTextResponse("Доступ запрещён", status_code=403)

        if not target_path.exists() or not target_path.is_file():
            return PlainTextResponse("Файл не найден", status_code=404)

        if target_path.suffix == ".zip":
            extract_dir = target_path.with_suffix("")
            extract_dir.mkdir(parents=True, exist_ok=True)

            with zipfile.ZipFile(target_path, "r") as archive:
                archive.extractall(path=extract_dir)

            return PlainTextResponse(f"Архив успешно распакован в: {extract_dir}")

        return FileResponse(target_path)

    except Exception as e:
        return PlainTextResponse(f"Ошибка открытия файла: {e}", status_code=500)
