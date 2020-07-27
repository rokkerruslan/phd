import pytest
import requests
from tests import delete_all_db, create_valid_account_info, create_valid_event_info, HOST


def teardown_function():
    delete_all_db()


def test_create_offer_on_your_account():
    """
    Запрещенно создавать офер на собствееный ивент.
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
    По бизнес-логике запрещено создавать два офера на один ивент.
    """

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    account = create_valid_account_info()
    account["Email"] += "test"

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id_2 = sign_up_response.json()["Account"]["ID"]
    account_2_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_events_response = requests.post(
        f"{HOST}/events",
        headers=account_token,
        json=create_valid_account_info()
    )

    assert create_events_response.status_code == 200, create_events_response.text

    event_id = create_events_response.json()["ID"]

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


@pytest.mark.xfail(reason="issue #42")
def test_create_offer_400():
    """
    Запрещено создавать офер на событие, созданное аккаунтом, который был удален.
    """

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    account_1_token = {"X-Auth-Token": sign_up_response.json()["Token"]}
    info = create_valid_event_info()
    info['OwnerID'] = account_id

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

    account = create_valid_account_info()
    account["Email"] += "test"

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id_2 = sign_up_response.json()["Account"]["ID"]
    account_2_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_offer_response = requests.post(
        f"{HOST}/offers",
        headers=account_2_token,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )

    assert create_offer_response.status_code == 400, create_offer_response.text
