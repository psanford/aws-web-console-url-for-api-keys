# aws-web-console-url-for-api-keys

This tool converts aws api keys to an aws web console login url.
It uses the federation api to generate a temporary login to access
the aws console with the api keys.

You need the following environment variables set:
```
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
AWS_SESSION_TOKEN
```
