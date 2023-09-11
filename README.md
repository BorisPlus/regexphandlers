# Regex Handlers

Мне не хватало возможности в URI-path пакета `net\http` возможности использования регулярных выражений, походобных REST Api:

* "http://domain.com/api/item/{numeric}";
* "http://domain.com/api/item/{numeric}/update";
* "http://domain.com/api/item/{numeric}/delete";
* "http://domain.com/api/notify/{email}".

И если там соответственно, например, не `{numeric}` или `{email}`, то даже не идти по тому пути, передавая вызов дальше.
  
## Если коротко

Итак, допустим уже [имееются](tests/handler.go) стандартные `net\http`-хендлеры:

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
            ids,
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

и всё, теперь их можно **спокойно** использовать в штатном `ServeMux` пакета `net\http`:

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

Да, `Handlers` выглядит громоздкой. Но я уверен, что со временем можно будет подмешать немного синтаксического сахара для упразднения длинного объявления `*regexphandlers.NewRegexpHandler` в конструируемом "перекрестке" регулярных путей `RegexpHandlers`.
