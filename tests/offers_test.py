import pytest
import requests
from tests import create_valid_account_info, create_valid_event_info, HOST, sign_up


def test_create_offer():
    """
    Тест проверяет функцию создания оффера.
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

    sign_up_response_2 = sign_up()

    assert sign_up_response_2.status_code == 200, sign_up_response_2.text

    account_id_2 = sign_up_response_2.json()["Account"]["ID"]
    account_2_token = {"X-Auth-Token": sign_up_response_2.json()["Token"]}

    create_offer_response = requests.post(
        f"{HOST}/offers",
        headers=account_2_token,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )

    assert create_offer_response.status_code == 200, create_offer_response.text


def test_create_offer_on_your_account_400():
    """
    Запрещено создавать оффер на собственный ивент.
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

    create_offer_response = requests.post(
        f"{HOST}/offers",
        headers=x_auth_token,
        json={
            "AccountID": account_id,
            "EventID": event_id
        }
    )

    assert create_offer_response.status_code == 400, create_offer_response.text
    assert "fails" in create_offer_response.text


@pytest.mark.xfail(reason="issue #37")
def test_repeated_create_offer_400():
    """
    По бизнес-логике запрещено создавать два оффера c одного аккаунта на один ивент.
    """

    sign_up_response = sign_up()

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_events_response = requests.post(
        f"{HOST}/events",
        headers=account_token,
        json=create_valid_account_info()
    )

    assert create_events_response.status_code == 200, create_events_response.text

    event_id = create_events_response.json()["ID"]

    sign_up_response_2 = sign_up()

    assert sign_up_response_2.status_code == 200, sign_up_response_2.text

    account_id_2 = sign_up_response_2.json()["Account"]["ID"]
    account_2_token = {"X-Auth-Token": sign_up_response_2.json()["Token"]}

    create_offer_response_1 = requests.post(
        f"{HOST}/offers",
        headers=account_2_token,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )

    assert create_offer_response_1.status_code == 200, create_offer_response_1.text

    create_offer_response_2 = requests.post(
        f"{HOST}/offers",
        headers=account_2_token,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )
    assert create_offer_response_2.status_code == 400, create_offer_response_2.text
    assert "fails" in create_offer_response_2.text


def test_create_offer_200_from_two_different_accounts():
    """
    Тест проверяет фунцию создания офферов с двух разных аккаунтов на один ивент.
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

    sign_up_response_2 = sign_up()

    assert sign_up_response_2.status_code == 200, sign_up_response_2.text

    account_id_2 = sign_up_response_2.json()["Account"]["ID"]
    x_auth_token_2 = {"X-Auth-Token": sign_up_response_2.json()["Token"]}

    create_offer_response = requests.post(
        f"{HOST}/offers",
        headers=x_auth_token_2,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )

    assert create_offer_response.status_code == 200, create_offer_response.text

    sign_up_response_3 = sign_up()

    assert sign_up_response_3.status_code == 200, sign_up_response_2.text

    account_id_3 = sign_up_response_3.json()["Account"]["ID"]
    x_auth_token_3 = {"X-Auth-Token": sign_up_response_3.json()["Token"]}

    create_offer_response = requests.post(
        f"{HOST}/offers",
        headers=x_auth_token_3,
        json={
            "AccountID": account_id_3,
            "EventID": event_id
        }
    )

    assert create_offer_response.status_code == 200, create_offer_response.text


def test_create_offer_400_account_that_created_the_event_deleted():
    """
    Запрещено создавать оффер на событие, созданное аккаунтом, который был удален.
    """

    sign_up_response = sign_up()

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    account_1_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)

    create_events_response = requests.post(
        f"{HOST}/events",
        headers=account_1_token,
        json=info
    )

    assert create_events_response.status_code == 200, create_events_response.text

    event_id = create_events_response.json()["ID"]

    delete_account_response = requests.delete(
        f"{HOST}/accounts/{account_id}",
        headers=account_1_token
    )

    assert delete_account_response.status_code == 204, delete_account_response.text

    sign_up_response_2 = sign_up()

    assert sign_up_response_2.status_code == 200, sign_up_response_2.text

    account_id_2 = sign_up_response_2.json()["Account"]["ID"]
    account_2_token = {"X-Auth-Token": sign_up_response_2.json()["Token"]}

    create_offer_response = requests.post(
        f"{HOST}/offers",
        headers=account_2_token,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )

    assert create_offer_response.status_code == 400, create_offer_response.text


def test_create_offer_400_id_doesnt_match():
    """
    Запрещено создавать оффер без подтверждения ID пользователя.
    """
    sign_up_response = sign_up()

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info(account_id)

    create_events_response = requests.post(
        f"{HOST}/events",
        headers=x_auth_token,
        json=info
    )

    assert create_events_response.status_code == 200, create_events_response.text

    event_id = create_events_response.json()["ID"]

    info['OwnerID'] = account_id + 1

    create_offer_response = requests.post(
        f"{HOST}/offers",
        headers=x_auth_token,
        json={
            "AccountID": info["OwnerID"],
            "EventID": event_id
        }
    )

    assert create_offer_response.status_code == 400, create_offer_response.text
    assert "this account" in create_offer_response.text
