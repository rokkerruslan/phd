import psycopg2
import os
from dotenv import load_dotenv
from pathlib import Path


env_path = Path(".env.tests")
load_dotenv(dotenv_path=env_path)

HOST = f"http://localhost{os.environ.get('ADDR')}"


def delete_all_db():
    """
    Фунция очистки всех баз данных.
    """

    connect = psycopg2.connect(os.environ.get('DATABASE_URL'))
    cursor = connect.cursor()

    query = """
        delete from tokens;
        delete from offers;
        delete from images;
        delete from timelines;
        delete from events;
        delete from accounts;
    """

    cursor.execute(query)

    connect.commit()
    cursor.close()
    connect.close()


def create_valid_account_info():
    """
    Функция создания словаря с валидными данными аккаунта.
    """
    return {
        "Email": "test1@gmail.com",
        "Password": "1",
        "Name": "test1"
    }


def create_valid_event_info():
    """
    Функция создания словаря с валидными данными ивента.
    """
    return {
        "Name": "1",
        "Description": "1",
        "OwnerID": 0,
        "Timelines": [
            {
                "Start": "2006-01-02T17:05:05Z",
                "End": "2006-01-02T18:06:05Z",
                "Place": "Saint Petersburg"
            }
        ]
    }
