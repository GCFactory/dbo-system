#!/bin/bash

openssl req -x509 -newkey rsa:4096 -nodes -keyout ./mail/mail_certs/key.pem -out ./mail/mail_certs/cert.pem -days 365 -subj "/CN=localhost"