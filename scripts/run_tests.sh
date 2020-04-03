#!/usr/bin/env bash

flake8 \
    --ignore=W503,W504,F723 \
    --exclude maintainer/tests/ \
    maintainer && \
mypy \
    maintainer/ \
    --disallow-untyped-defs \
    --strict-equality \
    --show-error-codes \
    --warn-return-any \
    --ignore-missing-imports && \
coverage erase && \
pytest \
    maintainer/ \
    -q \
    --cov-config .coveragerc \
    --cov-report= \
    --cov && \
coverage report && \
coverage html
