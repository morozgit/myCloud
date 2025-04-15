from fastapi import APIRouter
import os, pwd
from pathlib import Path
from fastapi import FastAPI, HTTPException

navigation_router = APIRouter(
    prefix="/navigation",
    tags=["navigation"],
)


BASE_DIR = Path(os.getenv("BASE_DIR", "/home")).resolve()
    
@navigation_router.get("/")
async def list_home_directory():
    if not BASE_DIR.exists() or not BASE_DIR.is_dir():
        raise HTTPException(status_code=404, detail="Базовая директория не найдена")

    contents = []
    for item in BASE_DIR.iterdir():
        contents.append({
            "name": item.name,
            "is_dir": item.is_dir(),
            "is_file": item.is_file(),
            "size": item.stat().st_size if item.is_file() else None,
        })

    return {
        "path": str(BASE_DIR),
        "items": contents
    }