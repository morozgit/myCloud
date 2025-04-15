import uvicorn
from fastapi import FastAPI
from src.navigation import navigation_router
app = FastAPI(
    title="Backend MyCloud",
    openapi_url="/api/openapi.json",
    docs_url="/api/docs"
)

app.include_router(navigation_router)

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)