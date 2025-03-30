#!/bin/bash

openssl req -x509 -newkey rsa:4096 -nodes -keyout ./mail/certs/key.pem -out ./mail/certs/cert.pem -days 365 -subj "/CN=localhost"