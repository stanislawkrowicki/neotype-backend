## Users service is responsible for everything connected to users account, including authorization.

> **POST /register**
>
> > **Request body:**
> >
> > - login: string: user login
> > - password: string: user password
>
> > **Returns:**
> >
> > - 400 Bad Request:
> >  - message: what went wrong
> > - 500 Internal Server Error:
> >  - message: what went wrong
> > - 200 OK:
> >   - message: register successful

> **POST /login**
>
> > **Request body:**
> >
> > - login: string: user login
> > - password: string: user password
>
> > **Returns:**
> >
> > - 400 Bad Request:
> >  - message: what went wrong
> > - 500 Internal Server Error:
> >  - message: what went wrong
> > - 200 OK:
> >  - message: register successful
> >  - token: JWT token used for Authorization

> **GET /data**
>
> > **Headers:**
> >
> > - Authorization: Bearer token
>
> > **Returns:**
> >
> > - 401 Unauthorized:
> >  - message: what went wrong
> > - 500 Internal Server Error:
> >  - message: what went wrong
> > - 200 OK:
> >  - JSON object with user data

> **GET /username**
>
> > **Headers:**
> >
> > - Authorization: Bearer token
>
> > **Returns:**
> >
> > - 401 Unauthorized:
> >  - message
> > - 404 Not Found:
> >  - message: user was not found
> > - 200 OK:
> >  - JSON object with username key
