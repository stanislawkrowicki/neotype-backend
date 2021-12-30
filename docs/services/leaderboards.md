## Leaderboards service holds best scores in Redis and lets you fetch best of them.

> **POST /leaderboards**
> THIS ENDPOINT IS USED INTERNALLY ONLY. YOU CAN NOT ACCESS IT VIA GATEWAY.
>
> > **Headers:**
> >
> > - Authorization: Bearer token
>
> > **Request body:**
> >
> > - wpm: float32
> > - accuracy: float32
> > - time: int
>
> > **Returns:**
> >
> > - 400 Bad Request:
> > - 500 Internal Server Error:
> > - 200 OK:

> **GET /leaderboards/{number}**
>
> > **Params**
> >
> > - number: how many leaders to get (between 1 and 50)
>
> > **Returns:**
> >
> > - 400 Bad Request:
> >  - message
> > - 500 Internal Server Error:
> >  - message
> > - 200 OK:
> >  - JSON array with leaders
