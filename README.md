# GoSheets
![Gopher](/img/favicon.png)

## General informations
The purpose of this web app is to send some data to Google Sheets with the available API.  
It is designed to record expenses.

The HTML files are currently only in french.


## Technical informations

#### Built with 
- [go](https://go.dev/) & [gin-gonic](https://gin-gonic.com/)
- [htmx](https://htmx.org/)
- [pico](https://picocss.com/)
- [gopherize](https://gopherize.me/) for the nice logo


#### Prerequisite
You need to make a [Google Service Account](https://developers.google.com/workspace/guides/create-credentials#service-account) to get the following credentials.


#### To run the app
- generate environment variables :
    ```bash
    export type="service_account"
    export project_id="project"
    export private_key_id="XY"
    export private_key="-----BEGIN PRIVATE KEY-----\nXYZ\n-----END PRIVATE KEY-----\n"
    export client_email="X@Y.iam.gserviceaccount.com"
    export client_id="1"
    export auth_uri="https://accounts.google.com/o/oauth2/auth"
    export token_uri="https://oauth2.googleapis.com/token"
    export auth_provider_x509_cert_url="https://www.googleapis.com/oauth2/v1/certs"
    export client_x509_cert_url="https://www.googleapis.com/robot/v1/metadata/x509/X%Y.iam.gserviceaccount.com"
    export universe_domain="googleapis.com" 
    ```
- locally :
    ```bash
    go run .
    ```
- with Docker :
    ```bash
    docker build --tag name/gosheets:tag .

    docker run --detach -e type -e project_id -e private_key_id -e private_key -e client_email -e client_id -e auth_uri -e token_uri -e auth_provider_x509_cert_url -e client_x509_cert_url -e universe_domain --publish 127.0.0.1:8082:8082 imageIdJustBuilt
    ```

## TODO
- Améliorer page Insert Rows
- Ajout onglet de paramétrage dans Gsheets.
    - variable sur la liste des catégories
- Voir utilité GET Sheets


## Changelog
- 2023-09-24 : add read all gsheet, start to use params in a new gsheet.
- 2023-09-13 : initialize project