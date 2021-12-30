## Results service is responsible for inserting and retrieving results from database.

### results-publisher works as an API, while results-consumer only consumes results from RabbitMQ.

### You will need both of them up and running.

> **POST /result**
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
> >  - message
> > - 401 Unauthorized:
> >  - message
> > - 500 Internal Server Error:
> >  - message
> > - 200 OK:
> >   - message: successfully added to the queue

> **GET /results/{number}**
>
> > **Headers:**
> >
> > - Authorization: Bearer token
>
> > **Params**
> >
> > - number: how many results to fetch(between 1 and 50)
>
> > **Returns:**
> >
> > - 401 Unauthorized:
> >  - message
> > - 200 OK:
> >  - JSON array with results
