import pytest
import requests
import psycopg2

HOST = "http://localhost:3000"


def teardown_function():
    delete_all_db()


def delete_all_db():
    connect = psycopg2.connect(dbname="postgres",
                               user="postgres",
                               password="postgres",
                               host='localhost')
    cursor = connect.cursor()

    cursor.execute('delete from tokens;')
    cursor.execute('delete from offers;')
    cursor.execute('delete from images;')
    cursor.execute('delete from timelines;')
    cursor.execute('delete from events;')
    cursor.execute('delete from accounts;')

    connect.commit()
    cursor.close()
    connect.close()


delete_all_db()

account_info = {
    "Email": "test1@gmail.com",
    "Password": "1",
    "Name": "test1"
}


def test_sign_up_200():
    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert r.status_code == 200, r.text


def test_sign_up_400_password_length_check_fails():
    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "",
            "name": "test1"
        }
    )

    assert r.status_code == 400, r.text
    assert "Password" in r.text


def test_sign_up_400_account_already_exists():
    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert r.status_code == 200, r.text

    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert r.status_code == 400, r.text
    assert "already" in r.text


def test_sign_out_204():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    sign_out_response = requests.delete(
        "http://localhost:3000/accounts/sign-out",
        headers=x_auth_token
    )

    assert sign_out_response.status_code == 204, sign_out_response.text


def test_sign_in_200():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    sign_in_response = requests.post(
        f"{HOST}/accounts/sign-in",
        json=account_info
    )

    assert sign_in_response.status_code == 200, sign_in_response.text


def test_sign_in_400_wrong_password():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    sign_in_response = requests.post(
        f"{HOST}/accounts/sign-in",
        json={
            "Email": "test1@gmail.com",
            "Password": "2",
            "name": "test1"
        }
    )

    assert sign_in_response.status_code == 400, sign_in_response.text
    assert "password" in sign_in_response.text


def test_sign_in_400_account_does_not_exist():
    """

    Запрещено осуществлять вход в систему, не зарегистрировавшись.
    """

    r = requests.post(
        f"{HOST}/accounts/sign-in",
        json=account_info
    )
    assert r.status_code == 400, r.text
    assert "exist" in r.text


def test_account_info_200():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    info_response = requests.get(
        f"{HOST}/accounts/{account_id}",
        headers=x_auth_token
    )

    assert info_response.status_code == 200, info_response.text


def test_delete_account_200():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    delete_account_response = requests.delete(
        f"{HOST}/accounts/{account_id}",
        headers=x_auth_token
    )

    assert delete_account_response.status_code == 204, delete_account_response.text


def create_event_200():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_event_response = requests.post(
        f"{HOST}/accounts/sign-up",
        headers=x_auth_token,
        json={
            "Name": "1",
            "Description": "1",
            "OwnerID": account_id,
            "Timelines": [
                {
                    "Start": "2006-01-02T17:05:05Z",
                    "End": "2006-01-02T18:06:05Z",
                    "Place": "Saint Petersburg"
                }
            ]
        }
    )

    assert create_event_response.status_code == 200, create_event_response.text


@pytest.mark.xfail(reason="issue #37")
def test_create_offer_400():
    """

    По бизнес-логике запрещено создавать два офера на один ивент.
    """

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account_info
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id_1 = sign_up_response.json()["Account"]["ID"]
    account_1_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json={
            "Email": "test2@gmail.com",
            "Password": "2",
            "name": "test2"
        }
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id_2 = sign_up_response.json()["Account"]["ID"]
    account_2_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_events_response = requests.post(
        f"{HOST}/events",
        headers=account_1_token,
        json={
            "Name": "1",
            "Description": "1",
            "OwnerID": account_id_1,
            "Timelines": [
                {
                    "Start": "2006-01-02T17:05:05Z",
                    "End": "2006-01-02T18:06:05Z",
                    "Place": "Saint Petersburg"
                }
            ]
        }
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
