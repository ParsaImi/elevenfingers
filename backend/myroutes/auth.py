import os
from fastapi import security
import jwt
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from datetime import timedelta, datetime
from fastapi import APIRouter, HTTPException, status
from starlette.status import HTTP_201_CREATED
from database.database import get_db
from models.schemas import UserCreate, User, UserResponse
from passlib.context import CryptContext
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from dotenv import load_dotenv
from fastapi import Depends

load_dotenv()
db = get_db()

router = APIRouter(
        prefix="/auth",
        tags=["auth"]
)

# Password Hashing 

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")
ACCESS_TOKEN_EXPIRE_MINUTES = 30
ALGORITHM = "HS256"



# Token model
class Token(BaseModel):
    access_token: str
    token_type: str

@router.post("/signup" , response_model=UserResponse, status_code=HTTP_201_CREATED)
def signup(user: UserCreate):
    if get_user_by_email(user.email):
        raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Email already in use"
        )
    if get_user_by_username(user.username):
        raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Username already in use"
    )
    db_user = User(
            email = user.email,
            username = user.username,
            hashed_password=get_hashed_password(user.password)
    )

    db.add(db_user)
    db.commit()
    db.refresh(db_user)

    return db_user


@router.post("/verify")
def verify(credentials : HTTPAuthorizationCredentials):
    try:
        token = credentials.credentials
        payload = jwt.decode(token, os.getenv("SECRET_KEY"), algorithms=ALGORITHM)
        return {"verify" : True , "username": payload.get("sub")}
    except ExpiredSignatureError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token has expired",
            headers={"WWW-Authenticate": "Bearer"},
        )
    except jwt.PyJWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )

@router.post("/login", response_model=Token)
def login(form_data : OAuth2PasswordRequestForm = Depends()):
    user = get_user_by_email(form_data.username)
    if not user :
        user = get_user_by_username(form_data.username)
    if not user or not verify_password(form_data.password, user.hashed_password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username/email or password",
            headers={"WWW-Authenticate": "Bearer"},
        )
    access_token_expires = timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    access_token = create_access_token(
            data={"sub": user.username},
            expires_delta=access_token_expires
    ) 
    return {"access_token": access_token, "token_type": "bearer"}


def create_access_token(data: dict , expires_delta):
    to_encode = data.copy()
    if expires_delta:
        expire = datetime.now() + expires_delta
    else:
        expire = datetime.now() + timedelta(minutes=15)

    to_encode.update({"exp" : expire})
    encoded_jwt = jwt.encode(to_encode, os.getenv("SECRET_KEY"), algorithm=ALGORITHM)
    return encoded_jwt

def get_hashed_password(password: str):
    return pwd_context.hash(password)


def get_user_by_email(email: str):
    return db.query(User).filter(User.email == email).first()


def get_user_by_username(username: str):
    return db.query(User).filter(User.username == username).first()

def verify_password(plain_password, hashed_password):
    return pwd_context.verify(plain_password, hashed_password)
