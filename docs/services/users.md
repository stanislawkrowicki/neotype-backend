## Users service is responsible for everything connected to users account, including authorization.

> **POST /register**
>>
>> **Request params:**
>> - login: string: user login
>> - password: string: user password
>
>> **Returns:**
>> - 400 Bad Request:
>>  - message: what went wrong
>> - 500 Internal Server Error:
>>  - message: what went wrong
>> - 200 OK:
>>    - message: register successful


> **POST /login**
>>
>> **Request params:**
>> - login: string: user login
>> - password: string: user password
>
>> **Returns:**
>> - 400 Bad Request:
>>  - message: what went wrong
>> - 500 Internal Server Error:
>>  - message: what went wrong
>> - 200 OK:
>>  - message: register successful
>>  - token: JWT token used for Authorization


> **GET /data**
>>
>> **Headers:**
>> - Authorization: Bearer token
>
>> **Returns:**
>> - 401 Unauthorized:
>>  - message: what went wrong
>> - 500 Internal Server Error:
>>  - message: what went wrong
>> - 200 OK:
>>  - JSON object with user data
