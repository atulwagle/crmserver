# This Go program is an implementation of CRM Server.

## This program supports the following APIs:

* Getting a single customer through a `/customers/{id}` path
    ```curl --location --request GET 'http://localhost:3000/customers/f3f1cc7d-1f32-4652-b016-52dc3ac5bf50'
        Response:
        {
            "id": "f3f1cc7d-1f32-4652-b016-52dc3ac5bf50",
            "name": "John",
            "role": "Supervisor",
            "email": "john@yahoo.com",
            "phone": 4083457834,
            "contacted": true
        }
    ```

* Getting all customers through a the `/customers` path
    ```
        curl --location --request GET 'http://localhost:3000/customers'
        Response:
                [
                    {
                        "id": "f3f1cc7d-1f32-4652-b016-52dc3ac5bf50",
                        "name": "John",
                        "role": "Supervisor",
                        "email": "john@yahoo.com",
                        "phone": 4083457834,
                        "contacted": true
                    },
                    {
                        "id": "2f7a9959-084a-41d7-a85f-5430cbf6d90a",
                        "name": "Bill",
                        "role": "Tester",
                        "email": "bill@yahoo.com",
                        "phone": 4083557844,
                        "contacted": true
                    },
                    {
                        "id": "2532aacc-a5ae-4540-a22c-d04d889f142e",
                        "name": "Harry",
                        "role": "Manager",
                        "email": "harry@gmail.com",
                        "phone": 4083657854,
                        "contacted": false
                    }
                ]
    ```
* Creating a customer through a `/customers `path
    ```
        curl --location --request POST 'http://localhost:3000/customers' \
            --header 'Content-Type: application/json' \
            --data-raw '{
                "name": "Atul",
                "role": "developer",
                "email": "atul@example.com",
                "phone": 4084872632,
                "contacted": true
            }'
        Response:
            For successfully creating the user, API will return 201
            {
                "id": "137c5121-7aff-4ae5-9dd6-da4f562c9a21",
                "name": "Atul",
                "role": "developer",
                "email": "atul@example.com",
                "phone": 4084872632,
                "contacted": true
            }

            In case the same payload is submitted once again, API will return 409 and 
            Record already present.
    ```
* Updating a customer through a `/customers/{id} `path
    ```
    This API support both PATCH and PUT http methods. 

        curl --location --request PATCH 'http://localhost:3000/customers/137c5121-7aff-4ae5-9dd6-da4f562c9a21' 
            --header 'Content-Type: application/json'
            --data-raw '{
                "name": "Atul",
                "role": "developer",
                "email": "atul@gmail.com",
                "phone": 4084872655,
                "contacted": false
            }'

        Response:
            If the resource is properly updated the API returns 200 and the updated resource:
                {
                    "id": "137c5121-7aff-4ae5-9dd6-da4f562c9a21",
                    "name": "Atul",
                    "role": "developer",
                    "email": "atul@gmail.com",
                    "phone": 4084872655,
                    "contacted": false
                }

            If the ID is not associated with a resource then API return 404 and Record not found.
    ```
* Updating multiple customer through a `/customers` path
    ```
        This API support both PATCH and PUT http methods.
        Using this API user can update multiple resources.
        In the response of the API, there will be a http status associated with every customer entry.

            curl --location --request PATCH 'http://localhost:3000/customers' \
                --header 'Content-Type: application/json' \
                --data-raw '[
                    {
                        "id": "f3f1cc7d-1f32-4652-b016-52dc3ac5bf50",
                        "name": "John",
                        "role": "Supervisor",
                        "email": "john@yahoo.com",
                        "phone": 4083457834,
                        "contacted": true
                    },
                    {
                        "id": "2f7a9959-084a-41d7-a85f-5430cbf6d90a",
                        "name": "Bill",
                        "role": "Tester",
                        "email": "bill@yahoo.com",
                        "phone": 4083557844,
                        "contacted": true
                    }
                ]'

            Response:
                If the API gets executed without any exceptions then it will return 200.
                    [
                        {
                            "customerentry": {
                                "id": "f3f1cc7d-1f32-4652-b016-52dc3ac5bf50",
                                "name": "John",
                                "role": "Supervisor",
                                "email": "john@yahoo.com",
                                "phone": 4083457834,
                                "contacted": true
                            },
                            "status": 200
                        },
                        {
                            "customerentry": {
                                "id": "2f7a9959-084a-41d7-a85f-5430cbf6d90a",
                                "name": "Bill",
                                "role": "Tester",
                                "email": "bill@yahoo.com",
                                "phone": 4083557844,
                                "contacted": true
                            },
                            "status": 200
                        }
                    ]
                
                If any resource could not be updated/found, API will return response such as:

                    [
                        {
                            "customerentry": {
                                "id": "f3f1cc7d-1f32-4652-b016-52dc3ac5bggg",
                                "name": "John",
                                "role": "Supervisor",
                                "email": "john@yahoo.com",
                                "phone": 4083457834,
                                "contacted": true
                            },
                            "status": 404
                        },
                        {
                            "customerentry": {
                                "id": "2f7a9959-084a-41d7-a85f-5430cbf6d90a",
                                "name": "Bill",
                                "role": "Tester",
                                "email": "bill@yahoo.com",
                                "phone": 4083557844,
                                "contacted": true
                            },
                            "status": 200
                        }
                    ]

    ```
* Deleting a customer through a `/customers/{id}` path
    ```
        curl --location --request DELETE 'http://localhost:3000/customers/f3f1cc7d-1f32-4652-b016-52dc3ac5bf50'

        Response:
            If the API is able to delete the resource it return 200 and contents of the customers db.

            [
                {
                    "id": "2f7a9959-084a-41d7-a85f-5430cbf6d90a",
                    "name": "Bill",
                    "role": "Tester",
                    "email": "bill@yahoo.com",
                    "phone": 4083557844,
                    "contacted": true
                },
                {
                    "id": "2532aacc-a5ae-4540-a22c-d04d889f142e",
                    "name": "Harry",
                    "role": "Manager",
                    "email": "harry@gmail.com",
                    "phone": 4083657854,
                    "contacted": false
                },
                {
                    "id": "0ea8f336-7c45-4d81-a17d-b736b1ca4d28",
                    "name": "Atul",
                    "role": "developer",
                    "email": "atul@example.com",
                    "phone": 4084872632,
                    "contacted": true
                }
            ]

            If we try to delete the same resource once again or try to delete using an invalid id, then the API returns 404 and Record not found.
    ```

    ## This program creates logs.txt that records the api calls and logs the customer id. It gets created in the current directory.

    ## References:

    * Go References: 
        * [Introduction to Programming in GO](https://www.golang-book.com/books/intro)
        * [Adding logs](https://www.honeybadger.io/blog/golang-logging/)
        * [Go Tour](https://go.dev/tour/list)
        * [Serving Static Sites](https://www.alexedwards.net/blog/serving-static-sites-with-go)
        * [Struts and JSON](https://drstearns.github.io/tutorials/gojson/)
    * Create Markdown:
        * [Markdown Live Preview](https://markdownlivepreview.com/)