# Planning Center Authentication CLI Example

This example uses the `pco-api-auth` package to perform authentication from a command line application.

`$ pco-auth -port 8080 -client_id 219b174b5a701c3e235dbd2401425c4738e48b9933fa6120bb573157fc5f692d -client_secret 3a95d09e3a7c8f28aed5d9a3c36d3916dd7cbacb357fdc67cc11126065a11582 -scope people,services`

STDOUT will receive the output:

```text
{
 "AccessToken": {
  "Token": "3b89b92bc72c10e1f6ce091f55835fecfee237d78b98c577665e33567962f006",
  "Kind": "Bearer",
  "ExpiresIn": 7200,
  "Refresh": "f9d7cb9cd382d9ecf605a3260f047c8e89d4d16aa489d02c2c4475e7b17c06ae",
  "Scope": "people services",
  "CreatedAt": 1569358429
 },
 "CurrentPerson": {
  "OrganizationID": "123",
  "PersonID": "456"
 }
}
```
