from fastapi import APIRouter
import os
from pathlib import Path
from fastapi import FastAPI, HTTPException

navigation_router = APIRouter(
    prefix="/api/navigation",
    tags=["navigation"],
)

@navigation_router.get("/")
def list_home_directory():
    username = os.getlogin()
    home_dir = Path(f"/home/{username}").resolve()

    if not home_dir.exists() or not home_dir.is_dir():
        raise HTTPException(status_code=404, detail="Домашняя директория не найдена")

    contents = []
    for item in home_dir.iterdir():
        contents.append({
            "name": item.name,
            "is_dir": item.is_dir(),
            "is_file": item.is_file(),
            "size": item.stat().st_size if item.is_file() else None,
        })

    return {
        "username": username,
        "path": str(home_dir),
        "items": contents
    }  