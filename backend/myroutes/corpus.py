import os
from fastapi import APIRouter, HTTPException, status

router = APIRouter(
        prefix="/corpus",
        tags=["corpus"]
)

@router.get("/text", resp)
