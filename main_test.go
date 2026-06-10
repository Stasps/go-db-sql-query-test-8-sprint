package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// TestSelectClientWhenOk() проверяет,
// что данные из таблицы корректно извлекаются.
// Идентификатор записи уже определён в тесте в переменной clientID.
func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	clientID := 1

	// Получаем клиента функцией selectClient()
	res, err := selectClient(db, clientID)
	// Если функция вернула ошибку — завершаем тест
	if err != nil {
		t.Fatalf("error calling selectClient(): %v", err)
	}

	// Проверяем, что ID клиента совпадает с ожидаемым
	assert.Equal(t, clientID, res.ID, "client ID should match the expected value")

	// Проверяем, что остальные поля не пустые
	assert.NotEmpty(t, res.FIO, "FIO field must not be empty")
	assert.NotEmpty(t, res.Login, "Login field must not be empty")
	assert.NotEmpty(t, res.Birthday, "Birthday field must not be empty")
	assert.NotEmpty(t, res.Email, "Email field must not be empty")
}

// TestSelectClientWhenNoClient() проверяет,
// что в случае отсутствия записи возвращается ошибка.
// Идентификатор записи уже определён в тесте в переменной clientID.
// Здесь намеренно задан идентификатор, которого нет в таблице.
func Test_SelectClient_WhenNoClient(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	clientID := -1

	// Получаем клиента функцией selectClient()
	res, err := selectClient(db, clientID)

	// Проверяем, что функция вернула ошибку
	if err == nil {
		t.Fatal("expected an error when client not found, but got nil")
	}

	// Проверяем, что ошибка равна sql.ErrNoRows
	assert.Equal(t, sql.ErrNoRows, err, "expected sql.ErrNoRows when client not found")

	// Проверяем, что все поля объекта Client пустые.
	assert.Empty(t, res.FIO, "FIO field must be empty")
	assert.Empty(t, res.Login, "Login field must be empty")
	assert.Empty(t, res.Birthday, "Birthday field must be empty")
	assert.Empty(t, res.Email, "Email field must be empty")
}

// TestInsertClientThenSelectAndCheck() добавляет запись в таблицу,
// а затем проверяет, что она добавилась корректно.
// Объект Client уже определён в тесте в переменной cl.
func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// Добавляем запись в таблицу функцией insertClient().
	res, err := insertClient(db, cl)
	// Если функция вернула ошибку — завершаем тест
	require.NoError(t, err, "insertClient() should not return an error")
	// Если функция вернула не пустой идентификатор — завершаем тест
	// Проверяем, что ID > 0
	require.NotZero(t, res, "insertClient() should return non-zero ID")
	// Сохраняем полученный идентификатор в поле ID в переменной cl.
	cl.ID = int(res)

	// Функцией selectClient() получаем объект Client по идентификатору.
	got, err := selectClient(db, cl.ID)
	// Проверяем, что функция вернула пустую ошибку. Иначе завершить тест.
	require.NoError(t, err, "selectClient() should not return an error when client exists")
	// Проверяем, что значения всех полей полученного объекта совпадают со значениями полей объекта в переменной cl.
	require.Equal(t, cl.FIO, got.FIO, "FIO field should match")
	require.Equal(t, cl.Login, got.Login, "Login field should match")
	require.Equal(t, cl.Birthday, got.Birthday, "Birthday field should match")
	require.Equal(t, cl.Email, got.Email, "Email field should match")
	require.Equal(t, cl.ID, got.ID, "ID field should match")
}

// TestInsertClientDeleteClientThenCheck() проверяет удаление записи.
// Объект Client уже определён в тесте в переменной cl.
func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// Добавляем запись в таблицу функцией insertClient().
	res, err := insertClient(db, cl)
	// Если функция вернула ошибку — завершаем тест
	require.NoError(t, err, "insertClient() should not return an error")
	// Если функция вернула не пустой идентификатор — завершаем тест
	// Проверяем, что ID > 0
	require.NotZero(t, res, "insertClient() should return non-zero ID")
	// Сохраняем полученный идентификатор в поле ID в переменной cl.
	cl.ID = int(res)

	// Функцией selectClient() получаем объект Client по идентификатору.
	got, err := selectClient(db, cl.ID)
	// Проверяем, что функция вернула пустую ошибку. Иначе завершить тест.
	require.NoError(t, err, "selectClient() should not return an error when client exists")

	// Удаляем запись функцией deleteClient(). Если функция вернула ошибку - завершаем тест.
	err = deleteClient(db, cl.ID)
	require.NoError(t, err)

	// Получаем объект клиента функцией selectClient().
	got, err = selectClient(db, cl.ID)
	// Проверяем, что функция вернула ошибку и ошибка равна sql.ErrNoRows. Иначе завершить тест.
	require.Equal(t, sql.ErrNoRows, err, "expected sql.ErrNoRows after deletion")
	// Проверяем, что поля объекта пустые
	assert.Empty(t, got.FIO, "FIO should be empty after deletion")
	assert.Empty(t, got.Login, "Login should be empty after deletion")
	assert.Empty(t, got.Birthday, "Birthday should be empty after deletion")
	assert.Empty(t, got.Email, "Email should be empty after deletion")
}
