from fastapi import APIRouter, HTTPException, Query
from pathlib import Path
import os

navigation_router = APIRouter(
    prefix="/api/navigation",
    tags=["navigation"],
)

BASE_DIR = Path(os.getenv("BASE_DIR", "/home")).resolve()

@navigation_router.get("/")
async def list_directory(path: str = Query("/")):
    requested_path = (BASE_DIR / path.lstrip("/")).resolve()

    # Защита от выхода за пределы BASE_DIR
    if not str(requested_path).startswith(str(BASE_DIR)):
        raise HTTPException(status_code=403, detail="Доступ к этой директории запрещён")

    if not requested_path.exists() or not requested_path.is_dir():
        raise HTTPException(status_code=404, detail="Директория не найдена")

    contents = []
    for item in requested_path.iterdir():
        contents.append({
            "name": item.name,
            "is_dir": item.is_dir(),
            "is_file": item.is_file(),
            "size": item.stat().st_size if item.is_file() else None,
        })

    return {
        "path": "/" + str(requested_path.relative_to(BASE_DIR)),
        "items": contents
    }
