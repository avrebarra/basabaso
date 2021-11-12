# Endpoint Spec Draft

## Users

```bash
@request
POST /api/users/find?

@response:ok
{
    "code": "1",
    "msg": "operation success",
    "data": {
        "users": [],
        "total": 12,
    }
}
```

```bash
@request
POST /api/users/persist
{
    "user": {}
}

@response:ok
{
    "code": "1",
    "msg": "operation success",
    "data": {
        "user": {}
    }
}
```

## Words

```bash
@request
POST /api/words/find?

@response:ok
{
    "code": "1",
    "msg": "operation success",
    "data": {
        "words": [],
        "total": 1500,
    }
}
```

```bash
@request
POST /api/words/persist
{
    "word": {}
}

@response:ok
{
    "code": "1",
    "msg": "operation success",
    "data": {
        "word": {}
    }
}
```

```bash
@request
POST /api/words/vote
{
    "word_id": "xid",
    "direction": "up"
}

@response:ok
{
    "code": "1",
    "msg": "operation success",
    "data": {}
}
```

```bash
@request
POST /api/words/report
{
    "word_id": "xid",
    "reason_group": "offensive"
    "reason_explanation": "it is not necessarily offensive"
}

@response:ok
{
    "code": "1",
    "msg": "operation success",
    "data": {}
}
```
