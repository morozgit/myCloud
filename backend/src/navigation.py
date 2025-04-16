from fastapi import APIRouter, HTTPException, Query
from pathlib import Path
import os
from fastapi.responses import JSONResponse
from fastapi import Request

navigation_router = APIRouter(
    prefix="/api/navigation",
    tags=["navigation"],
)

BASE_DIR = Path(os.getenv("BASE_DIR", "/home")).resolve()


@navigation_router.get("/")
async def list_directory(request: Request):
    try:
        # Получаем относительный путь из query параметра
        rel_path = request.query_params.get("path", "").lstrip("/")
        target_path = (BASE_DIR / rel_path).resolve()

        # Проверка безопасности — чтобы не выйти за пределы BASE_DIR
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

            # Если это папка — посчитаем количество вложенных элементов
            if item.is_dir():
                try:
                    item_info["children_count"] = len(list(item.iterdir()))
                except Exception:
                    item_info["children_count"] = None  # вдруг нет доступа

            contents.append(item_info)

        return {
            "path": str(target_path.relative_to(BASE_DIR)),
            "items": contents
        }

    except Exception as e:
        return JSONResponse(status_code=500, content={"detail": f"Ошибка чтения директории: {e}"})