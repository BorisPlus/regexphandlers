# Regex Handlers

[![Go Report Card](https://goreportcard.com/badge/github.com/BorisPlus/regexphandlers)](https://goreportcard.com/report/github.com/BorisPlus/regexphandlers)

Мне не хватало задейсвования в URI-path пакета `net\http` возможности их парсинга посредством регулярных выражений, подобно REST Api:

* http://domain.com/api/item/{numeric};
* http://domain.com/api/item/{numeric}/update;
* http://domain.com/api/item/{numeric}/delete;
* http://domain.com/api/notify/{email}.

И если там в URI-path, например, не `{numeric}` или `{email}` шаблоны, то даже не идти по этому пути, передавая вызов дальше по списку априори невалидных данных.

> __Замечание__: В итоговом проекте задейстовал шаблон `http://domain.com/{numeric}/{numeric}/{any}`.

## Если коротко

Итак, допустим уже [имееются](tests/handler.go) __штатные__ `net\http`-хендлеры:

* DefaultHandler
* VersionHandler
* GetHandler

То для задействования `Regex Handlers` их (`net\http`-хендлеры) достаточно [обернуть](tests/router.go) в конструкцию с регулярными выражениями для URI-path:

```go
func Handlers() regexphandlers.RegexpHandlers {
    return regexphandlers.NewRegexpHandlers(
        DefaultHandler{},
        *regexphandlers.NewRegexpHandler(
            `/api/version`,
            none,
            VersionHandler{},
        ),
        *regexphandlers.NewRegexpHandler(
            `/api/get/{numeric}/{string}`, // "parent_id", "child_name"
            ids, // { parent_id, child_name }
            GetHandler{},
        ),
    )
}
```

при этом указав порядок следования параметров в URI-path:

```go
var (
    none = regexphandlers.Params{}
    ids  = regexphandlers.Params{"parent_id", "child_name"}
)
```

и всё, теперь их можно использовать опять же в __штатном__ `ServeMux` пакета `net\http`:

```go
mux := http.NewServeMux()
mux.Handle("/api/", Handlers()) // Вот тут новый ServeMux на основе регулярок
server := http.Server{
    Addr:    net.JoinHostPort(host, fmt.Sprint(port)),
    Handler: mux,
}
```

## Тестирование

```bash
make test
go clean -testcache && go test -race -cover ./tests
ok      github.com/BorisPlus/regexphandlers/tests       2.097s  coverage: 76.5% of statements
```

## Минус балл за "код"

Да, необходимая к конструированию сторонним разработчиком (пользователем моего пакета) функция `func Handlers() regexphandlers.RegexpHandlers` выглядит громоздкой. Но я уверен, что со временем можно будет подмешать немного синтаксического сахара для упразднения длинного объявления "перекрестка" `*regexphandlers.NewRegexpHandler` из шаблонизированных регулярными выражениями путей `NewRegexpHandler`.
