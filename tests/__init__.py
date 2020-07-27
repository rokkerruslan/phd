import psycopg2

HOST = "http://localhost:3000"


def delete_all_db():
    """
    Фунция очистки всех баз данных.
    """
    connect = psycopg2.connect(dbname="postgres",
                               user="postgres",
                               password="postgres",
                               host='localhost')
    cursor = connect.cursor()

    query = """
        delete from tokens;
        delete from offers;
        delete from images;
        delete from timelines;
        delete from events;
        delete from accounts;
    """

    cursor.execute(query)

    connect.commit()
    cursor.close()
    connect.close()


def create_valid_account_info():
    """
    Функция создания словаря с валидными данными аккаунта.
    """
    return {
        "Email": "test1@gmail.com",
        "Password": "1",
        "Name": "test1"
    }


def create_valid_event_info():
    """
    Функция создания словаря с валидными данными ивента.
    """
    return {
        "Name": "1",
        "Description": "1",
        "OwnerID": 0,
        "Timelines": [
            {
                "Start": "2006-01-02T17:05:05Z",
                "End": "2006-01-02T18:06:05Z",
                "Place": "Saint Petersburg"
            }
        ]
    }