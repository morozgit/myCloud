from .models.model import File
from fastapi import APIRouter, HTTPException, Query
from pathlib import Path
import os
from fastapi.responses import JSONResponse
from fastapi import Request
from file_job.task import connectRabbitMQ

navigation_router = APIRouter(
    prefix="/api/navigation",
    tags=["navigation"],
)

BASE_DIR = Path(os.getenv("BASE_DIR", "/home"))


@navigation_router.get("/")
async def list_directory(request: Request):
    try:
        rel_path = request.query_params.get("path", "").lstrip("/")
        target_path = (BASE_DIR / rel_path).resolve()

        if not str(target_path).startswith(str(BASE_DIR)):
            return JSONResponse(status_code=403, content={"detail": "Доступ запрещён"})

        if not target_path.exists() or not target_path.is_dir():
            return JSONResponse(status_code=404, content={"detail": "Папка не найдена"})

        contents = []
        for item in target_path.iterdir():
            item_info = {
                "name": item.name,
                "is_dir": item.is_dir(),
                "is_file": item.is_file(),
                "size": item.stat().st_size if item.is_file() else None,
            }
            if item.is_dir():
                try:
                    item_info["children_count"] = len(list(item.iterdir()))
                except Exception:
                    item_info["children_count"] = None

            contents.append(item_info)

        return {
            "path": "/" + str(target_path.relative_to(BASE_DIR)),
            "items": contents
        }

    except Exception as e:
        return JSONResponse(status_code=500, content={"detail": f"Ошибка чтения директории: {e}"})
    

@navigation_router.post("/download")
async def download_file(file: File):
    connectRabbitMQ(str(file.path), str(file.name))
    print("Получен payload:", file)
    return {"message": "Задание на скачивание получено"}