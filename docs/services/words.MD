## Words service is used to generate words for typing tests.

> **GET /words/{number}**
> 
> **Params:**
>   - number: int: the number of words to fetch.
> 
> **Returns:**
> - 400 Bad Request: _number_ is not an integer
> - 200 OK:
>   - JSON array containing words.

