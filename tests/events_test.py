import pytest
import requests
from datetime import datetime
from tests import delete_all_db, create_valid_account_info, create_valid_event_info, HOST


def teardown_function():
    delete_all_db()


def test_create_event_200():
    """
    Тест проверяет функцию создания ивентов с валидными данными.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    info = create_valid_event_info(account_id)
    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text


def test_create_event_400_name_cant_be_empty():
    """
    Тест проверяет функцию создания ивентов с невалидными данными (Отсутствует значение ключа "Name" в теле запроса).
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)
    info["Name"] = ""

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "Title" in create_event_response.text


def test_create_event_400_description_cant_be_empty():
    """
    Тест проверяет функцию создания ивентов с невалидными данными (Отсутствует значение ключа "Description" в теле
    запроса).
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)
    info["Description"] = ""

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "Description" in create_event_response.text


def test_create_event_400_timlines_cant_be_empty():
    """
    Тест проверяет функцию создания ивентов с невалидными данными (Отсутствует значение ключа "Timelines" в теле
    запроса).
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)
    info["Timelines"] = []

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "Timelines" in create_event_response.text


def test_create_event_400_didnt_authorized():
    """
    Тест проверяет функцию создания ивентов не авторизованным пользователем.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)
    info['OwnerID'] = account_id + 1

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "authorized" in create_event_response.text


def test_event_info_200():
    """
    Тест проверяет функцию получения информации об ивенте.
    """

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text

    event_id = create_event_response.json()["ID"]
    event_info_response = requests.get(f"{HOST}/events/{event_id}")

    assert event_info_response.status_code == 200, event_info_response.text


def test_list_events_200():
    """
    Тест проверяет функцию получения списка ифентов у одного аккаунта.
    """

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text

    list_events_response = requests.get(f"{HOST}/events", params={"account_id": account_id})

    assert list_events_response.status_code == 200, list_events_response.text
    assert type(list_events_response.json()) == list


@pytest.mark.xfail(reason="issue #41")
def test_check_timlines_not_nil():
    """
    Таймлайн ивента не должен быть нулём. Issue #41.
    """

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text

    event_id = create_event_response.json()["ID"]

    event_info_response = requests.get(f"{HOST}/events/{event_id}")

    event_info = event_info_response.json()
    timlelines = event_info["Timelines"]

    assert event_info_response.status_code == 200, event_info_response.text
    assert timlelines


def test_create_events_400_start_after_end():
    """

    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json={
            "Name": "1",
            "Description": "1",
            "OwnerID": account_id,
            "IsPublic": False,
            "isHidden": False,
            "Timelines": [
                {
                    "Start": "2006-01-02T17:05:05Z",
                    "End": "2005-01-02T18:06:05Z",
                    "Place": "Saint Petersburg"
                }
            ]
        }
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "Start" in create_event_response.text


@pytest.mark.xfail(reason="issue #57")
def test_create_events_400_timeline_is_not_in_the_past():
    """
    Тест проверяет возниконовение ошибки при создании ивента с таймлайном который находится в прошедшем
    промежутке времени.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )
    event = create_event_response.json()

    created = datetime.strptime(event["Created"], "%Y-%m-%dT%H:%M:%S.%fZ")
    start = datetime.strptime(event["Timelines"][0]["Start"], "%Y-%m-%dT%H:%M:%SZ")
    end = datetime.strptime(event["Timelines"][0]["End"], "%Y-%m-%dT%H:%M:%SZ")

    assert create_event_response.status_code == 400, create_event_response.text
    assert created > start and end


@pytest.mark.xfail()
def test_create_events_400_intersection_of_timelines():
    """
    Тест проверяет возниконовение ошибки при пересечении временнх промежутков.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    event = create_valid_event_info(account_id)
    event["Timelines"].append(
        {
            "Start": "2018-01-02T17:05:05Z",
            "End": "2018-01-02T18:06:05Z",
            "Place": "Saint Petersburg"
        }
    )

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=event
    )

    create_event = create_event_response.json()

    start_1 = datetime.strptime(create_event["Timelines"][0]["Start"], "%Y-%m-%dT%H:%M:%SZ")
    end_1 = datetime.strptime(create_event["Timelines"][0]["End"], "%Y-%m-%dT%H:%M:%SZ")
    start_2 = datetime.strptime(create_event["Timelines"][1]["Start"], "%Y-%m-%dT%H:%M:%SZ")
    end_2 = datetime.strptime(create_event["Timelines"][1]["End"], "%Y-%m-%dT%H:%M:%SZ")

    assert create_event_response.status_code == 400, create_event_response.text
    assert start_1 and end_1 > end_2 or start_1 and end_1 < start_2
