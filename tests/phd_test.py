import pytest
import requests
import psycopg2


def find_test(variables):
    """
    :param dict variables:
    """

    functions = []

    for name in variables:
        function = variables[name]
        if name.startswith("test_") and callable(function):
            functions.append(function)

    return functions


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


def test_sign_up_200ok():
    delete_all_db()

    r = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert r.status_code == 200, r.text

    delete_all_db()


def test_sign_up_400_password_length_check_fails():
    delete_all_db()

    r = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "",
            "name": "test1"
        }
    )
    assert r.status_code == 400, r.text

    delete_all_db()


def test_sign_up_400br_account_already_exists():
    delete_all_db()

    r = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert r.status_code == 200, r.text

    r = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert r.status_code == 400, r.text

    delete_all_db()


def test_sign_out_204ok():
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    sign_out_response = requests.delete(
        "http://localhost:3000/accounts/sign-out",
        headers=x_auth_token
    )

    assert sign_out_response.status_code == 204

    delete_all_db()


def test_sign_in_200ok():
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200

    sign_in_response = requests.post(
        "http://localhost:3000/accounts/sign-in",
        json={"Email": "test1@gmail.com",
              "Password": "1"
              }
    )

    assert sign_in_response.status_code == 200

    delete_all_db()


def test_sign_in_400br_wrong_password():
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200

    sign_in_response = requests.post(
        "http://localhost:3000/accounts/sign-in",
        json={"Email": "test1@gmail.com",
              "Password": "2"
              }
    )

    assert sign_in_response.status_code == 400

    delete_all_db()


def test_sign_in_400br_account_does_not_exist():
    """

    Запрещено осуществлять вход в систему, не зарегистрировавшись.
    """
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200

    sign_in_response = requests.post(
        "http://localhost:3000/accounts/sign-in",
        json={"Email": "test2@gmail.com",
              "Password": "2"
              }
    )
    assert sign_in_response.status_code == 400

    delete_all_db()


def test_info_200ok():
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    info_response = requests.get(
        f"http://localhost:3000/accounts/{account_id}",
        headers=x_auth_token
    )

    assert info_response.status_code == 200

    delete_all_db()


def test_delete_account_200ok():
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    delete_account_response = requests.delete(
        f"http://localhost:3000/accounts/{account_id}",
        headers=x_auth_token
    )

    assert delete_account_response.status_code == 204

    delete_all_db()


def create_event_200ok():
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_event_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
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

    assert create_event_response.status_code == 200


@pytest.mark.xfail(reason="issue #37")
def test_create_offer_400br():
    """

    По бизнес-логике запрещено создавать два офера на один ивент.
    """
    delete_all_db()

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test1@gmail.com",
            "Password": "1",
            "name": "test1"
        }
    )

    assert sign_up_response.status_code == 200

    account_id_1 = sign_up_response.json()["Account"]["ID"]
    account_1_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    sign_up_response = requests.post(
        "http://localhost:3000/accounts/sign-up",
        json={
            "Email": "test2@gmail.com",
            "Password": "2",
            "name": "test2"
        }
    )

    assert sign_up_response.status_code == 200

    account_id_2 = sign_up_response.json()["Account"]["ID"]
    account_2_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    create_events_response = requests.post(
        "http://localhost:3000/events",
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

    assert create_events_response.status_code == 200

    event_id = create_events_response.json()["ID"]

    create_offer_response_1 = requests.post(
        "http://localhost:3000/offers",
        headers=account_2_token,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )

    assert create_offer_response_1.status_code == 200

    create_offer_response_2 = requests.post(
        "http://localhost:3000/offers",
        headers=account_2_token,
        json={
            "AccountID": account_id_2,
            "EventID": event_id
        }
    )
    assert create_offer_response_2.status_code == 400

    delete_all_db()


if __name__ == '__main__':

    find_test_list = find_test(locals().copy())
    for test in find_test_list:
        test()
