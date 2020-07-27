import pytest
import requests
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

    info = create_valid_event_info()
    info['OwnerID'] = account_id
    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text


def test_create_event_400_title_cant_be_empty():
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
    info = create_valid_event_info()
    info['OwnerID'] = account_id
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
    info = create_valid_event_info()
    info['OwnerID'] = account_id
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
    info = create_valid_event_info()
    info['OwnerID'] = account_id
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
    Тест проверяет функцию создания ивентов не авторизированным пользователем.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info()
    info['OwnerID'] = account_id + 1

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "authorized" in create_event_response.text


@pytest.mark.xfail(reason="issue #40")
def test_update_events_200():
    """
    Тест проверяет функцию обновления данных ивента с валидными данными.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info()
    info['OwnerID'] = account_id

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text

    event_id = create_event_response.json()["ID"]

    event = create_valid_event_info()
    del event["OwnerID"]

    update_event_response = requests.put(
        f"{HOST}/events{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 200, update_event_response.text


@pytest.mark.xfail(reason="issue #41")
def test_event_info():
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
    info = create_valid_event_info()
    info['OwnerID'] = account_id

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text

    event_id = create_event_response.json()["ID"]
    event_info_response = requests.get(f"{HOST}/events/{event_id}")

    assert event_info_response.status_code == 200, event_info_response.text


def test_list_events():
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
    info = create_valid_event_info()
    info['OwnerID'] = account_id

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text

    list_events_response = requests.get(f"{HOST}/events", params={"account_id": account_id})

    assert list_events_response.status_code == 200, list_events_response.text
    assert type(list_events_response.json()) == list
