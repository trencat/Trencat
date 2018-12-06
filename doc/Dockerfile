FROM python:3.6.7-alpine3.8 as builder
COPY requirements.txt .
RUN pip3 install -r requirements.txt && \
    rm -r /root/.cache/
# Create user with no privileges, no password (-D) and no home directory (-H)
RUN adduser -D -H trencat_doc_builder
# Execute container with user with no privileges
USER trencat_doc_builder