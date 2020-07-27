import requests
from tests import delete_all_db, create_valid_account_info, HOST


def teardown_function():
    delete_all_db()


def test_sign_up_200():
    """
    Тест проверяет фунцию регистрации аккаунта пользователя.
    """
    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert r.status_code == 200, r.text


def test_sign_up_400_password_length_check_fails():
    """
    Тест проверяет фунцию регистрации аккаунта с нарушением минимальной длины пароля.
    """
    account = create_valid_account_info()
    account["Password"] = ""

    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account
    )

    assert r.status_code == 400, r.text
    assert "Password" in r.text


def test_sign_up_400_account_already_exists():
    """
    Тест проверяет фунцию регистрации аккаунта с уже существующими данными пользователя.
    """
    account = create_valid_account_info()

    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account
    )

    assert r.status_code == 200, r.text

    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account
    )

    assert r.status_code == 400, r.text
    assert "already" in r.text


def test_sign_up_400_name_is_empty():
    """
    Тест проверяет фунцию регистрации аккаунта без ввода имени пользователя.
    """
    account = create_valid_account_info()
    account["name"] = ""

    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account
    )

    assert r.status_code == 400, r.text
    assert "Name" in r.text


def test_sign_up_400_email_is_empty():
    """
    Тест проверяет фунцию регистрации аккаунта без ввода e-mail пользователя.
    """
    account = create_valid_account_info()
    account["Email"] = ""

    r = requests.post(
        f"{HOST}/accounts/sign-up",
        json=account
    )

    assert r.status_code == 400, r.text
    assert "Email" in r.text


def test_sign_out_204():
    """
    Тест проверяет фунцию выхода из аккаунта пользователя.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    sign_out_response = requests.delete(
        "http://localhost:3000/accounts/sign-out",
        headers=x_auth_token
    )

    assert sign_out_response.status_code == 204, sign_out_response.text


def test_sign_in_200():
    """
    Тест проверяет фунцию входа в аккаунт пользователя.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account = create_valid_account_info()
    del account["Name"]

    sign_in_response = requests.post(
        f"{HOST}/accounts/sign-in",
        json=account
    )

    assert sign_in_response.status_code == 200, sign_in_response.text


def test_sign_in_400_wrong_password():
    """
    Тест проверяет фунцию входа в аккаунт пользователя с вводом неверного пароля.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account = create_valid_account_info()
    account["Password"] = "2"

    sign_in_response = requests.post(
        f"{HOST}/accounts/sign-in",
        json=account
    )

    assert sign_in_response.status_code == 400, sign_in_response.text
    assert "password" in sign_in_response.text


def test_sign_in_400_account_does_not_exist():
    """
    Запрещено осуществлять вход в систему, не зарегистрировавшись.
    """

    r = requests.post(
        f"{HOST}/accounts/sign-in",
        json=create_valid_account_info()
    )
    assert r.status_code == 400, r.text
    assert "exist" in r.text


def test_account_info_200():
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
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
    """
    Тест проверяет фунцию удаления аккаунта пользователя.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    delete_account_response = requests.delete(
        f"{HOST}/accounts/{account_id}",
        headers=x_auth_token
    )

    assert delete_account_response.status_code == 204, delete_account_response.text


def test_sign_in_deleted_account():
    """
    Вход в удаленный аккаунт запрещен.
    """
    sign_up_response = requests.post(
        f"{HOST}/accounts/sign-up",
        json=create_valid_account_info()
    )

    assert sign_up_response.status_code == 200, sign_up_response.text

    account_id = sign_up_response.json()["Account"]["ID"]
    x_auth_token = {"X-Auth-Token": sign_up_response.json()["Token"]}

    delete_account_response = requests.delete(
        f"{HOST}/accounts/{account_id}",
        headers=x_auth_token
    )

    assert delete_account_response.status_code == 204, delete_account_response.text

    account = create_valid_account_info()
    del account["Name"]

    sign_in_response = requests.post(
        f"{HOST}/accounts/sign-in",
        json=account
    )

    assert sign_in_response.status_code == 400, sign_in_response.text
    assert "exist" in sign_in_response.text
