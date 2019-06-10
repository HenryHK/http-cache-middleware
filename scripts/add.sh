curl \
    --include \
    --header "autopilotapikey:  ACCESS_KEY" \
    -d '
    {
        "contact": {
            "FirstName": "Lin",
            "LastName": "Han",
            "Email": "henryhan.hku@gmail.com",
            "custom": {
                "string--Test--Field": "This is a test"
            }
        }
    }
    ' \
    -H "Content-Type: application/json" \
    -X POST \
    http://localhost:8080/contact
