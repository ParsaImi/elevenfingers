from fastapi import FastAPI
from routes import auth
import os
from dotenv import load_dotenv
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()
load_dotenv()


origins = [
    "http://localhost",
    "http://localhost:5173",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
class Settings:
    SECRET_KEY: str = os.getenv('SECRET_KEY' , "fallback_secret_key")

settings = Settings()


app.include_router(auth.router)

@app.get("/")
async def root():
    return {"message": "Hellow nig"}
