import pytest
import requests
from tests import create_valid_account_info, create_valid_event_info, HOST, sign_up


def test_update_events_200():
    """
    Тест проверяет функцию обновления данных ивента с валидными данными.
    """
    sign_up_response = sign_up()

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


def test_update_events_400_name_is_empty():
    """
    Тест проверяет возникновение ошибки при обновлении ивента без ввода имени ивента.
    """
    sign_up_response = sign_up()

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
    sign_up_response = sign_up()

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
    sign_up_response = sign_up()

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


@pytest.mark.xfail(reason="issue #53")
def test_update_events_400_not_authorized():
    """
    Запрещено обновлять ивент без авторизации. issues 53
    """
    sign_up_response = sign_up()

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
    assert "X-Auth-Token" in create_event_response.text


@pytest.mark.xfail(reason="issue #54")
def test_update_events_400_id_doesnt_match():
    """
    Запрещено обновлять ивент без подверждения ID пользователя. issues 54
    """
    sign_up_response = sign_up()

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
    sign_up_response = sign_up()

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


@pytest.mark.xfail(reason="issue #56")
def test_update_events_info():
    """
    Тест проверяет обновление информации ивента. issue #56
    """
    sign_up_response = sign_up()

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
    sign_up_response = sign_up()

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
    event["Timelines"][0]["Start"] = "2022-08-02T13:43:09.535504Z"
    event["Timelines"][0]["End"] = "2022-08-03T13:43:09.535504Z"
    event["Timelines"][0]["Place"] = "Moscow"

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )

    assert update_event_response.status_code == 200, update_event_response.text

    event_info_response = requests.get(f"{HOST}/events/{event_id}")
    assert event_info_response.status_code == 200, event_info_response.text

    event_info = event_info_response.json()
    assert event_info["Timelines"] == event["Timelines"]


@pytest.mark.xfail
def test_update_events_added_two_timelines():
    """
    Тест проверяет добавление дополнительного таймлана в ивент.
    """
    sign_up_response = sign_up()

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


@pytest.mark.xfail
def test_update_events_del_one_timelines():
    """
    Тест проверяет удаление одного из таймлайнов ивента.
    """
    sign_up_response = sign_up()

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

    assert create_event_response.status_code == 200, create_event_response.text

    event_id = create_event_response.json()["ID"]
    del event["Timelines"][1]

    update_event_response = requests.put(
        f"{HOST}/events/{event_id}",
        headers=x_auth_token,
        json=event
    )
    assert update_event_response.status_code == 200, update_event_response.text

    event_info = update_event_response.json()

    assert len(event_info["Timelines"]) == 1
