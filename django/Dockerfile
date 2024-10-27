FROM python:3.12.5-slim

USER root

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONFAULTHANDLER 1
ENV PIPENV_VENV_IN_PROJECT 1

RUN pip install pipenv

WORKDIR /home/my-user/app

CMD tail -f /dev/null
