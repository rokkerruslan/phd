import os
from dotenv import load_dotenv
from pathlib import Path
from mimesis import Person
from mimesis import Text
from datetime import datetime, timedelta
import requests


person = Person("ru")
text = Text("ru")
env_path = Path(".env.tests")
load_dotenv(dotenv_path=env_path)

HOST = f"http://localhost{os.environ.get('ADDR')}"

format_time = "%Y-%m-%dT%H:%M:%SZ"


def sign_up():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    return sign_up_response


def create_valid_account_info():
    """
    Функция создания словаря с валидными данными аккаунта.
    """

    return {
        "Email": person.email(),
        "Password": person.password(length=10),
        "Name": person.name()
    }


def create_valid_event_info(account_id):
    """
    Функция создания словаря с валидными данными ивента.
    """

    start = datetime.now() + timedelta(minutes=65)
    end = datetime.now() + timedelta(minutes=165)

    return {
        "Name": text.title(),
        "Description": text.text(quantity=1),
        "OwnerID": account_id,
        "IsPublic": False,
        "isHidden": False,
        "Timelines": [
            {
                "Start": start.strftime(format_time),
                "End": end.strftime(format_time),
                "Place": "Saint Petersburg"
            }
        ]
    }


def create_invalid_event_info(account_id):
    """
    Функция создания словаря с валидными данными ивента.
    """
    return {
        "Name": text.title(),
        "Description": text.text(quantity=1),
        "OwnerID": account_id,
        "IsPublic": False,
        "isHidden": False,
        "Timelines": [
            {
                "Start": "2019-01-02T17:05:05Z",
                "End": "2019-01-02T18:06:05Z",
                "Place": "Saint Petersburg"
            }
        ]
    }
