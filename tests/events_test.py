import pytest
import requests
from datetime import datetime, timedelta
from tests import delete_all_db, create_valid_account_info, create_valid_event_info, create_invalid_event_info, HOST, \
    format_time


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


def test_create_event_400_id_doesnt_match():
    """
    Запрещено создавать ивент без подтверждения ID пользователя.
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


def test_create_event_400_not_authorized():
    """
    Тест проверяет функцию создания ивентов не авторизованным пользователем.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    info = create_valid_event_info(account_id)

    create_event_response = requests.post(
        f"{HOST}/events",
        json=info
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "X-Auth-Token" in create_event_response.text


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
    assert event_info_response.status_code == 200, event_info_response.text

    event_info = event_info_response.json()
    timlelines = event_info["Timelines"]

    assert timlelines


def test_create_events_400_start_after_end():
    """
    Запрещено задавать конец ивента раньше старта
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


def test_create_events_400_timeline_is_not_in_the_past():
    """
    Тест проверяет возниконовение ошибки при создании ивента с таймлайном который находится в прошедшем времени.
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
        json=create_invalid_event_info(account_id)
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "early" in create_event_response.text


@pytest.mark.xfail(reason="issue #58")
def test_create_events_400_intersection_of_timelines():
    """
    Тест проверяет возниконовение ошибки при пересечении временных промежутков.
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
            "Start": "2021-01-02T17:05:05Z",
            "End": "2021-01-02T18:06:05Z",
            "Place": "Saint Petersburg"
        }
    )

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=event
    )

    assert create_event_response.status_code == 400, create_event_response.text
    assert "crossing" in create_event_response.text


@pytest.mark.xfail
def test_create_event_400_timeline_is_not_in_the_present():
    """
    Таймлайн в ивенте разрешается устанавливать на 1 час позже чем время создания ивента. Тест проверяет
    возниконовение ошибки при создании ивента с таймлайном до этого часа.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    info = create_valid_event_info(account_id)
    forbidden_time = datetime.now() + timedelta(minutes=55)

    info["Timelines"][0]["Start"] = forbidden_time.strftime(format_time)

    create_event_response_1 = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response_1.status_code == 400, create_event_response_1.text
    assert "early" in create_event_response_1.text

    forbidden_time = datetime.now() + timedelta(minutes=65)

    info["Timelines"][0]["Start"] = forbidden_time.strftime(format_time)

    create_event_response_2 = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response_2.status_code == 200, create_event_response_2.text
