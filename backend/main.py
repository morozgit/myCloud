import uvicorn
from fastapi import FastAPI
from src.navigation import navigation_router
from src.files import files_router
from fastapi.middleware.cors import CORSMiddleware


app = FastAPI(
    title="Backend MyCloud",
    openapi_url="/api/openapi.json",
    docs_url="/api/docs"
)

app.include_router(navigation_router)
app.include_router(files_router)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["Authorization", "Content-Type", "X-Requested-With"],
)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)