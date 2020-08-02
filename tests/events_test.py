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
    info = create_valid_event_info(account_id)

    create_event_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_event_response.status_code == 200, create_event_response.text

    event_id = create_event_response.json()["ID"]

    event = create_valid_event_info(account_id)
    del event["OwnerID"]

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 200, update_event_response.text


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


def test_update_events_400_name_is_empty():
    """
    Тест проверяет возникновение ошибки при обновлении ивента без ввода имени ивента.
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

    event = create_valid_event_info(account_id)
    event["Name"] = ""
    del event["OwnerID"]

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 400, update_event_response.text
    assert "Title" in update_event_response.text


def test_update_events_400_description_is_empty():
    """
    Тест проверяет возникновение ошибки при обновлении ивента без ввода описания ивента.
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

    event = create_valid_event_info(account_id)
    event["Description"] = ""
    del event["OwnerID"]

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 400, update_event_response.text
    assert "Description" in update_event_response.text


@pytest.mark.xfail(reason="issue #52")
def test_update_events_400_timelines_is_empty():
    """
    Тест проверяет возникновение ошибки при обновлении ивента без ввода таймлайна ивента. issue #52
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

    event = create_valid_event_info(account_id)
    event["Timeline"] = []
    del event["OwnerID"]

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 400, update_event_response.text
    assert "Timeline" in update_event_response.text


@pytest.mark.xfail
def test_update_events_400_not_authorized():
    """
    Запрещено обновлять ивент без авторизации.
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

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        json=info
    )

    assert update_event_response.status_code == 400, update_event_response.text


@pytest.mark.xfail
def test_update_events_400_id_doesnt_match():
    """
    Запрещено обновлять ивент без авторизации.
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
    info['OwnerID'] = account_id + 1

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=info
    )

    assert update_event_response.status_code == 400, update_event_response.text
    assert "authorized" in update_event_response.text


@pytest.mark.xfail
def test_update_events_owner_id_not_nil():
    """
    ID аккаунта не может ровняться нулю.
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

    event = create_valid_event_info(account_id)
    del event["OwnerID"]

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 200, update_event_response.text
    assert update_event_response.json()["OwnerID"] != 0


@pytest.mark.xfail
def test_update_events_info():
    """
    Тест проверяет обновление информации ивента.
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

    event = create_valid_event_info(account_id)
    event["Name"] += "test"
    event["Description"] += "test"
    event["IsPublic"] = True
    event["IsHidden"] = True

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 200, update_event_response.text

    event_info_response = requests.get(f"{HOST}/events/{event_id}")
    event_info = event_info_response.json()

    assert event_info_response.status_code == 200, event_info_response.text
    assert event_info["Name"] == event["Name"]
    assert event_info["Description"] == event["Description"]
    assert event_info["IsPublic"] == event["IsPublic"]
    assert event_info["IsHidden"] == event["IsHidden"]


@pytest.mark.xfail
def test_update_events_timelines():
    """
    Тест проверяет обновление данных в таймлайне ивента.
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

    event = create_valid_event_info(account_id)
    event["Timelines"][0]["Start"] = "2021-08-02T13:43:09.535504Z"
    event["Timelines"][0]["End"] = "2021-08-03T13:43:09.535504Z"
    event["Timelines"][0]["Place"] = "Moscow"

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 200, update_event_response.text

    event_info_response = requests.get(f"{HOST}/events/{event_id}")
    event_info = event_info_response.json()

    assert event_info_response.status_code == 200, event_info_response.text
    assert event_info["Timelines"] == event["Timelines"]


@pytest.mark.xfail
def test_update_events_added_two_timelines():
    """
    Тест проверяет добавление дополнительного таймлана в ивент.
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

    timelines_2 = {
                "Start": "2021-01-02T17:05:05Z",
                "End": "2021-01-02T18:06:05Z",
                "Place": "Saint Petersburg"
            }

    event = create_valid_event_info(account_id)
    event["Timelines"].append(timelines_2)
    del event["OwnerID"]

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    event_info = update_event_response.json()

    assert update_event_response.status_code == 200, update_event_response.text
    assert len(event_info["Timelines"]) == 2
