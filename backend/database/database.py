from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker 


engine = create_engine("postgresql+psycopg2://postgres:123@103.68.108.252:5432/postgres")
SessionLocal = sessionmaker(autoflush=False, bind=engine)


def get_db():
    db_session = SessionLocal()
    return db_session

def get_engine():
    return engine

