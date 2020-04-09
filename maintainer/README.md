Relay Maintainer

## Setup

install `pipenv` and `pyenv`

```sh
$ pipenv install --python=$(pyenv which python3.7)
```

## Testing
Current testscript runs linting and typechecking only.

```sh
$ pipenv run test
```
