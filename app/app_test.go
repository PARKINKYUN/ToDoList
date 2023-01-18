package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
	"todolist/model"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	os.Remove("./test.db")
	assert := assert.New(t)
	// 테스트 서버 open
	appHandler := MakeHandler()
	defer appHandler.Close()

	ts := httptest.NewServer(appHandler)
	defer ts.Close()

	// "/todos" 경로로 "name"에 값을 넣어 POST 요청
	res, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	// 응답 메시지의 Body 에서 값을 읽어와 todo 에 저장하여 비교
	var todo model.Todo
	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	id1 := todo.ID

	res, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo 2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo 2")
	id2 := todo.ID

	// "/todos" 로 GET 요청을 보내 리스트 배열을 가져와 todos 에 저장하여 id와 내용 비교
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos := []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.Equal(t.Name, "Test todo")
		} else if t.ID == id2 {
			assert.Equal(t.Name, "Test todo 2")
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}

	// "/complete-todo/" + id 의 체크박스 상태 변경
	res, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos = []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)

	for _, t := range todos {
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	// "/todos/" + id 로 delete 요청을 보냈을 때 삭제 여부 체크
	// http 패키지는 delete 와 put, patch 를 지원하지 않기 때문에 NewRequest 로 해결한다.
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos = []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)

	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}
}
